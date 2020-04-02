package notify

import "encoding/json"

type payloadItem struct {
	field   string
	message interface{}
}

type payload []payloadItem

func (p payload) MarshalJSON() ([]byte, error) {
	dict := map[string]interface{}{}

	for _, item := range p {
		dict[item.field] = item.message
	}

	return json.Marshal(dict)
}

// Personalisation is a slice of structs used to define placeholder values in a
// template, such as name or reference number.
// The struct should be structured such that the key is the name of the value
// in your template, and the value is what you expect to be substituted in the
// message.
type Personalisation []struct {
	Key   string
	Value interface{}
}

func (personalisation Personalisation) updatePayload(p payload) payload {
	dict := map[string]interface{}{}

	for _, item := range personalisation {
		dict[item.Key] = item.Value
	}

	return append(p, payloadItem{"personalisation", dict})
}

func (personalisation Personalisation) updateSMSPayload(p payload) payload {
	return personalisation.updatePayload(p)
}

func (personalisation Personalisation) updateEmailPayload(p payload) payload {
	return personalisation.updatePayload(p)
}

func (personalisation Personalisation) updatePersonalisationPayload(p payload) payload {
	return personalisation.updatePayload(p)
}

// Reference is a unique identifier you create. It identifies a single unique
// notification or a batch of notifications.
func Reference(referenceID string) CommonOption {
	return updatePayloadFunc(func(p payload) payload {
		return append(p, payloadItem{"reference", referenceID})
	})
}

// StatusCallback allows you to pass delivery status callback details at the
// time of sending a notification.
func StatusCallback(url, bearerToken string) CommonOption {
	return updatePayloadFunc(func(p payload) payload {
		return append(p, payloadItem{"status_callback_url", url}, payloadItem{"status_callback_bearer_token", bearerToken})
	})
}

// EmailReplyToID is the ID of the reply-to address to receive replies from
// users.
func EmailReplyToID(address string) SendEmailOption {
	return updatePayloadFunc(func(p payload) payload {
		return append(p, payloadItem{"email_reply_to_id", address})
	})
}

// SMSSenderID is a unique identifier for the sender of a text message.
func SMSSenderID(senderID string) SendSMSOption {
	return updatePayloadFunc(func(p payload) payload {
		return append(p, payloadItem{"sms_sender_id", senderID})
	})
}

type CommonOption interface {
	SendSMSOption
	SendEmailOption
}

type SendSMSOption interface {
	updateSMSPayload(payload) payload
}

type SendEmailOption interface {
	updateEmailPayload(payload) payload
}

type PersonalisationOption interface {
	updatePersonalisationPayload(payload) payload
}

type updatePayloadFunc func(payload) payload

func (fn updatePayloadFunc) updateSMSPayload(p payload) payload {
	return fn(p)
}

func (fn updatePayloadFunc) updateEmailPayload(p payload) payload {
	return fn(p)
}
