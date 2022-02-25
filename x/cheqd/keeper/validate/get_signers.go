package validate

func (msg *MsgUpdateDidPayload) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, controller := range msg.Controller {
			if controller == msg.Id {
				result[i] = Signer{
					Signer:             controller,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: controller,
				}
			}
		}

		return result
	}

	if len(msg.Authentication) > 0 {
		return []Signer{
			{
				Signer:             msg.Id,
				Authentication:     msg.Authentication,
				VerificationMethod: msg.VerificationMethod,
			},
		}
	}

	return []Signer{}
}

func (msg *MsgCreateDidPayload) GetSigners() []Signer {
	if len(msg.Controller) > 0 {
		result := make([]Signer, len(msg.Controller))

		for i, controller := range msg.Controller {
			if controller == msg.Id {
				result[i] = Signer{
					Signer:             controller,
					Authentication:     msg.Authentication,
					VerificationMethod: msg.VerificationMethod,
				}
			} else {
				result[i] = Signer{
					Signer: controller,
				}
			}
		}

		return result
	}

	if len(msg.Authentication) > 0 {
		return []Signer{
			{
				Signer:             msg.Id,
				Authentication:     msg.Authentication,
				VerificationMethod: msg.VerificationMethod,
			},
		}
	}

	return []Signer{}
}
