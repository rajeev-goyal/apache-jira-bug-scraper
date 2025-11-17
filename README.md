## Bug Analysis Tool for Apache JIRA Issues

The following are required to run the tool:
1. Clone the Git repository of the target project.

2. Edit the `config.json` file to specify project details and analysis parameters. (Given below)

3. Run using Go. The commands are:
   ```bash
   go mod tidy
   go run .
   ```

### Edit the `config.json` file
Before running the tool, make sure to edit the `config.json` file for your system. An example configuration is:

```json
{
  "project_name": "KAFKA",
  "local_repo_path": "../kafka",
  "max_bugs_to_find": 100,
  "output_csv_file": "kafka_results.csv",
  "jira_base_url": "https://issues.apache.org/jira/browse/",
  "repo_commit_url": "https://github.com/apache/kafka/commit/",
  "analysis_keywords": [
    "EmbeddedZookeeper",
    "KafkaServerTestHarness",
    "kafka.server.KafkaServer",
    "IntegrationTestUtils",
    "kafka.utils.TestUtils"
  ]
}

```
### Get all commit hashes without any keywords

You can leave the analysis keywords empty to get the commit hashes for the bugs without filtering for minicluster tests. The tool finds the hashes by looking for the issue IDs in the commit messages.

- `project_name`: The JIRA project name to analyze (e.g., "KAFKA").
- `local_repo_path`: The local path to the cloned Git repository of the project.
- `max_bugs_to_find`: The maximum number of bugs to analyze.
- `output_csv_file`: The name of the output CSV file to store results.
- `jira_base_url`: The base URL for the JIRA issues. It will remain the same for Apache projects.
- `repo_commit_url`: The base URL for the repository commits.
- `analysis_keywords`: A list of keywords to search for the minicluster tests in the patch.


