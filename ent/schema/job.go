package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Job holds the schema definition for the Job entity.
type Job struct {
	ent.Schema
}

/*
    id: 1,
    hiring: false,
    categoryId: 2,
    title: 'Senior Full Stack Developer',
    slug: 'senior-full-stack-developer',
    location: 'Remote, Earth',
    summary:
      'Are you looking to make an impact with a dynamic team that is 100% remote with challenging opportunities? Then consider joining our team',
    type: 'Full-Time',
    category: 'Engineering',
    thumbnail: '/assets/images/cardOne.png',
    weHave: [
      'Employer-paid health care',
      'Perks for electricity, internet and more',
      'Casual and diverse workplace',
      'Paid vacation and team events',
    ],
    requirements: [
      'Energized to join a startup',
      'Excited to mentor more junior developers',
      'Good at problem selection, problem',
      'Solving, and course correcting',
      'Focused on best practices Highly pragmatic and collaborative',
    ],
    youHave: [
      'At least 4+ years of professional work experience as a Software Developer',
      'Proficient in Python and PHP',
      'Understanding of OOP principles and practices',
      'Experience working with relational databases MySQL and PostgreSQL',
      'You enjoy working with legacy code. Perhaps you are that one guy who has read Micheal Feathers',
      'Version control and Git workflow',
      'You enjoy Front-End development',
      'Understand Containers, Infrastructure as Code and have worked with various cloud providers',
      'You enjoy writing tests',
    ],
  },
*/

// Fields of the Job.
func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).
			Default(uuid.New),
	
		field.Time("created_at").
			Default(time.Now),

		field.Time("updated_at").
			Default(time.Now),
		
		field.Time("deleted_at"),

		field.Bool("hiring").
			Default(false),

		field.String("title"),

		field.String("slug"),

		field.String("location").
			Default("Remote, Earth"),

		field.String("summary"),

		field.Enum("employment").
			Values("Part Time", "Full Time", "Contract"),

		field.Enum("category").
			Values("Engineering", "Product & Design"),

		field.String("thumbnail"),

		field.JSON("wehave", []string{}),

		field.JSON("requirements", []string{}),

		field.JSON("youhave", []string{}),
	}
}

// Edges of the Job.
func (Job) Edges() []ent.Edge {
	return nil
}
