package b24

import (
	"path"
	"strconv"
	"time"

	"github.com/asm2me/asterlink/connect"
	log "github.com/sirupsen/logrus"
)

func (b *b24) OrigStart(c *connect.Call, oID string) {
	b.ent[c.LID] = &entity{ID: oID, cID: c.CID, log: b.log.WithField("lid", c.LID)}
}

func (b *b24) Start(c *connect.Call) {
	b.ent[c.LID] = &entity{cID: c.CID, log: b.log.WithField("lid", c.LID)}
	e := b.ent[c.LID]

	e.mux.Lock()

	uID, ok := b.eUID[c.Ext]
	if !ok {
		uID = b.cfg.DefUID
	}

	var params struct {
		UID   int    `json:"USER_ID"`
		Phone string `json:"PHONE_NUMBER"`
		Type  int    `json:"TYPE"`
		DID   string `json:"LINE_NUMBER"`
	}

	params.UID = uID
	params.Phone = c.CID
	params.DID = c.DID

	if c.Dir == connect.Out {
		params.Type = 1
	} else {
		params.Type = 2
	}

	var r struct {
		Result struct {
			ID string `json:"CALL_ID"`
		}
	}
	err := b.req("telephony.externalcall.register", params, &r)
	// TODO: ERROR HANDLING!!!
	e.ID = r.Result.ID
	e.log.WithField("id", e.ID).Debug("uID:"+strconv.Itoa(uID)+"===========Call.register !!=============")

	if err != nil {
		delete(b.ent, c.LID)
		e.mux.Unlock()

		return
	}


	e.mux.Unlock()
}

func (b *b24) Dial(c *connect.Call, ext string) {
	b.handleDial(c, ext, true)
}

func (b *b24) StopDial(c *connect.Call, ext string) {
	b.handleDial(c, ext, false)
}

func (b *b24) Answer(c *connect.Call, ext string) {
}

func (b *b24) End(c *connect.Call, cause string) {
	e, ok := b.ent[c.LID]
	if !ok || !e.isRegistred() {
		return
	}
	defer delete(b.ent, c.LID)

//	uID, ok := b.eUID[c.Ext]
	
	
	var r struct {
		Result []struct {
			ID    int    `json:"ID,string"`
			Phone string `json:"UF_PHONE_INNER"`
		}
	}
	err := b.req("user.get", map[string]map[string]string{
		"filter": {"UF_PHONE_INNER": c.Ext},
	}, &r)

	if err != nil {
		b.log.Error("Failed to Get UserID from Extension")
		return
	}	

 

	var uID = r.Result[0].ID
	e.log.WithField("id", e.ID).Debug("uID:"+strconv.Itoa(uID)+"===========Call Finished=============")
	
	
	
	
	var params struct {
		ID     string `json:"CALL_ID"`
		UID    int    `json:"USER_ID"`
		Dur    int    `json:"DURATION"`
		Status string `json:"STATUS_CODE"`
		Vote   int    `json:"VOTE,omitempty"`
	}

	params.ID = e.ID
	params.UID = uID

	if !c.TimeAnswer.IsZero() {
		params.Dur = int(time.Since(c.TimeAnswer).Seconds())
		params.Status = "200"
		e.log.Debug("call time=0 ---------------------")		
	} else {
		params.Dur = int(time.Since(c.TimeCall).Seconds())
		e.log.Debug("call time:"+strconv.Itoa( params.Dur) )		

		if cause == "16" {
			if c.Dir == connect.In {
				params.Status = "304" // This call was skipped
			} else {
				params.Status = "603-S" // This call was canceled
			}
		} else {
			params.Status, ok = b.causeCode[cause]
			if !ok {
				params.Status = "505" // Undefined
			}
		}
	}

	if c.Vote != "" && c.Vote != "-" {
		params.Vote, _ = strconv.Atoi(c.Vote)
	}

	b.req("telephony.externalcall.finish", params, nil)
	
	
	// TODO: HANDLE ERROR!!!!
	// upload recording
	e.log.Debug("b.cfg.RecUp:"+b.cfg.RecUp )		
	e.log.Debug("c.TimeAnswer:"+c.TimeAnswer.String() )
	e.log.Debug("c.Rec:"+c.Rec )
	
	if b.cfg.RecUp != "" &&params.Dur>0&&  c.Rec != "" {
		file := path.Base(c.Rec)
		url := b.cfg.RecUp + c.Rec

		e.log.WithFields(log.Fields{url: url}).Debug("URL:"+url+"====================Attaching call record=====================")
		b.req("telephony.externalCall.attachRecord", map[string]string{
			"CALL_ID":    e.ID,
			"FILENAME":   file,
			"RECORD_URL": url,
		}, nil)
	}
}

func (b *b24) handleDial(c *connect.Call, ext string, isDial bool) {
	e, ok := b.ent[c.LID]
	if !ok || !e.isRegistred() {
		return
	}
	
	e.log.WithField("id", e.ID).Debug("EXT:"+ c.Ext  +"===========Call Dialled=============")
	
	
	
	
	
	var r struct {
		Result []struct {
			ID    int    `json:"ID,string"`
			Phone string `json:"UF_PHONE_INNER"`
		}
	}
	err := b.req("user.get", map[string]map[string]string{
		"filter": {"UF_PHONE_INNER": ext},
	}, &r)

	if err != nil {
		b.log.Error("Failed to update users list")
		return
	}	

 

	var uID = r.Result[0].ID
	e.log.WithField("id", e.ID).Debug("uID:"+strconv.Itoa(uID)+"......................Call Dialled ")
	
	method := "telephony.externalcall."

	if isDial {
		method += "show"
	} else {
		method += "hide"
	}

	var params struct {
		ID  string `json:"CALL_ID"`
		UID int    `json:"USER_ID"`
	}
	params.ID = e.ID
	params.UID = uID

	b.req(method, params, nil)
}
