package server

import (
	"encoding/json"
	"fmt"
	"github.com/emqx/kuiper/common"
	"github.com/emqx/kuiper/plugins"
	"github.com/emqx/kuiper/xstream/sinks"
	"strings"
	"time"
)

const QUERY_RULE_ID = "internal-xstream_query_rule"

type Server int

func (t *Server) CreateQuery(sql string, reply *string) error {
	if _, ok := registry.Load(QUERY_RULE_ID); ok {
		stopQuery()
	}
	tp, err := ruleProcessor.ExecQuery(QUERY_RULE_ID, sql)
	if err != nil {
		return err
	} else {
		rs := &RuleState{Name: QUERY_RULE_ID, Topology: tp, Triggered: true}
		registry.Store(QUERY_RULE_ID, rs)
		msg := fmt.Sprintf("Query was submit successfully.")
		logger.Println(msg)
		*reply = fmt.Sprintf(msg)
	}
	return nil
}

func stopQuery() {
	if rs, ok := registry.Load(QUERY_RULE_ID); ok {
		logger.Printf("stop the query.")
		(*rs.Topology).Cancel()
		registry.Delete(QUERY_RULE_ID)
	}
}

/**
 * qid is not currently used.
 */
func (t *Server) GetQueryResult(qid string, reply *string) error {
	if rs, ok := registry.Load(QUERY_RULE_ID); ok {
		c := (*rs.Topology).GetContext()
		if c != nil && c.Err() != nil {
			return c.Err()
		}
	}

	sinks.QR.LastFetch = time.Now()
	sinks.QR.Mux.Lock()
	if len(sinks.QR.Results) > 0 {
		*reply = strings.Join(sinks.QR.Results, "")
		sinks.QR.Results = make([]string, 10)
	} else {
		*reply = ""
	}
	sinks.QR.Mux.Unlock()
	return nil
}

func (t *Server) Stream(stream string, reply *string) error {
	content, err := streamProcessor.ExecStmt(stream)
	if err != nil {
		return fmt.Errorf("Stream command error: %s", err)
	} else {
		for _, c := range content {
			*reply = *reply + fmt.Sprintln(c)
		}
	}
	return nil
}

func (t *Server) CreateRule(rule *common.RuleDesc, reply *string) error {
	r, err := ruleProcessor.ExecCreate(rule.Name, rule.Json)
	if err != nil {
		return fmt.Errorf("Create rule error : %s.", err)
	} else {
		*reply = fmt.Sprintf("Rule %s was created, please use 'cli getstatus rule $rule_name' command to get rule status.", rule.Name)
	}
	//Start the rule
	rs, err := createRuleState(r)
	if err != nil {
		return err
	}
	err = doStartRule(rs)
	if err != nil {
		return err
	}
	return nil
}

func (t *Server) GetStatusRule(name string, reply *string) error {
	if r, err := getRuleStatus(name); err != nil {
		return err
	} else {
		*reply = r
	}
	return nil
}

func (t *Server) StartRule(name string, reply *string) error {
	if err := startRule(name); err != nil {
		return err
	} else {
		*reply = fmt.Sprintf("Rule %s was started", name)
	}
	return nil
}

func (t *Server) StopRule(name string, reply *string) error {
	*reply = stopRule(name)
	return nil
}

func (t *Server) RestartRule(name string, reply *string) error {
	err := restartRule(name)
	if err != nil {
		return err
	}
	*reply = fmt.Sprintf("Rule %s was restarted.", name)
	return nil
}

func (t *Server) DescRule(name string, reply *string) error {
	r, err := ruleProcessor.ExecDesc(name)
	if err != nil {
		return fmt.Errorf("Desc rule error : %s.", err)
	} else {
		*reply = r
	}
	return nil
}

func (t *Server) ShowRules(_ int, reply *string) error {
	r, err := ruleProcessor.ExecShow()
	if err != nil {
		return fmt.Errorf("Show rule error : %s.", err)
	} else {
		*reply = r
	}
	return nil
}

func (t *Server) DropRule(name string, reply *string) error {
	stopRule(name)
	r, err := ruleProcessor.ExecDrop(name)
	if err != nil {
		return fmt.Errorf("Drop rule error : %s.", err)
	} else {
		err := t.StopRule(name, reply)
		if err != nil {
			return err
		}
	}
	*reply = r
	return nil
}

func (t *Server) CreatePlugin(arg *common.PluginDesc, reply *string) error {
	pt := plugins.PluginType(arg.Type)
	p, err := getPluginByJson(arg)
	if err != nil {
		return fmt.Errorf("Create plugin error: %s", err)
	}
	if p.File == "" {
		return fmt.Errorf("Create plugin error: Missing plugin file url.")
	}
	err = pluginManager.Register(pt, p)
	if err != nil {
		return fmt.Errorf("Create plugin error: %s", err)
	} else {
		*reply = fmt.Sprintf("Plugin %s is created.", p.Name)
	}
	return nil
}

func (t *Server) DropPlugin(arg *common.PluginDesc, reply *string) error {
	pt := plugins.PluginType(arg.Type)
	p, err := getPluginByJson(arg)
	if err != nil {
		return fmt.Errorf("Drop plugin error: %s", err)
	}
	err = pluginManager.Delete(pt, p.Name)
	if err != nil {
		return fmt.Errorf("Drop plugin error: %s", err)
	} else {
		*reply = fmt.Sprintf("Plugin %s is dropped.", p.Name)
	}
	return nil
}

func (t *Server) ShowPlugins(arg int, reply *string) error {
	pt := plugins.PluginType(arg)
	l, err := pluginManager.List(pt)
	if err != nil {
		return fmt.Errorf("Drop plugin error: %s", err)
	} else {
		if len(l) == 0 {
			l = append(l, "No plugin is found.")
		}
		*reply = strings.Join(l, "\n")
	}
	return nil
}

func getPluginByJson(arg *common.PluginDesc) (*plugins.Plugin, error) {
	var p plugins.Plugin
	if arg.Json != "" {
		if err := json.Unmarshal([]byte(arg.Json), &p); err != nil {
			return nil, fmt.Errorf("Parse plugin %s error : %s.", arg.Json, err)
		}
	}
	p.Name = arg.Name
	return &p, nil
}

func init() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-ticker.C
			if _, ok := registry.Load(QUERY_RULE_ID); !ok {
				continue
			}

			n := time.Now()
			w := 10 * time.Second
			if v := n.Sub(sinks.QR.LastFetch); v >= w {
				logger.Printf("The client seems no longer fetch the query result, stop the query now.")
				stopQuery()
				ticker.Stop()
				return
			}
		}
	}()
}
