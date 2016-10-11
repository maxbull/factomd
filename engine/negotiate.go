// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package engine

import (
	"time"

	"github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/state"
)

func Negotiate(s *state.State) {
	zeroBytes := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	time.Sleep(3 * time.Second)
	for {
		pl := s.ProcessLists.LastList()
		if pl != nil && pl.LenFaultMap() > 0 {
			if pl.ChosenNegotiation != zeroBytes {
				faultState := pl.GetFaultState(pl.ChosenNegotiation)
				if !faultState.IsNil() {
					if faultState.AmINegotiator {
						fullFault := state.CraftAndSubmitFullFault(pl, pl.ChosenNegotiation)
						if faultState.HasEnoughSigs(s) && faultState.PledgeDone {

							ack := s.NewAck(fullFault).(*messages.Ack)
							ack.SetVMIndex(int(fullFault.VMIndex))
							s.NetworkOutMsgQueue() <- ack

							break
						}
					}
				}
			} else {
				faultIDs := pl.GetKeysFaultMap()
				for _, faultID := range faultIDs {
					faultState := pl.GetFaultState(faultID)
					if faultState.AmINegotiator {
						fullFault := state.CraftAndSubmitFullFault(pl, faultID)
						if faultState.HasEnoughSigs(s) && faultState.PledgeDone {
							pl.ChosenNegotiation = faultID

							ack := s.NewAck(fullFault).(*messages.Ack)
							ack.SetVMIndex(int(fullFault.VMIndex))
							s.NetworkOutMsgQueue() <- ack

							break
						}
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
