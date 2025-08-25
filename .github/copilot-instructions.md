# GitHub Copilot Instructions

## Service Context
**IMPORTANT**: Always read and understand the `README.md` file first to get complete context about this service architecture, API endpoints, and implementation details before providing any assistance.

## Pull Request Description Generator

When a developer provides a PR title, generate a comprehensive PR description using this template:

```markdown
## 📋 Summary
[Brief description of what this PR accomplishes]

## 🎯 Agenda
- [ ] [List specific tasks/changes being made]

## 💼 Business Impact
- **Benefits**: [What business value this provides]
- **Risk Assessment**: [Low/Medium/High with brief explanation]

## 🔧 Technical Changes
- **Modified Files**: [List changed files with purpose]
- **API Changes**: [New/modified endpoints if any]
- **Database Changes**: [Schema/migration details if any]
- **Configuration**: [Environment variables or config changes]

## 🧪 Testing Strategy
- **Testing Approach**: [How to verify these changes]
- **Commands**: [Relevant make commands or curl examples]

## ⚠️ Deployment Notes
- **Migration Required**: [Yes/No with details]
- **Service Restart**: [Required/Not required]
- **Rollback Plan**: [Brief rollback approach]

## 🔍 Review Checklist
- [ ] Code follows project patterns from README
- [ ] Proper error handling implemented
- [ ] Tests added/updated as needed
- [ ] Documentation updated if required
```

## Code Standards
- Follow the patterns and conventions shown in the README.md
- Use the same error handling, logging, and configuration approaches
- Maintain consistency with existing codebase structure
- Reference the API examples and architecture described in README

## Analysis Instructions
1. **Read README.md** for service context
2. **Analyze PR title** to understand change scope
3. **Identify change type**: feature, bugfix, optimization, etc.
4. **Generate appropriate description** using the template above
5. **Suggest relevant testing** based on README examples
