package event

import amqp "github.com/rabbitmq/amqp091-go"

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

func (e *Emitter) Push(event, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emitter, nil
}
