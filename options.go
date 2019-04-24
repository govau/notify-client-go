package notify

import "encoding/json"

type PayloadItem struct {
	field   string
	message interface{}
}

type Payload []PayloadItem

func (payload Payload) MarshalJSON() ([]byte, error) {
	dict := map[string]interface{}{}

	for _, item := range payload {
		dict[item.field] = item.message
	}

	return json.Marshal(dict)
}

// Personalisation is a slice of structs used to define placeholder values in a template,
// such as name or reference number.
// The struct should be structured such that they key is the name of the value in your template, and the value is what you expect to be substituted in the message.
type Personalisation []struct{ Key, Value string }

func (personalisation Personalisation) UpdatePayload(payload Payload) Payload {
	dict := map[string]string{}

	for _, item := range personalisation {
		dict[item.Key] = item.Value
	}

	return append(payload, PayloadItem{"personalisation", dict})
}

func (personalisation Personalisation) UpdateSMSPayload(p Payload) Payload {
	return personalisation.UpdatePayload(p)
}

func (personalisation Personalisation) UpdateEmailPayload(p Payload) Payload {
	return personalisation.UpdatePayload(p)
}

// Reference is a unique identifier you create. It identifies a single unique notification or a batch of notifications.
func Reference(referenceID string) CommonOption {
	return CommonOption{
		UpdatePayloadFunc(func(payload Payload) Payload {
			return append(payload, PayloadItem{"reference", referenceID})
		}),
	}
}

// EmailReplyToID is the id of the reply-to address to receive replies from users.
func EmailReplyToID(address string) SendEmailOption {
	return UpdatePayloadFunc(func(payload Payload) Payload {
		return append(payload, PayloadItem{"email_reply_to_id", address})
	})
}

// SMSSenderID is a unique identifier for the sender of a text message.
func SMSSenderID(senderID string) SendSMSOption {
	return UpdatePayloadFunc(func(payload Payload) Payload {
		return append(payload, PayloadItem{"sms_sender_id", senderID})
	})
}

type PayloadUpdater interface {
	UpdatePayload(Payload) Payload
}

type CommonOption struct {
	PayloadUpdater
}

type SendSMSOption interface {
	UpdateSMSPayload(Payload) Payload
}

type SendEmailOption interface {
	UpdateEmailPayload(Payload) Payload
}

func (co CommonOption) UpdateSMSPayload(p Payload) Payload {
	return co.UpdatePayload(p)
}

func (co CommonOption) UpdateEmailPayload(p Payload) Payload {
	return co.UpdatePayload(p)
}

type UpdatePayloadFunc func(Payload) Payload

func (fn UpdatePayloadFunc) UpdatePayload(p Payload) Payload {
	return fn(p)
}

func (fn UpdatePayloadFunc) UpdateSMSPayload(p Payload) Payload {
	return fn(p)
}

func (fn UpdatePayloadFunc) UpdateEmailPayload(p Payload) Payload {
	return fn(p)
}
