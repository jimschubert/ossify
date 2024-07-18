## TODO

- Start with conventions
  - Build out the framework: report, support policies (OPA), auto-apply/fix?
- Restructure license, convention, and models to be flatter structures and packages by behaviors
- Full text search for existing licenses
- Create a generator or other script for syncing license data
  - Don't run on every build (maybe once per week?)
  - Pull only plain text, not _all_ license data from the opensource org's repo
- Move config code for versioning into cmd package
- Validate licenses behavior
  - Create, List, Search
  - Comment on the open source license of these licenses (provide attribution in root command --help)
- Fix the 'new' command for creating new licenses or conventions
