# Easy Commits

An AI-powered git commit message generator that creates detailed, structured commit messages based on your git diff and optional user context.

## Features

- **Multiple AI Providers**: Supports OpenAI, Anthropic Claude, and local Ollama models
- **Smart Diff Detection**: Automatically detects staged changes, falls back to unstaged changes if none staged
- **User Context**: Add additional context to help the AI understand your changes
- **Interactive Confirmation**: Shows generated message and asks for confirmation before committing
- **Automatic Staging**: Stages all changes before committing (runs `git add .`)
- **Claude Thinking Mode**: Advanced reasoning capabilities for Anthropic Claude models

## Installation

### Option 1: Build from Source
1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd easy-commits
   ```
2. Build the binary:
   ```bash
   go build -o easy-commits
   ```
3. (Optional) Move to your PATH:
   ```bash
   sudo mv easy-commits /usr/local/bin/
   ```

### Option 2: Using Go Install
```bash
go install github.com/your-username/easy-commits@latest
```

## Configuration

Configure your AI provider before first use:

```bash
easy-commits config
```

### Provider Options

#### OpenAI
- **Provider**: `openai`
- **API Key**: Your OpenAI API key
- **Model**: Automatically set to `gpt-3.5-turbo`

#### Anthropic Claude
- **Provider**: `anthropic`
- **API Key**: Your Anthropic API key
- **Model**: Choose from available Claude models (fetched dynamically)
- **Token Limit**: Configure max tokens (default: 500, thinking mode: 2048)
- **Thinking Mode**: Enable advanced reasoning for Opus/Sonnet models
  - Budget tokens: Minimum 1024 for thinking mode

#### Ollama (Local)
- **Provider**: `ollama`
- **Base URL**: Default `http://localhost:11434`
- **Model**: Your local model name (e.g., `llama2`, `codellama`, `mistral`)

Configuration is saved to `~/.easy-commits-config.json`

## Usage

### Basic Usage

Generate and create a commit with AI-generated message:
```bash
easy-commits commit
```

The tool will:
1. Check for staged changes (uses `git diff --cached`)
2. If no staged changes, check for unstaged changes (uses `git diff`)
3. Generate a commit message using AI
4. Show the message for review
5. Stage all changes and commit if approved

### With Additional Context

Provide extra context to help the AI understand your changes:
```bash
easy-commits commit --context "Fixed the authentication bug that was causing login failures"
```

### Example Workflow

1. Make your changes to files
2. Run `easy-commits commit`
3. Review the generated commit message
4. Type `y` to confirm or `n` to cancel
5. Changes are automatically staged and committed

## Generated Commit Message Format

The tool generates structured commit messages with:

1. **Subject line** (≤50 characters) - Brief summary in imperative mood
2. **Bullet points** - Summarized changes from the diff
3. **Detailed body** - Context and reasoning for the changes

Example output:
```
Implement user authentication system

- Add login and registration forms
- Create user model and database migration
- Implement password hashing and verification
- Set up session management

This change lays the foundation for user accounts and secure access to the application. It addresses the security requirements outlined in ticket #123.
```

## Commands

| Command | Description |
|---------|-------------|
| `easy-commits config` | Configure AI provider and credentials |
| `easy-commits commit` | Generate and create a commit |
| `easy-commits commit --context "description"` | Add user context to the commit generation |
| `easy-commits help` | Show help information |

## Requirements

- **Go**: 1.21 or later
- **Git**: Must be in a git repository
- **API Access**: Valid API key for chosen AI provider (or local Ollama setup)

## Error Handling

- **Not in git repo**: Tool checks for git repository before proceeding
- **No changes**: Detects when there are no changes to commit
- **API errors**: Provides clear error messages for API failures
- **Configuration**: Prompts to run config if no configuration found

## License

MIT License