package examples

import (
	"log"
	"os"
	"strings"

	"github.com/metamogul/ekstatic"
)

// Email composition workflow:

// States

type (
	emailEmpty emptyState

	emailWithSubject struct {
		subject string
	}

	emailWithGreeting struct {
		emailWithSubject
		greeting string
	}

	emailWithBodyCompositor struct {
		emailWithGreeting
		compositor textBodyCompositor
	}

	emailWithBody struct {
		emailWithGreeting
		textBody string
	}

	emailCompleted struct {
		emailWithBody
		closing string
	}
)

// Transition inputs

type (
	subject            string
	textBodyCompositor struct {
		*ekstatic.WorkflowInstance
	}
	greeting string
	closing  string

	printCommand emptyInput
)

// Transitions

type emailCompositingService struct {
	logger *log.Logger
}

func (c *emailCompositingService) addSubject(e emailEmpty, s subject) emailWithSubject {
	c.logger.Printf("Subject added: \"%s\"", s)
	return emailWithSubject{string(s)}
}

func (c *emailCompositingService) addGreeting(e emailWithSubject, g greeting) emailWithGreeting {
	c.logger.Printf("Greeting added: \"%s\"", g)
	return emailWithGreeting{e, string(g)}
}

func (c *emailCompositingService) addBodyCompositor(e emailWithGreeting, t textBodyCompositor) emailWithBodyCompositor {
	c.logger.Printf("Body compositor added")
	return emailWithBodyCompositor{e, t}
}

func (c *emailCompositingService) addParagraph(e emailWithBodyCompositor, p paragraph) emailWithBodyCompositor {
	c.logger.Printf("Paragraph added to email: \"%s\"", p)
	_ = e.compositor.ContinueWith(p)
	return e
}

func (c *emailCompositingService) buildBody(e emailWithBodyCompositor, t toStringWithLinebreak) emailWithBody {
	c.logger.Printf("Building body")
	_ = e.compositor.ContinueWith(t)
	c.logger.Printf("Body built:\n\n%v\n\n", e.compositor.CurrentState())
	return emailWithBody{e.emailWithGreeting, e.compositor.CurrentState().(string)}
}

func (c *emailCompositingService) addClosing(e emailWithBody, cl closing) emailCompleted {
	c.logger.Printf("Closing added: \"%s\"", cl)
	return emailCompleted{e, string(cl)}
}

func (c *emailCompositingService) printEmail(e emailCompleted, p printCommand) emailCompleted {
	email := "Subject: \"" + e.subject + "\"\n"
	email += e.greeting + "\n"
	email += e.textBody + "\n"
	email += e.closing + "\n"

	c.logger.Printf("Finished email:\n\n%s\n", email)

	return e
}

// Text body composition workflow:

// States

type textWithParagraphs struct {
	paragraphs []string
}

// Transition inputs

type (
	paragraph             string
	toStringWithLinebreak string
)

// Transitions

type bodyCompositingService struct {
	logger *log.Logger
}

func (b *bodyCompositingService) addParagraphToBody(t textWithParagraphs, p paragraph) textWithParagraphs {
	b.logger.Printf("Paragraph added to body compositor: \"%s\"", p)
	return textWithParagraphs{append(t.paragraphs, string(p))}
}

func (b *bodyCompositingService) convertToString(t textWithParagraphs, lineBreak toStringWithLinebreak) string {
	b.logger.Println("Building body in body compositor")
	return strings.Join(t.paragraphs, string(lineBreak))
}

// Global workflow definitions

var emailWorkflowSingleton *ekstatic.Workflow

func defineEmailCompositionWorkflow(e *emailCompositingService) *ekstatic.Workflow {
	if emailWorkflowSingleton != nil {
		return emailWorkflowSingleton
	}

	emailWorkflowSingleton = ekstatic.NewWorkflow()
	emailWorkflowSingleton.AddTransitions(
		e.addSubject,
		e.addGreeting,
		e.addBodyCompositor,
		e.addParagraph,
		e.buildBody,
		e.addClosing,
		e.printEmail,
	)

	return emailWorkflowSingleton
}

var emailBodyWorkflowSingleton *ekstatic.Workflow

func defineEmailBodyCompositionWorkflow(b *bodyCompositingService) *ekstatic.Workflow {
	if emailBodyWorkflowSingleton != nil {
		return emailBodyWorkflowSingleton
	}

	emailBodyWorkflowSingleton = ekstatic.NewWorkflow()
	emailBodyWorkflowSingleton.AddTransitions(
		b.addParagraphToBody,
		b.convertToString,
	)

	return emailBodyWorkflowSingleton
}

// Test

func ExampleWorkflow_submachine() {
	emailLogger := log.New(os.Stdout, "[EmailCompositingService] ", 0)
	emailCompositingService := &emailCompositingService{emailLogger}
	emailWorkflow := defineEmailCompositionWorkflow(emailCompositingService)

	bodyLogger := log.New(os.Stdout, "[BodyCompositingService]  ", 0)
	bodyCompositingService := &bodyCompositingService{bodyLogger}
	emailBodyWorkflow := defineEmailBodyCompositionWorkflow(bodyCompositingService)

	emailCompositor := emailWorkflow.New(emailEmpty{})

	_ = emailCompositor.ContinueWith(subject("ekstatic"))
	_ = emailCompositor.ContinueWith(greeting("Dear reader of this test,"))

	emailBodyCompositor := emailBodyWorkflow.New(textWithParagraphs{})

	_ = emailCompositor.ContinueWith(textBodyCompositor{emailBodyCompositor})
	_ = emailCompositor.ContinueWith(paragraph("This example is supposed to demonstrate how to implement a submachine."))
	_ = emailCompositor.ContinueWith(paragraph("I hope it will be insightful to you."))
	_ = emailCompositor.ContinueWith(toStringWithLinebreak("\n"))
	_ = emailCompositor.ContinueWith(closing("Best regards"))
	_ = emailCompositor.ContinueWith(printCommand{})

	// Output:
	// [EmailCompositingService] Subject added: "ekstatic"
	// [EmailCompositingService] Greeting added: "Dear reader of this test,"
	// [EmailCompositingService] Body compositor added
	// [EmailCompositingService] Paragraph added to email: "This example is supposed to demonstrate how to implement a submachine."
	// [BodyCompositingService]  Paragraph added to body compositor: "This example is supposed to demonstrate how to implement a submachine."
	// [EmailCompositingService] Paragraph added to email: "I hope it will be insightful to you."
	// [BodyCompositingService]  Paragraph added to body compositor: "I hope it will be insightful to you."
	// [EmailCompositingService] Building body
	// [BodyCompositingService]  Building body in body compositor
	// [EmailCompositingService] Body built:
	//
	// This example is supposed to demonstrate how to implement a submachine.
	// I hope it will be insightful to you.
	//
	// [EmailCompositingService] Closing added: "Best regards"
	// [EmailCompositingService] Finished email:
	//
	// Subject: "ekstatic"
	// Dear reader of this test,
	// This example is supposed to demonstrate how to implement a submachine.
	// I hope it will be insightful to you.
	// Best regards
}
