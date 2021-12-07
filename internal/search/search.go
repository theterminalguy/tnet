// Package search allows you to search and filter Repository objects using Ransack matchers
// Example matchers include:
//   - *_eq: Equal
//   - *_gt: Greater Than
//   - *_gte: Greater Than or Equal
//   - *_lt: Less Than
//   - *_lte: Less Than or Equal
//
// Aside from specifying matchers, you can also specify modifiers:
//   - *_asc: Sort Ascending
//   - *_desc: Sort Descending
package search

// TODO: Let's define interfaces which allows us to implment search using any data source. 
// For example we can use a database(postgres, MySQL, elastic search), a file, a cache, a service, etc.

// TODO: While still using this "rough" implementation we need to allow searching with a date-period range
// And also allow for searching with database column type of jsonb
// e.g PrimaryTechnology column in WorkExperience table.
