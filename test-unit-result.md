user@userland:/media/code/code/Go/garnizeh/englog$ make test-unit
Running unit tests...
go test -mod=mod -v -short -tags='!integration' ./internal/...
=== RUN   TestNewOllamaService
=== RUN   TestNewOllamaService/Valid_parameters
=== RUN   TestNewOllamaService/Empty_base_URL
=== RUN   TestNewOllamaService/Nil_logger
--- PASS: TestNewOllamaService (0.00s)
    --- PASS: TestNewOllamaService/Valid_parameters (0.00s)
    --- PASS: TestNewOllamaService/Empty_base_URL (0.00s)
    --- PASS: TestNewOllamaService/Nil_logger (0.00s)
=== RUN   TestGenerateInsight
=== RUN   TestGenerateInsight/Successful_generation
=== RUN   TestGenerateInsight/Empty_prompt
=== RUN   TestGenerateInsight/Server_error
=== RUN   TestGenerateInsight/Invalid_JSON_response
--- PASS: TestGenerateInsight (6.01s)
    --- PASS: TestGenerateInsight/Successful_generation (0.00s)
    --- PASS: TestGenerateInsight/Empty_prompt (0.00s)
    --- PASS: TestGenerateInsight/Server_error (3.01s)
    --- PASS: TestGenerateInsight/Invalid_JSON_response (3.00s)
=== RUN   TestGenerateInsightWithContext
=== RUN   TestGenerateInsightWithContext/Valid_request_with_string_context
=== RUN   TestGenerateInsightWithContext/Valid_request_with_structured_context
=== RUN   TestGenerateInsightWithContext/Valid_request_with_nil_context
=== RUN   TestGenerateInsightWithContext/Empty_prompt
--- PASS: TestGenerateInsightWithContext (0.00s)
    --- PASS: TestGenerateInsightWithContext/Valid_request_with_string_context (0.00s)
    --- PASS: TestGenerateInsightWithContext/Valid_request_with_structured_context (0.00s)
    --- PASS: TestGenerateInsightWithContext/Valid_request_with_nil_context (0.00s)
    --- PASS: TestGenerateInsightWithContext/Empty_prompt (0.00s)
=== RUN   TestBuildEnhancedPrompt
=== RUN   TestBuildEnhancedPrompt/Complete_request_with_string_context
=== RUN   TestBuildEnhancedPrompt/Request_with_structured_context
=== RUN   TestBuildEnhancedPrompt/Request_with_many_entry_IDs_(truncated_display)
=== RUN   TestBuildEnhancedPrompt/Request_with_minimal_data
=== RUN   TestBuildEnhancedPrompt/Request_with_team_collaboration_type
=== RUN   TestBuildEnhancedPrompt/Request_with_custom_struct_context
--- PASS: TestBuildEnhancedPrompt (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Complete_request_with_string_context (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Request_with_structured_context (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Request_with_many_entry_IDs_(truncated_display) (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Request_with_minimal_data (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Request_with_team_collaboration_type (0.00s)
    --- PASS: TestBuildEnhancedPrompt/Request_with_custom_struct_context (0.00s)
=== RUN   TestBuildEnhancedPromptExample
    ollama_test.go:426: === Enhanced Prompt Example ===
        Please analyze my productivity patterns and provide actionable insights for improvement.

        --- Request Information ---
        User ID: user-12345
        Insight Type: productivity
        Number of Log Entries: 6
        Log Entry IDs: [entry-001, entry-002, entry-003, ... (1 more), entry-005, entry-006]

        Insight Generation Guidelines for 'productivity':
        - Focus on efficiency patterns, time utilization, and value delivery
        - Identify high-impact activities and optimization opportunities
        - Analyze work-life balance and sustainable productivity patterns

        --- Additional Context ---
        Structured Context:
        {
          "date_range": {
            "end": "2025-07-31",
            "start": "2025-07-01"
          },
          "focus_areas": [
            "development",
            "meetings",
            "documentation"
          ],
          "performance_metrics": {
            "avg_daily_hours": 8.5,
            "productivity_score": 0.85
          },
          "time_blocks": [
            "morning",
            "afternoon",
            "evening"
          ]
        }

        --- Output Instructions ---
        Please provide a comprehensive analysis that includes:
        1. Key findings and patterns identified
        2. Specific, actionable recommendations
        3. Confidence level in your analysis (high/medium/low)
        4. Suggested next steps or areas for deeper investigation
        === End of Example ===
--- PASS: TestBuildEnhancedPromptExample (0.00s)
=== RUN   TestValidateInsightRequest
=== RUN   TestValidateInsightRequest/Valid_request
=== RUN   TestValidateInsightRequest/Empty_prompt
=== RUN   TestValidateInsightRequest/Empty_user_ID
=== RUN   TestValidateInsightRequest/Empty_insight_type
=== RUN   TestValidateInsightRequest/Empty_entry_IDs
--- PASS: TestValidateInsightRequest (0.00s)
    --- PASS: TestValidateInsightRequest/Valid_request (0.00s)
    --- PASS: TestValidateInsightRequest/Empty_prompt (0.00s)
    --- PASS: TestValidateInsightRequest/Empty_user_ID (0.00s)
    --- PASS: TestValidateInsightRequest/Empty_insight_type (0.00s)
    --- PASS: TestValidateInsightRequest/Empty_entry_IDs (0.00s)
=== RUN   TestValidateContextForInsightType
=== RUN   TestValidateContextForInsightType/Nil_context
=== RUN   TestValidateContextForInsightType/String_context
=== RUN   TestValidateContextForInsightType/Valid_productivity_context
=== RUN   TestValidateContextForInsightType/Invalid_productivity_context_-_wrong_type
=== RUN   TestValidateContextForInsightType/Invalid_productivity_context_-_not_array
=== RUN   TestValidateContextForInsightType/Valid_skill_development_context
=== RUN   TestValidateContextForInsightType/Invalid_skill_development_context
=== RUN   TestValidateContextForInsightType/Valid_time_management_context
=== RUN   TestValidateContextForInsightType/Invalid_time_management_context_-_missing_start
=== RUN   TestValidateContextForInsightType/Invalid_time_management_context_-_missing_end
=== RUN   TestValidateContextForInsightType/Unknown_insight_type_with_valid_JSON
=== RUN   TestValidateContextForInsightType/Non-JSON_serializable_context
--- PASS: TestValidateContextForInsightType (0.00s)
    --- PASS: TestValidateContextForInsightType/Nil_context (0.00s)
    --- PASS: TestValidateContextForInsightType/String_context (0.00s)
    --- PASS: TestValidateContextForInsightType/Valid_productivity_context (0.00s)
    --- PASS: TestValidateContextForInsightType/Invalid_productivity_context_-_wrong_type (0.00s)
    --- PASS: TestValidateContextForInsightType/Invalid_productivity_context_-_not_array (0.00s)
    --- PASS: TestValidateContextForInsightType/Valid_skill_development_context (0.00s)
    --- PASS: TestValidateContextForInsightType/Invalid_skill_development_context (0.00s)
    --- PASS: TestValidateContextForInsightType/Valid_time_management_context (0.00s)
    --- PASS: TestValidateContextForInsightType/Invalid_time_management_context_-_missing_start (0.00s)
    --- PASS: TestValidateContextForInsightType/Invalid_time_management_context_-_missing_end (0.00s)
    --- PASS: TestValidateContextForInsightType/Unknown_insight_type_with_valid_JSON (0.00s)
    --- PASS: TestValidateContextForInsightType/Non-JSON_serializable_context (0.00s)
=== RUN   TestGenerateWeeklyReport
=== RUN   TestGenerateWeeklyReport/Successful_generation
=== RUN   TestGenerateWeeklyReport/Empty_user_ID
=== RUN   TestGenerateWeeklyReport/Server_error
--- PASS: TestGenerateWeeklyReport (6.01s)
    --- PASS: TestGenerateWeeklyReport/Successful_generation (0.00s)
    --- PASS: TestGenerateWeeklyReport/Empty_user_ID (0.00s)
    --- PASS: TestGenerateWeeklyReport/Server_error (6.01s)
=== RUN   TestHealthCheck
=== RUN   TestHealthCheck/Healthy_service
=== RUN   TestHealthCheck/Unhealthy_service
--- PASS: TestHealthCheck (0.00s)
    --- PASS: TestHealthCheck/Healthy_service (0.00s)
    --- PASS: TestHealthCheck/Unhealthy_service (0.00s)
=== RUN   TestGenerateWithTimeout
=== RUN   TestGenerateWithTimeout/Successful_generation_within_timeout
=== RUN   TestGenerateWithTimeout/Timeout_exceeded
--- PASS: TestGenerateWithTimeout (0.30s)
    --- PASS: TestGenerateWithTimeout/Successful_generation_within_timeout (0.10s)
    --- PASS: TestGenerateWithTimeout/Timeout_exceeded (0.20s)
=== RUN   TestContextCancellation
=== RUN   TestContextCancellation/GenerateInsight_context_cancellation
=== RUN   TestContextCancellation/GenerateWeeklyReport_context_cancellation
=== RUN   TestContextCancellation/HealthCheck_context_cancellation
--- PASS: TestContextCancellation (2.20s)
    --- PASS: TestContextCancellation/GenerateInsight_context_cancellation (0.10s)
    --- PASS: TestContextCancellation/GenerateWeeklyReport_context_cancellation (0.10s)
    --- PASS: TestContextCancellation/HealthCheck_context_cancellation (0.10s)
=== RUN   TestRetryMechanism
--- PASS: TestRetryMechanism (3.00s)
=== RUN   TestConcurrentRequests
--- PASS: TestConcurrentRequests (0.00s)
=== RUN   TestJSONSerialization
=== RUN   TestJSONSerialization/GenerateRequest
=== RUN   TestJSONSerialization/GenerateResponse
=== RUN   TestJSONSerialization/Insight
=== RUN   TestJSONSerialization/InsightRequest
=== RUN   TestJSONSerialization/WeeklyReportRequest
=== RUN   TestJSONSerialization/WeeklyReport
--- PASS: TestJSONSerialization (0.00s)
    --- PASS: TestJSONSerialization/GenerateRequest (0.00s)
    --- PASS: TestJSONSerialization/GenerateResponse (0.00s)
    --- PASS: TestJSONSerialization/Insight (0.00s)
    --- PASS: TestJSONSerialization/InsightRequest (0.00s)
    --- PASS: TestJSONSerialization/WeeklyReportRequest (0.00s)
    --- PASS: TestJSONSerialization/WeeklyReport (0.00s)
=== RUN   TestEdgeCases
=== RUN   TestEdgeCases/Very_long_prompt
=== RUN   TestEdgeCases/Complex_nested_context
=== RUN   TestEdgeCases/Malformed_server_response
--- PASS: TestEdgeCases (3.00s)
    --- PASS: TestEdgeCases/Very_long_prompt (0.00s)
    --- PASS: TestEdgeCases/Complex_nested_context (0.00s)
    --- PASS: TestEdgeCases/Malformed_server_response (3.00s)
PASS
ok      github.com/garnizeh/englog/internal/ai  20.553s
=== RUN   TestRequireAuthMiddleware
=== RUN   TestRequireAuthMiddleware/valid_bearer_token
=== RUN   TestRequireAuthMiddleware/missing_authorization_header
=== RUN   TestRequireAuthMiddleware/invalid_bearer_format
=== RUN   TestRequireAuthMiddleware/bearer_without_token
=== RUN   TestRequireAuthMiddleware/invalid_token
--- PASS: TestRequireAuthMiddleware (0.00s)
    --- PASS: TestRequireAuthMiddleware/valid_bearer_token (0.00s)
    --- PASS: TestRequireAuthMiddleware/missing_authorization_header (0.00s)
    --- PASS: TestRequireAuthMiddleware/invalid_bearer_format (0.00s)
    --- PASS: TestRequireAuthMiddleware/bearer_without_token (0.00s)
    --- PASS: TestRequireAuthMiddleware/invalid_token (0.00s)
=== RUN   TestOptionalAuthMiddleware
=== RUN   TestOptionalAuthMiddleware/valid_bearer_token
=== RUN   TestOptionalAuthMiddleware/missing_authorization_header
=== RUN   TestOptionalAuthMiddleware/invalid_bearer_format
=== RUN   TestOptionalAuthMiddleware/invalid_token
--- PASS: TestOptionalAuthMiddleware (0.00s)
    --- PASS: TestOptionalAuthMiddleware/valid_bearer_token (0.00s)
    --- PASS: TestOptionalAuthMiddleware/missing_authorization_header (0.00s)
    --- PASS: TestOptionalAuthMiddleware/invalid_bearer_format (0.00s)
    --- PASS: TestOptionalAuthMiddleware/invalid_token (0.00s)
=== RUN   TestCaseInsensitiveBearerToken
=== RUN   TestCaseInsensitiveBearerToken/lowercase_bearer
=== RUN   TestCaseInsensitiveBearerToken/uppercase_bearer
=== RUN   TestCaseInsensitiveBearerToken/mixed_case_bearer
=== RUN   TestCaseInsensitiveBearerToken/weird_case_bearer
--- PASS: TestCaseInsensitiveBearerToken (0.00s)
    --- PASS: TestCaseInsensitiveBearerToken/lowercase_bearer (0.00s)
    --- PASS: TestCaseInsensitiveBearerToken/uppercase_bearer (0.00s)
    --- PASS: TestCaseInsensitiveBearerToken/mixed_case_bearer (0.00s)
    --- PASS: TestCaseInsensitiveBearerToken/weird_case_bearer (0.00s)
=== RUN   TestMultipleSpacesInAuthHeader
=== RUN   TestMultipleSpacesInAuthHeader/single_space
=== RUN   TestMultipleSpacesInAuthHeader/multiple_spaces
=== RUN   TestMultipleSpacesInAuthHeader/tab_character
--- PASS: TestMultipleSpacesInAuthHeader (0.00s)
    --- PASS: TestMultipleSpacesInAuthHeader/single_space (0.00s)
    --- PASS: TestMultipleSpacesInAuthHeader/multiple_spaces (0.00s)
    --- PASS: TestMultipleSpacesInAuthHeader/tab_character (0.00s)
=== RUN   TestNewAuthService
--- PASS: TestNewAuthService (0.00s)
=== RUN   TestCreateAccessToken
--- PASS: TestCreateAccessToken (0.00s)
=== RUN   TestCreateRefreshToken
--- PASS: TestCreateRefreshToken (0.00s)
=== RUN   TestValidateTokenWithValidAccessToken
--- PASS: TestValidateTokenWithValidAccessToken (0.00s)
=== RUN   TestValidateTokenWithInvalidToken
=== RUN   TestValidateTokenWithInvalidToken/completely_invalid_token
=== RUN   TestValidateTokenWithInvalidToken/empty_token
=== RUN   TestValidateTokenWithInvalidToken/malformed_token
--- PASS: TestValidateTokenWithInvalidToken (0.00s)
    --- PASS: TestValidateTokenWithInvalidToken/completely_invalid_token (0.00s)
    --- PASS: TestValidateTokenWithInvalidToken/empty_token (0.00s)
    --- PASS: TestValidateTokenWithInvalidToken/malformed_token (0.00s)
=== RUN   TestValidateTokenWithWrongSigningMethod
--- PASS: TestValidateTokenWithWrongSigningMethod (0.00s)
=== RUN   TestHashPassword
=== RUN   TestHashPassword/simple_password
=== RUN   TestHashPassword/complex_password
=== RUN   TestHashPassword/empty_password
=== RUN   TestHashPassword/unicode_password
--- PASS: TestHashPassword (0.01s)
    --- PASS: TestHashPassword/simple_password (0.00s)
    --- PASS: TestHashPassword/complex_password (0.00s)
    --- PASS: TestHashPassword/empty_password (0.00s)
    --- PASS: TestHashPassword/unicode_password (0.00s)
=== RUN   TestCheckPassword
=== RUN   TestCheckPassword/incorrect_password_wrong-password
=== RUN   TestCheckPassword/incorrect_password_test-password-124
=== RUN   TestCheckPassword/incorrect_password_TEST-PASSWORD-123
=== RUN   TestCheckPassword/incorrect_password_
=== RUN   TestCheckPassword/incorrect_password_test-password-12
--- PASS: TestCheckPassword (0.01s)
    --- PASS: TestCheckPassword/incorrect_password_wrong-password (0.00s)
    --- PASS: TestCheckPassword/incorrect_password_test-password-124 (0.00s)
    --- PASS: TestCheckPassword/incorrect_password_TEST-PASSWORD-123 (0.00s)
    --- PASS: TestCheckPassword/incorrect_password_ (0.00s)
    --- PASS: TestCheckPassword/incorrect_password_test-password-12 (0.00s)
=== RUN   TestTokenExpiration
--- PASS: TestTokenExpiration (0.00s)
=== RUN   TestMultipleTokensUniqueness
--- PASS: TestMultipleTokensUniqueness (1.10s)
=== RUN   TestGenerateJTI
--- PASS: TestGenerateJTI (0.00s)
PASS
ok      github.com/garnizeh/englog/internal/auth        1.134s
?       github.com/garnizeh/englog/internal/config      [no test files]
?       github.com/garnizeh/englog/internal/database    [no test files]
=== RUN   TestWorkerCapabilityMatching
=== RUN   TestWorkerCapabilityMatching/AI_Insights_task_to_AI_worker
=== RUN   TestWorkerCapabilityMatching/Weekly_report_task_to_report_worker
=== RUN   TestWorkerCapabilityMatching/Data_analysis_task_to_multi_worker
=== RUN   TestWorkerCapabilityMatching/Notification_task_to_generic_worker
=== RUN   TestWorkerCapabilityMatching/Cleanup_task_to_any_worker
--- PASS: TestWorkerCapabilityMatching (0.00s)
    --- PASS: TestWorkerCapabilityMatching/AI_Insights_task_to_AI_worker (0.00s)
    --- PASS: TestWorkerCapabilityMatching/Weekly_report_task_to_report_worker (0.00s)
    --- PASS: TestWorkerCapabilityMatching/Data_analysis_task_to_multi_worker (0.00s)
    --- PASS: TestWorkerCapabilityMatching/Notification_task_to_generic_worker (0.00s)
    --- PASS: TestWorkerCapabilityMatching/Cleanup_task_to_any_worker (0.00s)
=== RUN   TestTaskPriorityAndDeadlines
=== RUN   TestTaskPriorityAndDeadlines/tasks_with_different_priorities
=== RUN   TestTaskPriorityAndDeadlines/tasks_with_past_deadlines
=== RUN   TestTaskPriorityAndDeadlines/tasks_with_very_long_deadlines
--- PASS: TestTaskPriorityAndDeadlines (0.00s)
    --- PASS: TestTaskPriorityAndDeadlines/tasks_with_different_priorities (0.00s)
    --- PASS: TestTaskPriorityAndDeadlines/tasks_with_past_deadlines (0.00s)
    --- PASS: TestTaskPriorityAndDeadlines/tasks_with_very_long_deadlines (0.00s)
=== RUN   TestTaskMetadata
=== RUN   TestTaskMetadata/task_with_rich_metadata
=== RUN   TestTaskMetadata/task_with_empty_metadata
=== RUN   TestTaskMetadata/task_with_nil_metadata
=== RUN   TestTaskMetadata/task_with_special_characters_in_metadata
--- PASS: TestTaskMetadata (0.00s)
    --- PASS: TestTaskMetadata/task_with_rich_metadata (0.00s)
    --- PASS: TestTaskMetadata/task_with_empty_metadata (0.00s)
    --- PASS: TestTaskMetadata/task_with_nil_metadata (0.00s)
    --- PASS: TestTaskMetadata/task_with_special_characters_in_metadata (0.00s)
=== RUN   TestWorkerStatsAndStatus
=== RUN   TestWorkerStatsAndStatus/heartbeat_with_detailed_stats
=== RUN   TestWorkerStatsAndStatus/heartbeat_with_zero_stats
=== RUN   TestWorkerStatsAndStatus/heartbeat_with_high_resource_usage
=== RUN   TestWorkerStatsAndStatus/heartbeat_with_error_status
=== RUN   TestWorkerStatsAndStatus/heartbeat_without_stats
--- PASS: TestWorkerStatsAndStatus (0.00s)
    --- PASS: TestWorkerStatsAndStatus/heartbeat_with_detailed_stats (0.00s)
    --- PASS: TestWorkerStatsAndStatus/heartbeat_with_zero_stats (0.00s)
    --- PASS: TestWorkerStatsAndStatus/heartbeat_with_high_resource_usage (0.00s)
    --- PASS: TestWorkerStatsAndStatus/heartbeat_with_error_status (0.00s)
    --- PASS: TestWorkerStatsAndStatus/heartbeat_without_stats (0.00s)
=== RUN   TestHealthCheckWithDifferentScenarios
=== RUN   TestHealthCheckWithDifferentScenarios/health_check_with_no_workers
    comprehensive_test.go:469:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/comprehensive_test.go:469
                Error:          Not equal:
                                expected: "healthy"
                                actual  : "warning"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -healthy
                                +warning
                Test:           TestHealthCheckWithDifferentScenarios/health_check_with_no_workers
    comprehensive_test.go:471:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/comprehensive_test.go:471
                Error:          map[string]string{"grpc_server":"healthy", "ollama":"unknown", "task_queue":"healthy", "worker_connections":"no_workers"} does not contain "grpc"
                Test:           TestHealthCheckWithDifferentScenarios/health_check_with_no_workers
    comprehensive_test.go:472:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/comprehensive_test.go:472
                Error:          Not equal:
                                expected: "healthy"
                                actual  : ""

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -healthy
                                +
                Test:           TestHealthCheckWithDifferentScenarios/health_check_with_no_workers
=== RUN   TestHealthCheckWithDifferentScenarios/health_check_with_workers_of_different_ages
    comprehensive_test.go:519:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/comprehensive_test.go:519
                Error:          map[string]string{"grpc_server":"healthy", "ollama":"unknown", "task_queue":"healthy", "worker_connections":"no_workers"} does not contain "grpc"
                Test:           TestHealthCheckWithDifferentScenarios/health_check_with_workers_of_different_ages
--- FAIL: TestHealthCheckWithDifferentScenarios (0.00s)
    --- FAIL: TestHealthCheckWithDifferentScenarios/health_check_with_no_workers (0.00s)
    --- FAIL: TestHealthCheckWithDifferentScenarios/health_check_with_workers_of_different_ages (0.00s)
=== RUN   TestManager_NewManager
=== RUN   TestManager_NewManager/creates_new_manager_successfully
--- PASS: TestManager_NewManager (0.00s)
    --- PASS: TestManager_NewManager/creates_new_manager_successfully (0.00s)
=== RUN   TestManager_Start
=== RUN   TestManager_Start/start_without_TLS
=== RUN   TestManager_Start/start_with_invalid_TLS_configuration
=== RUN   TestManager_Start/start_with_port_already_in_use
--- PASS: TestManager_Start (0.10s)
    --- PASS: TestManager_Start/start_without_TLS (0.00s)
    --- PASS: TestManager_Start/start_with_invalid_TLS_configuration (0.00s)
    --- PASS: TestManager_Start/start_with_port_already_in_use (0.10s)
=== RUN   TestManager_Stop
=== RUN   TestManager_Stop/stop_without_starting
=== RUN   TestManager_Stop/stop_after_starting
=== RUN   TestManager_Stop/multiple_stops
--- PASS: TestManager_Stop (0.00s)
    --- PASS: TestManager_Stop/stop_without_starting (0.00s)
    --- PASS: TestManager_Stop/stop_after_starting (0.00s)
    --- PASS: TestManager_Stop/multiple_stops (0.00s)
=== RUN   TestManager_QueueInsightGenerationTask
=== RUN   TestManager_QueueInsightGenerationTask/valid_insight_generation_task
=== RUN   TestManager_QueueInsightGenerationTask/empty_user_ID
=== RUN   TestManager_QueueInsightGenerationTask/empty_entry_IDs
=== RUN   TestManager_QueueInsightGenerationTask/many_entry_IDs
--- PASS: TestManager_QueueInsightGenerationTask (0.00s)
    --- PASS: TestManager_QueueInsightGenerationTask/valid_insight_generation_task (0.00s)
    --- PASS: TestManager_QueueInsightGenerationTask/empty_user_ID (0.00s)
    --- PASS: TestManager_QueueInsightGenerationTask/empty_entry_IDs (0.00s)
    --- PASS: TestManager_QueueInsightGenerationTask/many_entry_IDs (0.00s)
=== RUN   TestManager_QueueWeeklyReportTask
=== RUN   TestManager_QueueWeeklyReportTask/valid_weekly_report_task
=== RUN   TestManager_QueueWeeklyReportTask/empty_user_ID
=== RUN   TestManager_QueueWeeklyReportTask/invalid_date_range_(end_before_start)
=== RUN   TestManager_QueueWeeklyReportTask/future_dates
--- PASS: TestManager_QueueWeeklyReportTask (0.00s)
    --- PASS: TestManager_QueueWeeklyReportTask/valid_weekly_report_task (0.00s)
    --- PASS: TestManager_QueueWeeklyReportTask/empty_user_ID (0.00s)
    --- PASS: TestManager_QueueWeeklyReportTask/invalid_date_range_(end_before_start) (0.00s)
    --- PASS: TestManager_QueueWeeklyReportTask/future_dates (0.00s)
=== RUN   TestManager_GetTaskResult
=== RUN   TestManager_GetTaskResult/get_non-existent_task_result
=== RUN   TestManager_GetTaskResult/get_task_result_after_reporting
--- PASS: TestManager_GetTaskResult (0.00s)
    --- PASS: TestManager_GetTaskResult/get_non-existent_task_result (0.00s)
    --- PASS: TestManager_GetTaskResult/get_task_result_after_reporting (0.00s)
=== RUN   TestManager_GetActiveWorkers
=== RUN   TestManager_GetActiveWorkers/no_active_workers
=== RUN   TestManager_GetActiveWorkers/with_registered_workers
--- PASS: TestManager_GetActiveWorkers (0.00s)
    --- PASS: TestManager_GetActiveWorkers/no_active_workers (0.00s)
    --- PASS: TestManager_GetActiveWorkers/with_registered_workers (0.00s)
=== RUN   TestManager_HealthCheck
=== RUN   TestManager_HealthCheck/health_check_without_starting_server
=== RUN   TestManager_HealthCheck/health_check_with_running_server
=== RUN   TestManager_HealthCheck/health_check_with_context_timeout
--- PASS: TestManager_HealthCheck (0.00s)
    --- PASS: TestManager_HealthCheck/health_check_without_starting_server (0.00s)
    --- PASS: TestManager_HealthCheck/health_check_with_running_server (0.00s)
    --- PASS: TestManager_HealthCheck/health_check_with_context_timeout (0.00s)
=== RUN   TestManager_TaskQueueIntegration
=== RUN   TestManager_TaskQueueIntegration/queue_multiple_task_types
=== RUN   TestManager_TaskQueueIntegration/queue_tasks_and_retrieve_results
--- PASS: TestManager_TaskQueueIntegration (0.00s)
    --- PASS: TestManager_TaskQueueIntegration/queue_multiple_task_types (0.00s)
    --- PASS: TestManager_TaskQueueIntegration/queue_tasks_and_retrieve_results (0.00s)
=== RUN   TestManager_ConcurrentOperations
=== RUN   TestManager_ConcurrentOperations/concurrent_task_queuing
=== RUN   TestManager_ConcurrentOperations/concurrent_health_checks
--- PASS: TestManager_ConcurrentOperations (0.00s)
    --- PASS: TestManager_ConcurrentOperations/concurrent_task_queuing (0.00s)
    --- PASS: TestManager_ConcurrentOperations/concurrent_health_checks (0.00s)
=== RUN   TestServer_NewServer
=== RUN   TestServer_NewServer/creates_new_server_successfully
--- PASS: TestServer_NewServer (0.00s)
    --- PASS: TestServer_NewServer/creates_new_server_successfully (0.00s)
=== RUN   TestServer_RegisterWorker
=== RUN   TestServer_RegisterWorker/valid_worker_registration
=== RUN   TestServer_RegisterWorker/missing_worker_ID
=== RUN   TestServer_RegisterWorker/missing_worker_name
=== RUN   TestServer_RegisterWorker/re-registration_of_existing_worker
--- PASS: TestServer_RegisterWorker (0.00s)
    --- PASS: TestServer_RegisterWorker/valid_worker_registration (0.00s)
    --- PASS: TestServer_RegisterWorker/missing_worker_ID (0.00s)
    --- PASS: TestServer_RegisterWorker/missing_worker_name (0.00s)
    --- PASS: TestServer_RegisterWorker/re-registration_of_existing_worker (0.00s)
=== RUN   TestServer_WorkerHeartbeat
=== RUN   TestServer_WorkerHeartbeat/valid_heartbeat
=== RUN   TestServer_WorkerHeartbeat/worker_not_found
=== RUN   TestServer_WorkerHeartbeat/invalid_session_token
=== RUN   TestServer_WorkerHeartbeat/status_change_heartbeat
--- PASS: TestServer_WorkerHeartbeat (0.00s)
    --- PASS: TestServer_WorkerHeartbeat/valid_heartbeat (0.00s)
    --- PASS: TestServer_WorkerHeartbeat/worker_not_found (0.00s)
    --- PASS: TestServer_WorkerHeartbeat/invalid_session_token (0.00s)
    --- PASS: TestServer_WorkerHeartbeat/status_change_heartbeat (0.00s)
=== RUN   TestServer_QueueTask
=== RUN   TestServer_QueueTask/valid_task
=== RUN   TestServer_QueueTask/missing_task_ID
--- PASS: TestServer_QueueTask (0.00s)
    --- PASS: TestServer_QueueTask/valid_task (0.00s)
    --- PASS: TestServer_QueueTask/missing_task_ID (0.00s)
=== RUN   TestServer_StreamTasks
=== RUN   TestServer_StreamTasks/worker_not_found
=== RUN   TestServer_StreamTasks/invalid_session_token
=== RUN   TestServer_StreamTasks/successful_task_streaming_with_context_cancellation
--- PASS: TestServer_StreamTasks (0.01s)
    --- PASS: TestServer_StreamTasks/worker_not_found (0.00s)
    --- PASS: TestServer_StreamTasks/invalid_session_token (0.00s)
    --- PASS: TestServer_StreamTasks/successful_task_streaming_with_context_cancellation (0.01s)
=== RUN   TestServer_ReportTaskResult
=== RUN   TestServer_ReportTaskResult/successful_task_completion
=== RUN   TestServer_ReportTaskResult/failed_task
=== RUN   TestServer_ReportTaskResult/missing_task_ID
=== RUN   TestServer_ReportTaskResult/missing_worker_ID
--- PASS: TestServer_ReportTaskResult (0.00s)
    --- PASS: TestServer_ReportTaskResult/successful_task_completion (0.00s)
    --- PASS: TestServer_ReportTaskResult/failed_task (0.00s)
    --- PASS: TestServer_ReportTaskResult/missing_task_ID (0.00s)
    --- PASS: TestServer_ReportTaskResult/missing_worker_ID (0.00s)
=== RUN   TestServer_UpdateTaskProgress
=== RUN   TestServer_UpdateTaskProgress/valid_progress_update
=== RUN   TestServer_UpdateTaskProgress/progress_at_0%
=== RUN   TestServer_UpdateTaskProgress/progress_at_100%
=== RUN   TestServer_UpdateTaskProgress/invalid_progress_percentage_(negative)
=== RUN   TestServer_UpdateTaskProgress/invalid_progress_percentage_(over_100)
--- PASS: TestServer_UpdateTaskProgress (0.00s)
    --- PASS: TestServer_UpdateTaskProgress/valid_progress_update (0.00s)
    --- PASS: TestServer_UpdateTaskProgress/progress_at_0% (0.00s)
    --- PASS: TestServer_UpdateTaskProgress/progress_at_100% (0.00s)
    --- PASS: TestServer_UpdateTaskProgress/invalid_progress_percentage_(negative) (0.00s)
    --- PASS: TestServer_UpdateTaskProgress/invalid_progress_percentage_(over_100) (0.00s)
=== RUN   TestServer_HealthCheck
=== RUN   TestServer_HealthCheck/health_check_with_no_workers
    server_test.go:536:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/server_test.go:536
                Error:          Not equal:
                                expected: "healthy"
                                actual  : "warning"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -healthy
                                +warning
                Test:           TestServer_HealthCheck/health_check_with_no_workers
    server_test.go:539:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/server_test.go:539
                Error:          map[string]string{"grpc_server":"healthy", "ollama":"unknown", "task_queue":"healthy", "worker_connections":"no_workers"} does not contain "grpc"
                Test:           TestServer_HealthCheck/health_check_with_no_workers
    server_test.go:540:
                Error Trace:    /media/code/code/Go/garnizeh/englog/internal/grpc/server_test.go:540
                Error:          Not equal:
                                expected: "healthy"
                                actual  : ""

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -healthy
                                +
                Test:           TestServer_HealthCheck/health_check_with_no_workers
=== RUN   TestServer_HealthCheck/health_check_with_registered_workers
--- FAIL: TestServer_HealthCheck (0.00s)
    --- FAIL: TestServer_HealthCheck/health_check_with_no_workers (0.00s)
    --- PASS: TestServer_HealthCheck/health_check_with_registered_workers (0.00s)
=== RUN   TestServer_HelperMethods
=== RUN   TestServer_HelperMethods/GetActiveWorkers_with_no_workers
=== RUN   TestServer_HelperMethods/GetActiveWorkers_with_registered_workers
=== RUN   TestServer_HelperMethods/GetTaskResult_for_non-existent_task
--- PASS: TestServer_HelperMethods (0.00s)
    --- PASS: TestServer_HelperMethods/GetActiveWorkers_with_no_workers (0.00s)
    --- PASS: TestServer_HelperMethods/GetActiveWorkers_with_registered_workers (0.00s)
    --- PASS: TestServer_HelperMethods/GetTaskResult_for_non-existent_task (0.00s)
=== RUN   TestServer_Start
=== RUN   TestServer_Start/start_with_invalid_address
=== RUN   TestServer_Start/start_with_available_port
--- PASS: TestServer_Start (0.01s)
    --- PASS: TestServer_Start/start_with_invalid_address (0.00s)
    --- PASS: TestServer_Start/start_with_available_port (0.01s)
=== RUN   TestServer_TaskQueueFull
--- PASS: TestServer_TaskQueueFull (0.00s)
=== RUN   TestServer_ConcurrentAccess
=== RUN   TestServer_ConcurrentAccess/concurrent_worker_registrations
=== RUN   TestServer_ConcurrentAccess/concurrent_task_queueing
--- PASS: TestServer_ConcurrentAccess (0.00s)
    --- PASS: TestServer_ConcurrentAccess/concurrent_worker_registrations (0.00s)
    --- PASS: TestServer_ConcurrentAccess/concurrent_task_queueing (0.00s)
FAIL
FAIL    github.com/garnizeh/englog/internal/grpc        0.145s
=== RUN   TestAnalyticsHandler_GetProductivityMetrics_Comprehensive
=== RUN   TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_activity_data
2025/08/01 15:40:18 github.com/testcontainers/testcontainers-go - Connected to docker:
  Server Version: 28.3.2
  API Version: 1.50
  Operating System: Ubuntu 24.04.2 LTS
  Total Memory: 128652 MB
  Testcontainers for Go Version: v0.38.0
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 16d98066cadbd84934afe64ee52063ba475dce50cfd47090414910502fb648c8
  Test ProcessID: 19eded8d-c698-4cdb-8e90-eea7522e63ca
2025/08/01 15:40:18 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:18 🐳 Creating container for image testcontainers/ryuk:0.12.0
2025/08/01 15:40:18 ✅ Container created: d7fda81c30fd
2025/08/01 15:40:18 🐳 Starting container: d7fda81c30fd
2025/08/01 15:40:18 ✅ Container started: d7fda81c30fd
2025/08/01 15:40:18 ⏳ Waiting for container id d7fda81c30fd image: testcontainers/ryuk:0.12.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false skipExternalCheck:false}
2025/08/01 15:40:18 🔔 Container is ready: d7fda81c30fd
2025/08/01 15:40:18 ✅ Container created: b5e387a09b9c
2025/08/01 15:40:18 🐳 Starting container: b5e387a09b9c
2025/08/01 15:40:18 ✅ Container started: b5e387a09b9c
2025/08/01 15:40:18 ⏳ Waiting for container id b5e387a09b9c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0004767f8 Strategies:[0xc000714720 0xc000595c80]}
2025/08/01 15:40:20 🔔 Container is ready: b5e387a09b9c
DSN: postgres://testuser:testpass@localhost:36282/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36282/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:20 OK   000001_create_users_table.sql (11.91ms)
2025/08/01 15:40:20 OK   000002_create_projects_table.sql (12.65ms)
2025/08/01 15:40:20 OK   000003_create_log_entries_table.sql (19.35ms)
2025/08/01 15:40:20 OK   000004_create_tags_system.sql (16.88ms)
2025/08/01 15:40:20 OK   000005_create_auth_tables.sql (20.88ms)
2025/08/01 15:40:20 OK   000006_create_insights_table.sql (24.43ms)
2025/08/01 15:40:20 OK   000007_create_performance_indexes.sql (15.98ms)
2025/08/01 15:40:20 OK   000008_create_analytics_views.sql (30.05ms)
2025/08/01 15:40:20 OK   000009_development_data.sql (29.02ms)
2025/08/01 15:40:20 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36282/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36282/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:20 🐳 Stopping container: b5e387a09b9c
2025/08/01 15:40:21 ✅ Container stopped: b5e387a09b9c
2025/08/01 15:40:21 🐳 Terminating container: b5e387a09b9c
2025/08/01 15:40:21 🚫 Container terminated: b5e387a09b9c
=== RUN   TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_empty_date_range
2025/08/01 15:40:21 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:21 ✅ Container created: 8f394f40b7c6
2025/08/01 15:40:21 🐳 Starting container: 8f394f40b7c6
2025/08/01 15:40:21 ✅ Container started: 8f394f40b7c6
2025/08/01 15:40:21 ⏳ Waiting for container id 8f394f40b7c6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc000265880 Strategies:[0xc000715b00 0xc000ea1800]}
2025/08/01 15:40:23 🔔 Container is ready: 8f394f40b7c6
DSN: postgres://testuser:testpass@localhost:36283/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36283/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:23 OK   000001_create_users_table.sql (11.92ms)
2025/08/01 15:40:23 OK   000002_create_projects_table.sql (12.62ms)
2025/08/01 15:40:23 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:40:23 OK   000004_create_tags_system.sql (17.06ms)
2025/08/01 15:40:23 OK   000005_create_auth_tables.sql (20.87ms)
2025/08/01 15:40:23 OK   000006_create_insights_table.sql (25.18ms)
2025/08/01 15:40:23 OK   000007_create_performance_indexes.sql (16.18ms)
2025/08/01 15:40:23 OK   000008_create_analytics_views.sql (29.08ms)
2025/08/01 15:40:23 OK   000009_development_data.sql (29.33ms)
2025/08/01 15:40:23 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36283/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36283/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:23 🐳 Stopping container: 8f394f40b7c6
2025/08/01 15:40:23 ✅ Container stopped: 8f394f40b7c6
2025/08/01 15:40:23 🐳 Terminating container: 8f394f40b7c6
2025/08/01 15:40:23 🚫 Container terminated: 8f394f40b7c6
=== RUN   TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_future_date_range
2025/08/01 15:40:23 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:23 ✅ Container created: 49168beed08f
2025/08/01 15:40:23 🐳 Starting container: 49168beed08f
2025/08/01 15:40:24 ✅ Container started: 49168beed08f
2025/08/01 15:40:24 ⏳ Waiting for container id 49168beed08f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc000ebf868 Strategies:[0xc0001e3ec0 0xc0005fe7e0]}
2025/08/01 15:40:25 🔔 Container is ready: 49168beed08f
DSN: postgres://testuser:testpass@localhost:36284/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36284/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:25 OK   000001_create_users_table.sql (12.44ms)
2025/08/01 15:40:25 OK   000002_create_projects_table.sql (15.41ms)
2025/08/01 15:40:25 OK   000003_create_log_entries_table.sql (21.88ms)
2025/08/01 15:40:25 OK   000004_create_tags_system.sql (17.37ms)
2025/08/01 15:40:26 OK   000005_create_auth_tables.sql (20.02ms)
2025/08/01 15:40:26 OK   000006_create_insights_table.sql (21.88ms)
2025/08/01 15:40:26 OK   000007_create_performance_indexes.sql (13.22ms)
2025/08/01 15:40:26 OK   000008_create_analytics_views.sql (23.22ms)
2025/08/01 15:40:26 OK   000009_development_data.sql (24.33ms)
2025/08/01 15:40:26 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36284/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36284/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:26 🐳 Stopping container: 49168beed08f
2025/08/01 15:40:26 ✅ Container stopped: 49168beed08f
2025/08/01 15:40:26 🐳 Terminating container: 49168beed08f
2025/08/01 15:40:26 🚫 Container terminated: 49168beed08f
--- PASS: TestAnalyticsHandler_GetProductivityMetrics_Comprehensive (8.65s)
    --- PASS: TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_activity_data (3.15s)
    --- PASS: TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_empty_date_range (2.74s)
    --- PASS: TestAnalyticsHandler_GetProductivityMetrics_Comprehensive/productivity_metrics_with_future_date_range (2.75s)
=== RUN   TestAnalyticsHandler_GetActivitySummary_Comprehensive
=== RUN   TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_multiple_activities
2025/08/01 15:40:26 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:26 ✅ Container created: 31687ea9e4f6
2025/08/01 15:40:26 🐳 Starting container: 31687ea9e4f6
2025/08/01 15:40:26 ✅ Container started: 31687ea9e4f6
2025/08/01 15:40:26 ⏳ Waiting for container id 31687ea9e4f6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001e752b8 Strategies:[0xc00200aae0 0xc001e7baa0]}
2025/08/01 15:40:28 🔔 Container is ready: 31687ea9e4f6
DSN: postgres://testuser:testpass@localhost:36285/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36285/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:28 OK   000001_create_users_table.sql (9.86ms)
2025/08/01 15:40:28 OK   000002_create_projects_table.sql (11.64ms)
2025/08/01 15:40:28 OK   000003_create_log_entries_table.sql (18.85ms)
2025/08/01 15:40:28 OK   000004_create_tags_system.sql (16.15ms)
2025/08/01 15:40:28 OK   000005_create_auth_tables.sql (20.49ms)
2025/08/01 15:40:28 OK   000006_create_insights_table.sql (24.53ms)
2025/08/01 15:40:28 OK   000007_create_performance_indexes.sql (15.57ms)
2025/08/01 15:40:28 OK   000008_create_analytics_views.sql (27.68ms)
2025/08/01 15:40:28 OK   000009_development_data.sql (28.12ms)
2025/08/01 15:40:28 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36285/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36285/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:29 🐳 Stopping container: 31687ea9e4f6
2025/08/01 15:40:29 ✅ Container stopped: 31687ea9e4f6
2025/08/01 15:40:29 🐳 Terminating container: 31687ea9e4f6
2025/08/01 15:40:29 🚫 Container terminated: 31687ea9e4f6
=== RUN   TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_date_range_filter
2025/08/01 15:40:29 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:29 ✅ Container created: a13198c41fe8
2025/08/01 15:40:29 🐳 Starting container: a13198c41fe8
2025/08/01 15:40:29 ✅ Container started: a13198c41fe8
2025/08/01 15:40:29 ⏳ Waiting for container id a13198c41fe8 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003a0d430 Strategies:[0xc001e9cb40 0xc003a03d10]}
2025/08/01 15:40:31 🔔 Container is ready: a13198c41fe8
DSN: postgres://testuser:testpass@localhost:36286/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36286/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:31 OK   000001_create_users_table.sql (9.89ms)
2025/08/01 15:40:31 OK   000002_create_projects_table.sql (12.01ms)
2025/08/01 15:40:31 OK   000003_create_log_entries_table.sql (18.04ms)
2025/08/01 15:40:31 OK   000004_create_tags_system.sql (14.33ms)
2025/08/01 15:40:31 OK   000005_create_auth_tables.sql (17.07ms)
2025/08/01 15:40:31 OK   000006_create_insights_table.sql (20.56ms)
2025/08/01 15:40:31 OK   000007_create_performance_indexes.sql (12.9ms)
2025/08/01 15:40:31 OK   000008_create_analytics_views.sql (23.34ms)
2025/08/01 15:40:31 OK   000009_development_data.sql (28.13ms)
2025/08/01 15:40:31 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36286/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36286/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:31 🐳 Stopping container: a13198c41fe8
2025/08/01 15:40:31 ✅ Container stopped: a13198c41fe8
2025/08/01 15:40:31 🐳 Terminating container: a13198c41fe8
2025/08/01 15:40:32 🚫 Container terminated: a13198c41fe8
=== RUN   TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_no_activities
2025/08/01 15:40:32 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:32 ✅ Container created: 091ed60f104b
2025/08/01 15:40:32 🐳 Starting container: 091ed60f104b
2025/08/01 15:40:32 ✅ Container started: 091ed60f104b
2025/08/01 15:40:32 ⏳ Waiting for container id 091ed60f104b image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003e63238 Strategies:[0xc001f534a0 0xc004391350]}
2025/08/01 15:40:34 🔔 Container is ready: 091ed60f104b
DSN: postgres://testuser:testpass@localhost:36287/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36287/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:34 OK   000001_create_users_table.sql (12.08ms)
2025/08/01 15:40:34 OK   000002_create_projects_table.sql (14.11ms)
2025/08/01 15:40:34 OK   000003_create_log_entries_table.sql (21.18ms)
2025/08/01 15:40:34 OK   000004_create_tags_system.sql (17.51ms)
2025/08/01 15:40:34 OK   000005_create_auth_tables.sql (21.53ms)
2025/08/01 15:40:34 OK   000006_create_insights_table.sql (25.72ms)
2025/08/01 15:40:34 OK   000007_create_performance_indexes.sql (15.93ms)
2025/08/01 15:40:34 OK   000008_create_analytics_views.sql (28.91ms)
2025/08/01 15:40:34 OK   000009_development_data.sql (29.09ms)
2025/08/01 15:40:34 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36287/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36287/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:34 🐳 Stopping container: 091ed60f104b
2025/08/01 15:40:34 ✅ Container stopped: 091ed60f104b
2025/08/01 15:40:34 🐳 Terminating container: 091ed60f104b
2025/08/01 15:40:35 🚫 Container terminated: 091ed60f104b
--- PASS: TestAnalyticsHandler_GetActivitySummary_Comprehensive (8.39s)
    --- PASS: TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_multiple_activities (2.80s)
    --- PASS: TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_date_range_filter (2.64s)
    --- PASS: TestAnalyticsHandler_GetActivitySummary_Comprehensive/activity_summary_with_no_activities (2.94s)
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive/various_invalid_date_formats_in_productivity_metrics
2025/08/01 15:40:35 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:35 ✅ Container created: ccdf9bb95991
2025/08/01 15:40:35 🐳 Starting container: ccdf9bb95991
2025/08/01 15:40:35 ✅ Container started: ccdf9bb95991
2025/08/01 15:40:35 ⏳ Waiting for container id ccdf9bb95991 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002481088 Strategies:[0xc001f0e480 0xc00481fd40]}
2025/08/01 15:40:37 🔔 Container is ready: ccdf9bb95991
DSN: postgres://testuser:testpass@localhost:36288/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36288/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:37 OK   000001_create_users_table.sql (10.56ms)
2025/08/01 15:40:37 OK   000002_create_projects_table.sql (12.14ms)
2025/08/01 15:40:37 OK   000003_create_log_entries_table.sql (18.98ms)
2025/08/01 15:40:37 OK   000004_create_tags_system.sql (17.18ms)
2025/08/01 15:40:37 OK   000005_create_auth_tables.sql (21.61ms)
2025/08/01 15:40:37 OK   000006_create_insights_table.sql (25.08ms)
2025/08/01 15:40:37 OK   000007_create_performance_indexes.sql (15.72ms)
2025/08/01 15:40:37 OK   000008_create_analytics_views.sql (28.68ms)
2025/08/01 15:40:37 OK   000009_development_data.sql (28.94ms)
2025/08/01 15:40:37 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36288/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36288/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:37 🐳 Stopping container: ccdf9bb95991
2025/08/01 15:40:37 ✅ Container stopped: ccdf9bb95991
2025/08/01 15:40:37 🐳 Terminating container: ccdf9bb95991
2025/08/01 15:40:37 🚫 Container terminated: ccdf9bb95991
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive/missing_day_in_date_format
2025/08/01 15:40:37 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:37 ✅ Container created: 3e235c0ae65c
2025/08/01 15:40:37 🐳 Starting container: 3e235c0ae65c
2025/08/01 15:40:37 ✅ Container started: 3e235c0ae65c
2025/08/01 15:40:37 ⏳ Waiting for container id 3e235c0ae65c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001f08378 Strategies:[0xc000ec0720 0xc0012d3680]}
2025/08/01 15:40:39 🔔 Container is ready: 3e235c0ae65c
DSN: postgres://testuser:testpass@localhost:36289/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36289/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:39 OK   000001_create_users_table.sql (11.25ms)
2025/08/01 15:40:39 OK   000002_create_projects_table.sql (12.73ms)
2025/08/01 15:40:39 OK   000003_create_log_entries_table.sql (20.92ms)
2025/08/01 15:40:39 OK   000004_create_tags_system.sql (17.06ms)
2025/08/01 15:40:39 OK   000005_create_auth_tables.sql (19.69ms)
2025/08/01 15:40:39 OK   000006_create_insights_table.sql (23.83ms)
2025/08/01 15:40:39 OK   000007_create_performance_indexes.sql (15.17ms)
2025/08/01 15:40:39 OK   000008_create_analytics_views.sql (28.74ms)
2025/08/01 15:40:39 OK   000009_development_data.sql (28.74ms)
2025/08/01 15:40:39 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36289/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36289/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:40 🐳 Stopping container: 3e235c0ae65c
2025/08/01 15:40:40 ✅ Container stopped: 3e235c0ae65c
2025/08/01 15:40:40 🐳 Terminating container: 3e235c0ae65c
2025/08/01 15:40:40 🚫 Container terminated: 3e235c0ae65c
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive/invalid_month_in_date
2025/08/01 15:40:40 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:40 ✅ Container created: e08c1e6562df
2025/08/01 15:40:40 🐳 Starting container: e08c1e6562df
2025/08/01 15:40:40 ✅ Container started: e08c1e6562df
2025/08/01 15:40:40 ⏳ Waiting for container id e08c1e6562df image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00250a968 Strategies:[0xc001200a80 0xc002529620]}
2025/08/01 15:40:42 🔔 Container is ready: e08c1e6562df
DSN: postgres://testuser:testpass@localhost:36290/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36290/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:42 OK   000001_create_users_table.sql (10.79ms)
2025/08/01 15:40:42 OK   000002_create_projects_table.sql (13.66ms)
2025/08/01 15:40:42 OK   000003_create_log_entries_table.sql (20.06ms)
2025/08/01 15:40:42 OK   000004_create_tags_system.sql (16.49ms)
2025/08/01 15:40:42 OK   000005_create_auth_tables.sql (21.06ms)
2025/08/01 15:40:42 OK   000006_create_insights_table.sql (24.87ms)
2025/08/01 15:40:42 OK   000007_create_performance_indexes.sql (15.78ms)
2025/08/01 15:40:42 OK   000008_create_analytics_views.sql (28.64ms)
2025/08/01 15:40:42 OK   000009_development_data.sql (29.19ms)
2025/08/01 15:40:42 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36290/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36290/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:43 🐳 Stopping container: e08c1e6562df
2025/08/01 15:40:43 ✅ Container stopped: e08c1e6562df
2025/08/01 15:40:43 🐳 Terminating container: e08c1e6562df
2025/08/01 15:40:43 🚫 Container terminated: e08c1e6562df
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive/invalid_day_in_date
2025/08/01 15:40:43 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:43 ✅ Container created: 94e215cc950c
2025/08/01 15:40:43 🐳 Starting container: 94e215cc950c
2025/08/01 15:40:43 ✅ Container started: 94e215cc950c
2025/08/01 15:40:43 ⏳ Waiting for container id 94e215cc950c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0047fab20 Strategies:[0xc001e9c540 0xc0012ae990]}
2025/08/01 15:40:45 🔔 Container is ready: 94e215cc950c
DSN: postgres://testuser:testpass@localhost:36291/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36291/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:45 OK   000001_create_users_table.sql (10.81ms)
2025/08/01 15:40:45 OK   000002_create_projects_table.sql (12.82ms)
2025/08/01 15:40:45 OK   000003_create_log_entries_table.sql (20.5ms)
2025/08/01 15:40:45 OK   000004_create_tags_system.sql (17.28ms)
2025/08/01 15:40:45 OK   000005_create_auth_tables.sql (21.11ms)
2025/08/01 15:40:45 OK   000006_create_insights_table.sql (25.62ms)
2025/08/01 15:40:45 OK   000007_create_performance_indexes.sql (16.67ms)
2025/08/01 15:40:45 OK   000008_create_analytics_views.sql (29.32ms)
2025/08/01 15:40:45 OK   000009_development_data.sql (28.21ms)
2025/08/01 15:40:45 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36291/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36291/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:46 🐳 Stopping container: 94e215cc950c
2025/08/01 15:40:46 ✅ Container stopped: 94e215cc950c
2025/08/01 15:40:46 🐳 Terminating container: 94e215cc950c
2025/08/01 15:40:46 🚫 Container terminated: 94e215cc950c
=== RUN   TestAnalyticsHandler_DateParsing_Comprehensive/malformed_date_with_text
2025/08/01 15:40:46 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:46 ✅ Container created: 50078daa638e
2025/08/01 15:40:46 🐳 Starting container: 50078daa638e
2025/08/01 15:40:46 ✅ Container started: 50078daa638e
2025/08/01 15:40:46 ⏳ Waiting for container id 50078daa638e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003b4ed70 Strategies:[0xc000ec0600 0xc0039825d0]}
2025/08/01 15:40:48 🔔 Container is ready: 50078daa638e
DSN: postgres://testuser:testpass@localhost:36292/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36292/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:48 OK   000001_create_users_table.sql (10.4ms)
2025/08/01 15:40:48 OK   000002_create_projects_table.sql (12.44ms)
2025/08/01 15:40:48 OK   000003_create_log_entries_table.sql (18.93ms)
2025/08/01 15:40:48 OK   000004_create_tags_system.sql (16.33ms)
2025/08/01 15:40:48 OK   000005_create_auth_tables.sql (22.05ms)
2025/08/01 15:40:48 OK   000006_create_insights_table.sql (25.94ms)
2025/08/01 15:40:48 OK   000007_create_performance_indexes.sql (15.21ms)
2025/08/01 15:40:48 OK   000008_create_analytics_views.sql (27.84ms)
2025/08/01 15:40:48 OK   000009_development_data.sql (28.61ms)
2025/08/01 15:40:48 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36292/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36292/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:48 🐳 Stopping container: 50078daa638e
2025/08/01 15:40:49 ✅ Container stopped: 50078daa638e
2025/08/01 15:40:49 🐳 Terminating container: 50078daa638e
2025/08/01 15:40:49 🚫 Container terminated: 50078daa638e
--- PASS: TestAnalyticsHandler_DateParsing_Comprehensive (14.15s)
    --- PASS: TestAnalyticsHandler_DateParsing_Comprehensive/various_invalid_date_formats_in_productivity_metrics (2.76s)
    --- PASS: TestAnalyticsHandler_DateParsing_Comprehensive/missing_day_in_date_format (2.77s)
    --- PASS: TestAnalyticsHandler_DateParsing_Comprehensive/invalid_month_in_date (2.93s)
    --- PASS: TestAnalyticsHandler_DateParsing_Comprehensive/invalid_day_in_date (2.86s)
    --- PASS: TestAnalyticsHandler_DateParsing_Comprehensive/malformed_date_with_text (2.83s)
=== RUN   TestAnalyticsHandler_GetProductivityMetrics
=== RUN   TestAnalyticsHandler_GetProductivityMetrics/successful_productivity_metrics_retrieval
2025/08/01 15:40:49 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:49 ✅ Container created: 5a8e92220946
2025/08/01 15:40:49 🐳 Starting container: 5a8e92220946
2025/08/01 15:40:49 ✅ Container started: 5a8e92220946
2025/08/01 15:40:49 ⏳ Waiting for container id 5a8e92220946 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0024803c8 Strategies:[0xc0001e3920 0xc003d5ef00]}
2025/08/01 15:40:51 🔔 Container is ready: 5a8e92220946
DSN: postgres://testuser:testpass@localhost:36293/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36293/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:51 OK   000001_create_users_table.sql (11.18ms)
2025/08/01 15:40:51 OK   000002_create_projects_table.sql (12.75ms)
2025/08/01 15:40:51 OK   000003_create_log_entries_table.sql (20.32ms)
2025/08/01 15:40:51 OK   000004_create_tags_system.sql (16.84ms)
2025/08/01 15:40:51 OK   000005_create_auth_tables.sql (20.74ms)
2025/08/01 15:40:51 OK   000006_create_insights_table.sql (25.19ms)
2025/08/01 15:40:51 OK   000007_create_performance_indexes.sql (15.11ms)
2025/08/01 15:40:51 OK   000008_create_analytics_views.sql (27.45ms)
2025/08/01 15:40:51 OK   000009_development_data.sql (28.27ms)
2025/08/01 15:40:51 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36293/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36293/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:51 🐳 Stopping container: 5a8e92220946
2025/08/01 15:40:51 ✅ Container stopped: 5a8e92220946
2025/08/01 15:40:51 🐳 Terminating container: 5a8e92220946
2025/08/01 15:40:51 🚫 Container terminated: 5a8e92220946
=== RUN   TestAnalyticsHandler_GetProductivityMetrics/productivity_metrics_with_date_range
2025/08/01 15:40:51 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:51 ✅ Container created: 15d178b5fa03
2025/08/01 15:40:51 🐳 Starting container: 15d178b5fa03
2025/08/01 15:40:52 ✅ Container started: 15d178b5fa03
2025/08/01 15:40:52 ⏳ Waiting for container id 15d178b5fa03 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001eb1fd0 Strategies:[0xc000709140 0xc000464f00]}
2025/08/01 15:40:53 🔔 Container is ready: 15d178b5fa03
DSN: postgres://testuser:testpass@localhost:36294/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36294/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:53 OK   000001_create_users_table.sql (10.72ms)
2025/08/01 15:40:53 OK   000002_create_projects_table.sql (12.32ms)
2025/08/01 15:40:54 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:40:54 OK   000004_create_tags_system.sql (16.74ms)
2025/08/01 15:40:54 OK   000005_create_auth_tables.sql (20.92ms)
2025/08/01 15:40:54 OK   000006_create_insights_table.sql (24.81ms)
2025/08/01 15:40:54 OK   000007_create_performance_indexes.sql (15.7ms)
2025/08/01 15:40:54 OK   000008_create_analytics_views.sql (27.79ms)
2025/08/01 15:40:54 OK   000009_development_data.sql (28.14ms)
2025/08/01 15:40:54 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36294/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36294/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:54 🐳 Stopping container: 15d178b5fa03
2025/08/01 15:40:54 ✅ Container stopped: 15d178b5fa03
2025/08/01 15:40:54 🐳 Terminating container: 15d178b5fa03
2025/08/01 15:40:54 🚫 Container terminated: 15d178b5fa03
--- PASS: TestAnalyticsHandler_GetProductivityMetrics (5.55s)
    --- PASS: TestAnalyticsHandler_GetProductivityMetrics/successful_productivity_metrics_retrieval (2.75s)
    --- PASS: TestAnalyticsHandler_GetProductivityMetrics/productivity_metrics_with_date_range (2.79s)
=== RUN   TestAnalyticsHandler_GetActivitySummary
=== RUN   TestAnalyticsHandler_GetActivitySummary/successful_activity_summary_retrieval
2025/08/01 15:40:54 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:54 ✅ Container created: 8ec6a5c00bb5
2025/08/01 15:40:54 🐳 Starting container: 8ec6a5c00bb5
2025/08/01 15:40:54 ✅ Container started: 8ec6a5c00bb5
2025/08/01 15:40:54 ⏳ Waiting for container id 8ec6a5c00bb5 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003e62bd0 Strategies:[0xc000ec0900 0xc003ea77a0]}
2025/08/01 15:40:56 🔔 Container is ready: 8ec6a5c00bb5
DSN: postgres://testuser:testpass@localhost:36295/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36295/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:56 OK   000001_create_users_table.sql (11.02ms)
2025/08/01 15:40:56 OK   000002_create_projects_table.sql (12.89ms)
2025/08/01 15:40:56 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:40:56 OK   000004_create_tags_system.sql (16.73ms)
2025/08/01 15:40:56 OK   000005_create_auth_tables.sql (20.24ms)
2025/08/01 15:40:56 OK   000006_create_insights_table.sql (24.17ms)
2025/08/01 15:40:56 OK   000007_create_performance_indexes.sql (15.37ms)
2025/08/01 15:40:56 OK   000008_create_analytics_views.sql (27.94ms)
2025/08/01 15:40:56 OK   000009_development_data.sql (28.59ms)
2025/08/01 15:40:56 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36295/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36295/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:40:57 🐳 Stopping container: 8ec6a5c00bb5
2025/08/01 15:40:57 ✅ Container stopped: 8ec6a5c00bb5
2025/08/01 15:40:57 🐳 Terminating container: 8ec6a5c00bb5
2025/08/01 15:40:57 🚫 Container terminated: 8ec6a5c00bb5
=== RUN   TestAnalyticsHandler_GetActivitySummary/activity_summary_with_date_range
2025/08/01 15:40:57 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:40:57 ✅ Container created: 13cdf31e35d8
2025/08/01 15:40:57 🐳 Starting container: 13cdf31e35d8
2025/08/01 15:40:57 ✅ Container started: 13cdf31e35d8
2025/08/01 15:40:57 ⏳ Waiting for container id 13cdf31e35d8 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00438bc30 Strategies:[0xc001188600 0xc003adbd10]}
2025/08/01 15:40:59 🔔 Container is ready: 13cdf31e35d8
DSN: postgres://testuser:testpass@localhost:36296/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36296/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:40:59 OK   000001_create_users_table.sql (10.04ms)
2025/08/01 15:40:59 OK   000002_create_projects_table.sql (11.93ms)
2025/08/01 15:40:59 OK   000003_create_log_entries_table.sql (19.02ms)
2025/08/01 15:40:59 OK   000004_create_tags_system.sql (16.2ms)
2025/08/01 15:40:59 OK   000005_create_auth_tables.sql (20.4ms)
2025/08/01 15:40:59 OK   000006_create_insights_table.sql (23.99ms)
2025/08/01 15:40:59 OK   000007_create_performance_indexes.sql (15.25ms)
2025/08/01 15:40:59 OK   000008_create_analytics_views.sql (30.03ms)
2025/08/01 15:40:59 OK   000009_development_data.sql (29.08ms)
2025/08/01 15:40:59 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36296/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36296/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:00 🐳 Stopping container: 13cdf31e35d8
2025/08/01 15:41:00 ✅ Container stopped: 13cdf31e35d8
2025/08/01 15:41:00 🐳 Terminating container: 13cdf31e35d8
2025/08/01 15:41:00 🚫 Container terminated: 13cdf31e35d8
--- PASS: TestAnalyticsHandler_GetActivitySummary (5.61s)
    --- PASS: TestAnalyticsHandler_GetActivitySummary/successful_activity_summary_retrieval (2.84s)
    --- PASS: TestAnalyticsHandler_GetActivitySummary/activity_summary_with_date_range (2.77s)
=== RUN   TestAnalyticsHandler_ErrorHandling
=== RUN   TestAnalyticsHandler_ErrorHandling/unauthorized_productivity_metrics_request
2025/08/01 15:41:00 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:00 ✅ Container created: 954adf3e27f7
2025/08/01 15:41:00 🐳 Starting container: 954adf3e27f7
2025/08/01 15:41:00 ✅ Container started: 954adf3e27f7
2025/08/01 15:41:00 ⏳ Waiting for container id 954adf3e27f7 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002480cb0 Strategies:[0xc000301f20 0xc0019e0e10]}
2025/08/01 15:41:02 🔔 Container is ready: 954adf3e27f7
DSN: postgres://testuser:testpass@localhost:36297/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36297/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:02 OK   000001_create_users_table.sql (9.73ms)
2025/08/01 15:41:02 OK   000002_create_projects_table.sql (11.94ms)
2025/08/01 15:41:02 OK   000003_create_log_entries_table.sql (19.66ms)
2025/08/01 15:41:02 OK   000004_create_tags_system.sql (17.14ms)
2025/08/01 15:41:02 OK   000005_create_auth_tables.sql (20.44ms)
2025/08/01 15:41:02 OK   000006_create_insights_table.sql (23.75ms)
2025/08/01 15:41:02 OK   000007_create_performance_indexes.sql (15.59ms)
2025/08/01 15:41:02 OK   000008_create_analytics_views.sql (28.42ms)
2025/08/01 15:41:02 OK   000009_development_data.sql (25.24ms)
2025/08/01 15:41:02 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36297/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36297/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:02 🐳 Stopping container: 954adf3e27f7
2025/08/01 15:41:02 ✅ Container stopped: 954adf3e27f7
2025/08/01 15:41:02 🐳 Terminating container: 954adf3e27f7
2025/08/01 15:41:02 🚫 Container terminated: 954adf3e27f7
=== RUN   TestAnalyticsHandler_ErrorHandling/unauthorized_activity_summary_request
2025/08/01 15:41:02 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:02 ✅ Container created: 057b81dd2030
2025/08/01 15:41:02 🐳 Starting container: 057b81dd2030
2025/08/01 15:41:03 ✅ Container started: 057b81dd2030
2025/08/01 15:41:03 ⏳ Waiting for container id 057b81dd2030 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001e74210 Strategies:[0xc001200480 0xc003868e70]}
2025/08/01 15:41:04 🔔 Container is ready: 057b81dd2030
DSN: postgres://testuser:testpass@localhost:36298/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36298/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:04 OK   000001_create_users_table.sql (10.58ms)
2025/08/01 15:41:04 OK   000002_create_projects_table.sql (12.51ms)
2025/08/01 15:41:04 OK   000003_create_log_entries_table.sql (19.97ms)
2025/08/01 15:41:04 OK   000004_create_tags_system.sql (16.89ms)
2025/08/01 15:41:04 OK   000005_create_auth_tables.sql (20.76ms)
2025/08/01 15:41:04 OK   000006_create_insights_table.sql (24.83ms)
2025/08/01 15:41:05 OK   000007_create_performance_indexes.sql (15.59ms)
2025/08/01 15:41:05 OK   000008_create_analytics_views.sql (28.07ms)
2025/08/01 15:41:05 OK   000009_development_data.sql (29.02ms)
2025/08/01 15:41:05 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36298/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36298/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:05 🐳 Stopping container: 057b81dd2030
2025/08/01 15:41:05 ✅ Container stopped: 057b81dd2030
2025/08/01 15:41:05 🐳 Terminating container: 057b81dd2030
2025/08/01 15:41:05 🚫 Container terminated: 057b81dd2030
=== RUN   TestAnalyticsHandler_ErrorHandling/invalid_date_format_in_productivity_metrics
2025/08/01 15:41:05 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:05 ✅ Container created: f329b3385ec0
2025/08/01 15:41:05 🐳 Starting container: f329b3385ec0
2025/08/01 15:41:05 ✅ Container started: f329b3385ec0
2025/08/01 15:41:05 ⏳ Waiting for container id f329b3385ec0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001f41bf8 Strategies:[0xc000519800 0xc003f20900]}
2025/08/01 15:41:07 🔔 Container is ready: f329b3385ec0
DSN: postgres://testuser:testpass@localhost:36299/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36299/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:07 OK   000001_create_users_table.sql (11.83ms)
2025/08/01 15:41:07 OK   000002_create_projects_table.sql (13.14ms)
2025/08/01 15:41:07 OK   000003_create_log_entries_table.sql (21.74ms)
2025/08/01 15:41:07 OK   000004_create_tags_system.sql (17.4ms)
2025/08/01 15:41:07 OK   000005_create_auth_tables.sql (21.3ms)
2025/08/01 15:41:07 OK   000006_create_insights_table.sql (25.77ms)
2025/08/01 15:41:07 OK   000007_create_performance_indexes.sql (15.88ms)
2025/08/01 15:41:07 OK   000008_create_analytics_views.sql (29.27ms)
2025/08/01 15:41:07 OK   000009_development_data.sql (29.48ms)
2025/08/01 15:41:07 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36299/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36299/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:07 🐳 Stopping container: f329b3385ec0
2025/08/01 15:41:08 ✅ Container stopped: f329b3385ec0
2025/08/01 15:41:08 🐳 Terminating container: f329b3385ec0
2025/08/01 15:41:08 🚫 Container terminated: f329b3385ec0
=== RUN   TestAnalyticsHandler_ErrorHandling/invalid_date_format_in_activity_summary
2025/08/01 15:41:08 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:08 ✅ Container created: 746c26a39ad2
2025/08/01 15:41:08 🐳 Starting container: 746c26a39ad2
2025/08/01 15:41:08 ✅ Container started: 746c26a39ad2
2025/08/01 15:41:08 ⏳ Waiting for container id 746c26a39ad2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00067b818 Strategies:[0xc000657020 0xc000fce630]}
2025/08/01 15:41:10 🔔 Container is ready: 746c26a39ad2
DSN: postgres://testuser:testpass@localhost:36300/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36300/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:10 OK   000001_create_users_table.sql (9.91ms)
2025/08/01 15:41:10 OK   000002_create_projects_table.sql (11.52ms)
2025/08/01 15:41:10 OK   000003_create_log_entries_table.sql (19.23ms)
2025/08/01 15:41:10 OK   000004_create_tags_system.sql (15.91ms)
2025/08/01 15:41:10 OK   000005_create_auth_tables.sql (20.29ms)
2025/08/01 15:41:10 OK   000006_create_insights_table.sql (24.16ms)
2025/08/01 15:41:10 OK   000007_create_performance_indexes.sql (15.74ms)
2025/08/01 15:41:10 OK   000008_create_analytics_views.sql (28.03ms)
2025/08/01 15:41:10 OK   000009_development_data.sql (28.89ms)
2025/08/01 15:41:10 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36300/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36300/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:10 🐳 Stopping container: 746c26a39ad2
2025/08/01 15:41:10 ✅ Container stopped: 746c26a39ad2
2025/08/01 15:41:10 🐳 Terminating container: 746c26a39ad2
2025/08/01 15:41:10 🚫 Container terminated: 746c26a39ad2
=== RUN   TestAnalyticsHandler_ErrorHandling/end_date_before_start_date_in_productivity_metrics
2025/08/01 15:41:10 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:11 ✅ Container created: e0ea7a91cc1f
2025/08/01 15:41:11 🐳 Starting container: e0ea7a91cc1f
2025/08/01 15:41:11 ✅ Container started: e0ea7a91cc1f
2025/08/01 15:41:11 ⏳ Waiting for container id e0ea7a91cc1f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00483cf58 Strategies:[0xc000300b40 0xc003953440]}
2025/08/01 15:41:12 🔔 Container is ready: e0ea7a91cc1f
DSN: postgres://testuser:testpass@localhost:36301/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36301/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:13 OK   000001_create_users_table.sql (10.64ms)
2025/08/01 15:41:13 OK   000002_create_projects_table.sql (12.02ms)
2025/08/01 15:41:13 OK   000003_create_log_entries_table.sql (18.86ms)
2025/08/01 15:41:13 OK   000004_create_tags_system.sql (16.35ms)
2025/08/01 15:41:13 OK   000005_create_auth_tables.sql (19.39ms)
2025/08/01 15:41:13 OK   000006_create_insights_table.sql (23.96ms)
2025/08/01 15:41:13 OK   000007_create_performance_indexes.sql (15.37ms)
2025/08/01 15:41:13 OK   000008_create_analytics_views.sql (26.96ms)
2025/08/01 15:41:13 OK   000009_development_data.sql (28.53ms)
2025/08/01 15:41:13 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36301/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36301/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:13 🐳 Stopping container: e0ea7a91cc1f
2025/08/01 15:41:13 ✅ Container stopped: e0ea7a91cc1f
2025/08/01 15:41:13 🐳 Terminating container: e0ea7a91cc1f
2025/08/01 15:41:13 🚫 Container terminated: e0ea7a91cc1f
=== RUN   TestAnalyticsHandler_ErrorHandling/end_date_before_start_date_in_activity_summary
2025/08/01 15:41:13 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:13 ✅ Container created: 8669bd1667b6
2025/08/01 15:41:13 🐳 Starting container: 8669bd1667b6
2025/08/01 15:41:13 ✅ Container started: 8669bd1667b6
2025/08/01 15:41:13 ⏳ Waiting for container id 8669bd1667b6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003f0f188 Strategies:[0xc001f53980 0xc003f20810]}
2025/08/01 15:41:15 🔔 Container is ready: 8669bd1667b6
DSN: postgres://testuser:testpass@localhost:36302/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36302/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:15 OK   000001_create_users_table.sql (10.64ms)
2025/08/01 15:41:15 OK   000002_create_projects_table.sql (12.94ms)
2025/08/01 15:41:15 OK   000003_create_log_entries_table.sql (19.67ms)
2025/08/01 15:41:15 OK   000004_create_tags_system.sql (16.11ms)
2025/08/01 15:41:15 OK   000005_create_auth_tables.sql (19.91ms)
2025/08/01 15:41:15 OK   000006_create_insights_table.sql (24.09ms)
2025/08/01 15:41:15 OK   000007_create_performance_indexes.sql (15.42ms)
2025/08/01 15:41:15 OK   000008_create_analytics_views.sql (27.92ms)
2025/08/01 15:41:15 OK   000009_development_data.sql (28.41ms)
2025/08/01 15:41:15 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36302/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36302/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /health                   --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).HealthCheck-fm (8 handlers)
[GIN-debug] GET    /ready                    --> github.com/garnizeh/englog/internal/handlers.(*HealthHandler).ReadinessCheck-fm (8 handlers)
[GIN-debug] POST   /v1/auth/register         --> github.com/garnizeh/englog/internal/auth.(*AuthService).RegisterHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/login            --> github.com/garnizeh/englog/internal/auth.(*AuthService).LoginHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/refresh          --> github.com/garnizeh/englog/internal/auth.(*AuthService).RefreshHandler-fm (8 handlers)
[GIN-debug] POST   /v1/auth/logout           --> github.com/garnizeh/englog/internal/auth.(*AuthService).LogoutHandler-fm (8 handlers)
[GIN-debug] GET    /v1/auth/me               --> github.com/garnizeh/englog/internal/auth.(*AuthService).MeHandler-fm (9 handlers)
[GIN-debug] POST   /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).CreateLogEntry-fm (9 handlers)
[GIN-debug] GET    /v1/logs                  --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntries-fm (9 handlers)
[GIN-debug] GET    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).GetLogEntry-fm (10 handlers)
[GIN-debug] PUT    /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).UpdateLogEntry-fm (10 handlers)
[GIN-debug] DELETE /v1/logs/:id              --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).DeleteLogEntry-fm (10 handlers)
[GIN-debug] POST   /v1/logs/bulk             --> github.com/garnizeh/englog/internal/handlers.(*LogEntryHandler).BulkCreateLogEntries-fm (9 handlers)
[GIN-debug] POST   /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).CreateProject-fm (9 handlers)
[GIN-debug] GET    /v1/projects              --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProjects-fm (9 handlers)
[GIN-debug] GET    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).GetProject-fm (10 handlers)
[GIN-debug] PUT    /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).UpdateProject-fm (10 handlers)
[GIN-debug] DELETE /v1/projects/:id          --> github.com/garnizeh/englog/internal/handlers.(*ProjectHandler).DeleteProject-fm (10 handlers)
[GIN-debug] GET    /v1/analytics/productivity --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetProductivityMetrics-fm (9 handlers)
[GIN-debug] GET    /v1/analytics/summary     --> github.com/garnizeh/englog/internal/handlers.(*AnalyticsHandler).GetActivitySummary-fm (9 handlers)
[GIN-debug] POST   /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).CreateTag-fm (9 handlers)
[GIN-debug] GET    /v1/tags                  --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/popular          --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetPopularTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/recent           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetRecentlyUsedTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/search           --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).SearchTags-fm (9 handlers)
[GIN-debug] GET    /v1/tags/usage            --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetUserTagUsage-fm (9 handlers)
[GIN-debug] GET    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).GetTag-fm (9 handlers)
[GIN-debug] PUT    /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).UpdateTag-fm (10 handlers)
[GIN-debug] DELETE /v1/tags/:id              --> github.com/garnizeh/englog/internal/handlers.(*TagHandler).DeleteTag-fm (9 handlers)
[GIN-debug] GET    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).GetProfile-fm (9 handlers)
[GIN-debug] PUT    /v1/users/profile         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).UpdateProfile-fm (9 handlers)
[GIN-debug] POST   /v1/users/change-password --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).ChangePassword-fm (9 handlers)
[GIN-debug] DELETE /v1/users/account         --> github.com/garnizeh/englog/internal/handlers.(*UserHandler).DeleteAccount-fm (9 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/swaggo/gin-swagger.CustomWrapHandler.func1 (8 handlers)
2025/08/01 15:41:16 🐳 Stopping container: 8669bd1667b6
2025/08/01 15:41:16 ✅ Container stopped: 8669bd1667b6
2025/08/01 15:41:16 🐳 Terminating container: 8669bd1667b6
2025/08/01 15:41:16 🚫 Container terminated: 8669bd1667b6
--- PASS: TestAnalyticsHandler_ErrorHandling (16.19s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/unauthorized_productivity_metrics_request (2.52s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/unauthorized_activity_summary_request (2.54s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/invalid_date_format_in_productivity_metrics (2.79s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/invalid_date_format_in_activity_summary (2.80s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/end_date_before_start_date_in_productivity_metrics (2.74s)
    --- PASS: TestAnalyticsHandler_ErrorHandling/end_date_before_start_date_in_activity_summary (2.80s)
=== RUN   TestHealthHandler_HealthCheck
=== RUN   TestHealthHandler_HealthCheck/successful_health_check
=== RUN   TestHealthHandler_HealthCheck/uptime_progression
=== RUN   TestHealthHandler_HealthCheck/consistent_version
--- PASS: TestHealthHandler_HealthCheck (0.01s)
    --- PASS: TestHealthHandler_HealthCheck/successful_health_check (0.00s)
    --- PASS: TestHealthHandler_HealthCheck/uptime_progression (0.01s)
    --- PASS: TestHealthHandler_HealthCheck/consistent_version (0.00s)
=== RUN   TestHealthHandler_ReadinessCheck
=== RUN   TestHealthHandler_ReadinessCheck/successful_readiness_check
=== RUN   TestHealthHandler_ReadinessCheck/multiple_readiness_checks
=== RUN   TestHealthHandler_ReadinessCheck/timestamp_format_validation
--- PASS: TestHealthHandler_ReadinessCheck (0.00s)
    --- PASS: TestHealthHandler_ReadinessCheck/successful_readiness_check (0.00s)
    --- PASS: TestHealthHandler_ReadinessCheck/multiple_readiness_checks (0.00s)
    --- PASS: TestHealthHandler_ReadinessCheck/timestamp_format_validation (0.00s)
=== RUN   TestHealthHandler_ErrorHandling
=== RUN   TestHealthHandler_ErrorHandling/wrong_http_method_health
=== RUN   TestHealthHandler_ErrorHandling/wrong_http_method_readiness
=== RUN   TestHealthHandler_ErrorHandling/non_existent_endpoint
--- PASS: TestHealthHandler_ErrorHandling (0.00s)
    --- PASS: TestHealthHandler_ErrorHandling/wrong_http_method_health (0.00s)
    --- PASS: TestHealthHandler_ErrorHandling/wrong_http_method_readiness (0.00s)
    --- PASS: TestHealthHandler_ErrorHandling/non_existent_endpoint (0.00s)
=== RUN   TestHealthHandler_ResponseHeaders
=== RUN   TestHealthHandler_ResponseHeaders/health_check_headers
=== RUN   TestHealthHandler_ResponseHeaders/readiness_check_headers
--- PASS: TestHealthHandler_ResponseHeaders (0.00s)
    --- PASS: TestHealthHandler_ResponseHeaders/health_check_headers (0.00s)
    --- PASS: TestHealthHandler_ResponseHeaders/readiness_check_headers (0.00s)
=== RUN   TestLogEntryHandler_CreateLogEntry
=== RUN   TestLogEntryHandler_CreateLogEntry/successful_log_entry_creation
2025/08/01 15:41:16 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:16 ✅ Container created: 90eee67d26da
2025/08/01 15:41:16 🐳 Starting container: 90eee67d26da
2025/08/01 15:41:16 ✅ Container started: 90eee67d26da
2025/08/01 15:41:16 ⏳ Waiting for container id 90eee67d26da image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0025813f8 Strategies:[0xc000657320 0xc0019e0300]}
2025/08/01 15:41:18 🔔 Container is ready: 90eee67d26da
DSN: postgres://testuser:testpass@localhost:36303/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36303/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:18 OK   000001_create_users_table.sql (10.77ms)
2025/08/01 15:41:18 OK   000002_create_projects_table.sql (12.34ms)
2025/08/01 15:41:18 OK   000003_create_log_entries_table.sql (19.97ms)
2025/08/01 15:41:18 OK   000004_create_tags_system.sql (17.27ms)
2025/08/01 15:41:18 OK   000005_create_auth_tables.sql (21.12ms)
2025/08/01 15:41:18 OK   000006_create_insights_table.sql (25.12ms)
2025/08/01 15:41:18 OK   000007_create_performance_indexes.sql (15.64ms)
2025/08/01 15:41:18 OK   000008_create_analytics_views.sql (28.36ms)
2025/08/01 15:41:18 OK   000009_development_data.sql (29.01ms)
2025/08/01 15:41:18 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36303/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36303/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:19 🐳 Stopping container: 90eee67d26da
2025/08/01 15:41:19 ✅ Container stopped: 90eee67d26da
2025/08/01 15:41:19 🐳 Terminating container: 90eee67d26da
2025/08/01 15:41:19 🚫 Container terminated: 90eee67d26da
=== RUN   TestLogEntryHandler_CreateLogEntry/unauthorized_access
2025/08/01 15:41:19 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:19 ✅ Container created: bed19aa76665
2025/08/01 15:41:19 🐳 Starting container: bed19aa76665
2025/08/01 15:41:19 ✅ Container started: bed19aa76665
2025/08/01 15:41:19 ⏳ Waiting for container id bed19aa76665 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003cc8c00 Strategies:[0xc003949e00 0xc003cdaae0]}
2025/08/01 15:41:21 🔔 Container is ready: bed19aa76665
DSN: postgres://testuser:testpass@localhost:36304/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36304/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:21 OK   000001_create_users_table.sql (11.96ms)
2025/08/01 15:41:21 OK   000002_create_projects_table.sql (13.52ms)
2025/08/01 15:41:21 OK   000003_create_log_entries_table.sql (21.24ms)
2025/08/01 15:41:21 OK   000004_create_tags_system.sql (17.62ms)
2025/08/01 15:41:21 OK   000005_create_auth_tables.sql (21.4ms)
2025/08/01 15:41:21 OK   000006_create_insights_table.sql (25.23ms)
2025/08/01 15:41:21 OK   000007_create_performance_indexes.sql (17.29ms)
2025/08/01 15:41:21 OK   000008_create_analytics_views.sql (28.08ms)
2025/08/01 15:41:21 OK   000009_development_data.sql (28.82ms)
2025/08/01 15:41:21 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36304/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36304/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:21 🐳 Stopping container: bed19aa76665
2025/08/01 15:41:21 ✅ Container stopped: bed19aa76665
2025/08/01 15:41:21 🐳 Terminating container: bed19aa76665
2025/08/01 15:41:21 🚫 Container terminated: bed19aa76665
=== RUN   TestLogEntryHandler_CreateLogEntry/invalid_request_body
2025/08/01 15:41:21 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:22 ✅ Container created: ca58f3124658
2025/08/01 15:41:22 🐳 Starting container: ca58f3124658
2025/08/01 15:41:22 ✅ Container started: ca58f3124658
2025/08/01 15:41:22 ⏳ Waiting for container id ca58f3124658 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00478c588 Strategies:[0xc003b88660 0xc003869230]}
2025/08/01 15:41:23 🔔 Container is ready: ca58f3124658
DSN: postgres://testuser:testpass@localhost:36305/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36305/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:24 OK   000001_create_users_table.sql (9.31ms)
2025/08/01 15:41:24 OK   000002_create_projects_table.sql (11.5ms)
2025/08/01 15:41:24 OK   000003_create_log_entries_table.sql (19.04ms)
2025/08/01 15:41:24 OK   000004_create_tags_system.sql (16.41ms)
2025/08/01 15:41:24 OK   000005_create_auth_tables.sql (21.22ms)
2025/08/01 15:41:24 OK   000006_create_insights_table.sql (25.89ms)
2025/08/01 15:41:24 OK   000007_create_performance_indexes.sql (16.34ms)
2025/08/01 15:41:24 OK   000008_create_analytics_views.sql (29.4ms)
2025/08/01 15:41:24 OK   000009_development_data.sql (30.79ms)
2025/08/01 15:41:24 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36305/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36305/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:24 🐳 Stopping container: ca58f3124658
2025/08/01 15:41:24 ✅ Container stopped: ca58f3124658
2025/08/01 15:41:24 🐳 Terminating container: ca58f3124658
2025/08/01 15:41:24 🚫 Container terminated: ca58f3124658
=== RUN   TestLogEntryHandler_CreateLogEntry/missing_required_fields
2025/08/01 15:41:24 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:24 ✅ Container created: b347e83e43e4
2025/08/01 15:41:24 🐳 Starting container: b347e83e43e4
2025/08/01 15:41:25 ✅ Container started: b347e83e43e4
2025/08/01 15:41:25 ⏳ Waiting for container id b347e83e43e4 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00478da28 Strategies:[0xc003b89da0 0xc0039839b0]}
2025/08/01 15:41:26 🔔 Container is ready: b347e83e43e4
DSN: postgres://testuser:testpass@localhost:36306/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36306/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:26 OK   000001_create_users_table.sql (10.03ms)
2025/08/01 15:41:26 OK   000002_create_projects_table.sql (11.47ms)
2025/08/01 15:41:26 OK   000003_create_log_entries_table.sql (18.21ms)
2025/08/01 15:41:26 OK   000004_create_tags_system.sql (15.76ms)
2025/08/01 15:41:26 OK   000005_create_auth_tables.sql (19.6ms)
2025/08/01 15:41:26 OK   000006_create_insights_table.sql (23.74ms)
2025/08/01 15:41:26 OK   000007_create_performance_indexes.sql (15.17ms)
2025/08/01 15:41:26 OK   000008_create_analytics_views.sql (26.3ms)
2025/08/01 15:41:27 OK   000009_development_data.sql (25.38ms)
2025/08/01 15:41:27 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36306/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36306/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:27 🐳 Stopping container: b347e83e43e4
2025/08/01 15:41:27 ✅ Container stopped: b347e83e43e4
2025/08/01 15:41:27 🐳 Terminating container: b347e83e43e4
2025/08/01 15:41:27 🚫 Container terminated: b347e83e43e4
--- PASS: TestLogEntryHandler_CreateLogEntry (11.00s)
    --- PASS: TestLogEntryHandler_CreateLogEntry/successful_log_entry_creation (2.94s)
    --- PASS: TestLogEntryHandler_CreateLogEntry/unauthorized_access (2.54s)
    --- PASS: TestLogEntryHandler_CreateLogEntry/invalid_request_body (2.80s)
    --- PASS: TestLogEntryHandler_CreateLogEntry/missing_required_fields (2.72s)
=== RUN   TestLogEntryHandler_GetLogEntry
=== RUN   TestLogEntryHandler_GetLogEntry/successful_log_entry_retrieval
2025/08/01 15:41:27 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:27 ✅ Container created: ef3c497c310d
2025/08/01 15:41:27 🐳 Starting container: ef3c497c310d
2025/08/01 15:41:27 ✅ Container started: ef3c497c310d
2025/08/01 15:41:27 ⏳ Waiting for container id ef3c497c310d image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00445c578 Strategies:[0xc000519c20 0xc003f210b0]}
2025/08/01 15:41:29 🔔 Container is ready: ef3c497c310d
DSN: postgres://testuser:testpass@localhost:36307/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36307/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:29 OK   000001_create_users_table.sql (10.57ms)
2025/08/01 15:41:29 OK   000002_create_projects_table.sql (12.81ms)
2025/08/01 15:41:29 OK   000003_create_log_entries_table.sql (19.14ms)
2025/08/01 15:41:29 OK   000004_create_tags_system.sql (16.06ms)
2025/08/01 15:41:29 OK   000005_create_auth_tables.sql (20.3ms)
2025/08/01 15:41:29 OK   000006_create_insights_table.sql (24.36ms)
2025/08/01 15:41:29 OK   000007_create_performance_indexes.sql (16.03ms)
2025/08/01 15:41:29 OK   000008_create_analytics_views.sql (29.16ms)
2025/08/01 15:41:29 OK   000009_development_data.sql (30.24ms)
2025/08/01 15:41:29 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36307/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36307/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:30 🐳 Stopping container: ef3c497c310d
2025/08/01 15:41:30 ✅ Container stopped: ef3c497c310d
2025/08/01 15:41:30 🐳 Terminating container: ef3c497c310d
2025/08/01 15:41:30 🚫 Container terminated: ef3c497c310d
=== RUN   TestLogEntryHandler_GetLogEntry/unauthorized_access
2025/08/01 15:41:30 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:30 ✅ Container created: 090edaee5e66
2025/08/01 15:41:30 🐳 Starting container: 090edaee5e66
2025/08/01 15:41:30 ✅ Container started: 090edaee5e66
2025/08/01 15:41:30 ⏳ Waiting for container id 090edaee5e66 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0024d8c38 Strategies:[0xc003b88360 0xc004404bd0]}
2025/08/01 15:41:32 🔔 Container is ready: 090edaee5e66
DSN: postgres://testuser:testpass@localhost:36308/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36308/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:32 OK   000001_create_users_table.sql (9.7ms)
2025/08/01 15:41:32 OK   000002_create_projects_table.sql (11.74ms)
2025/08/01 15:41:32 OK   000003_create_log_entries_table.sql (18.93ms)
2025/08/01 15:41:32 OK   000004_create_tags_system.sql (16.91ms)
2025/08/01 15:41:32 OK   000005_create_auth_tables.sql (20.13ms)
2025/08/01 15:41:32 OK   000006_create_insights_table.sql (24.32ms)
2025/08/01 15:41:32 OK   000007_create_performance_indexes.sql (15.68ms)
2025/08/01 15:41:32 OK   000008_create_analytics_views.sql (28.68ms)
2025/08/01 15:41:32 OK   000009_development_data.sql (31.01ms)
2025/08/01 15:41:32 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36308/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36308/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:32 🐳 Stopping container: 090edaee5e66
2025/08/01 15:41:32 ✅ Container stopped: 090edaee5e66
2025/08/01 15:41:32 🐳 Terminating container: 090edaee5e66
2025/08/01 15:41:33 🚫 Container terminated: 090edaee5e66
=== RUN   TestLogEntryHandler_GetLogEntry/log_entry_not_found
2025/08/01 15:41:33 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:33 ✅ Container created: ef1d79ab2830
2025/08/01 15:41:33 🐳 Starting container: ef1d79ab2830
2025/08/01 15:41:33 ✅ Container started: ef1d79ab2830
2025/08/01 15:41:33 ⏳ Waiting for container id ef1d79ab2830 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001e74c78 Strategies:[0xc003d0a9c0 0xc001e6af60]}
2025/08/01 15:41:34 🔔 Container is ready: ef1d79ab2830
DSN: postgres://testuser:testpass@localhost:36309/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36309/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:34 OK   000001_create_users_table.sql (9.79ms)
2025/08/01 15:41:34 OK   000002_create_projects_table.sql (11.53ms)
2025/08/01 15:41:34 OK   000003_create_log_entries_table.sql (18.39ms)
2025/08/01 15:41:35 OK   000004_create_tags_system.sql (16.16ms)
2025/08/01 15:41:35 OK   000005_create_auth_tables.sql (19.46ms)
2025/08/01 15:41:35 OK   000006_create_insights_table.sql (24.02ms)
2025/08/01 15:41:35 OK   000007_create_performance_indexes.sql (15.22ms)
2025/08/01 15:41:35 OK   000008_create_analytics_views.sql (28.24ms)
2025/08/01 15:41:35 OK   000009_development_data.sql (28.79ms)
2025/08/01 15:41:35 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36309/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36309/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:35 🐳 Stopping container: ef1d79ab2830
2025/08/01 15:41:35 ✅ Container stopped: ef1d79ab2830
2025/08/01 15:41:35 🐳 Terminating container: ef1d79ab2830
2025/08/01 15:41:35 🚫 Container terminated: ef1d79ab2830
=== RUN   TestLogEntryHandler_GetLogEntry/missing_log_entry_ID
2025/08/01 15:41:35 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:35 ✅ Container created: fc2e38734714
2025/08/01 15:41:35 🐳 Starting container: fc2e38734714
2025/08/01 15:41:35 ✅ Container started: fc2e38734714
2025/08/01 15:41:35 ⏳ Waiting for container id fc2e38734714 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc000ebf078 Strategies:[0xc00185a4e0 0xc003291110]}
2025/08/01 15:41:37 🔔 Container is ready: fc2e38734714
DSN: postgres://testuser:testpass@localhost:36310/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36310/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:37 OK   000001_create_users_table.sql (10.13ms)
2025/08/01 15:41:37 OK   000002_create_projects_table.sql (11.87ms)
2025/08/01 15:41:37 OK   000003_create_log_entries_table.sql (18.73ms)
2025/08/01 15:41:37 OK   000004_create_tags_system.sql (16.72ms)
2025/08/01 15:41:37 OK   000005_create_auth_tables.sql (21.09ms)
2025/08/01 15:41:37 OK   000006_create_insights_table.sql (24.54ms)
2025/08/01 15:41:37 OK   000007_create_performance_indexes.sql (13.81ms)
2025/08/01 15:41:37 OK   000008_create_analytics_views.sql (24.28ms)
2025/08/01 15:41:37 OK   000009_development_data.sql (24.65ms)
2025/08/01 15:41:37 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36310/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36310/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:38 🐳 Stopping container: fc2e38734714
2025/08/01 15:41:38 ✅ Container stopped: fc2e38734714
2025/08/01 15:41:38 🐳 Terminating container: fc2e38734714
2025/08/01 15:41:38 🚫 Container terminated: fc2e38734714
--- PASS: TestLogEntryHandler_GetLogEntry (10.92s)
    --- PASS: TestLogEntryHandler_GetLogEntry/successful_log_entry_retrieval (2.83s)
    --- PASS: TestLogEntryHandler_GetLogEntry/unauthorized_access (2.66s)
    --- PASS: TestLogEntryHandler_GetLogEntry/log_entry_not_found (2.67s)
    --- PASS: TestLogEntryHandler_GetLogEntry/missing_log_entry_ID (2.76s)
=== RUN   TestLogEntryHandler_GetLogEntries
=== RUN   TestLogEntryHandler_GetLogEntries/successful_log_entries_retrieval
2025/08/01 15:41:38 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:38 ✅ Container created: c0297a4f0c9e
2025/08/01 15:41:38 🐳 Starting container: c0297a4f0c9e
2025/08/01 15:41:38 ✅ Container started: c0297a4f0c9e
2025/08/01 15:41:38 ⏳ Waiting for container id c0297a4f0c9e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc000265838 Strategies:[0xc00003f080 0xc003f11740]}
2025/08/01 15:41:40 🔔 Container is ready: c0297a4f0c9e
DSN: postgres://testuser:testpass@localhost:36311/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36311/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:40 OK   000001_create_users_table.sql (10.96ms)
2025/08/01 15:41:40 OK   000002_create_projects_table.sql (12.52ms)
2025/08/01 15:41:40 OK   000003_create_log_entries_table.sql (19.42ms)
2025/08/01 15:41:40 OK   000004_create_tags_system.sql (16.77ms)
2025/08/01 15:41:40 OK   000005_create_auth_tables.sql (21.01ms)
2025/08/01 15:41:40 OK   000006_create_insights_table.sql (25.47ms)
2025/08/01 15:41:40 OK   000007_create_performance_indexes.sql (16.59ms)
2025/08/01 15:41:40 OK   000008_create_analytics_views.sql (28.64ms)
2025/08/01 15:41:40 OK   000009_development_data.sql (31.36ms)
2025/08/01 15:41:40 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36311/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36311/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:40 🐳 Stopping container: c0297a4f0c9e
2025/08/01 15:41:41 ✅ Container stopped: c0297a4f0c9e
2025/08/01 15:41:41 🐳 Terminating container: c0297a4f0c9e
2025/08/01 15:41:41 🚫 Container terminated: c0297a4f0c9e
=== RUN   TestLogEntryHandler_GetLogEntries/unauthorized_access
2025/08/01 15:41:41 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:41 ✅ Container created: f1b2b93733d3
2025/08/01 15:41:41 🐳 Starting container: f1b2b93733d3
2025/08/01 15:41:41 ✅ Container started: f1b2b93733d3
2025/08/01 15:41:41 ⏳ Waiting for container id f1b2b93733d3 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00478cc38 Strategies:[0xc003cff320 0xc003f3a1b0]}
2025/08/01 15:41:43 🔔 Container is ready: f1b2b93733d3
DSN: postgres://testuser:testpass@localhost:36312/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36312/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:43 OK   000001_create_users_table.sql (11.39ms)
2025/08/01 15:41:43 OK   000002_create_projects_table.sql (12.72ms)
2025/08/01 15:41:43 OK   000003_create_log_entries_table.sql (20.36ms)
2025/08/01 15:41:43 OK   000004_create_tags_system.sql (16.61ms)
2025/08/01 15:41:43 OK   000005_create_auth_tables.sql (20.34ms)
2025/08/01 15:41:43 OK   000006_create_insights_table.sql (24.39ms)
2025/08/01 15:41:43 OK   000007_create_performance_indexes.sql (15.44ms)
2025/08/01 15:41:43 OK   000008_create_analytics_views.sql (28.9ms)
2025/08/01 15:41:43 OK   000009_development_data.sql (29.39ms)
2025/08/01 15:41:43 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36312/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36312/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:43 🐳 Stopping container: f1b2b93733d3
2025/08/01 15:41:43 ✅ Container stopped: f1b2b93733d3
2025/08/01 15:41:43 🐳 Terminating container: f1b2b93733d3
2025/08/01 15:41:43 🚫 Container terminated: f1b2b93733d3
=== RUN   TestLogEntryHandler_GetLogEntries/with_filters
2025/08/01 15:41:43 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:43 ✅ Container created: 55800f5edc57
2025/08/01 15:41:43 🐳 Starting container: 55800f5edc57
2025/08/01 15:41:44 ✅ Container started: 55800f5edc57
2025/08/01 15:41:44 ⏳ Waiting for container id 55800f5edc57 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0022463f8 Strategies:[0xc00200b5c0 0xc002248420]}
2025/08/01 15:41:45 🔔 Container is ready: 55800f5edc57
DSN: postgres://testuser:testpass@localhost:36313/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36313/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:46 OK   000001_create_users_table.sql (10.65ms)
2025/08/01 15:41:46 OK   000002_create_projects_table.sql (12.93ms)
2025/08/01 15:41:46 OK   000003_create_log_entries_table.sql (19.67ms)
2025/08/01 15:41:46 OK   000004_create_tags_system.sql (15.48ms)
2025/08/01 15:41:46 OK   000005_create_auth_tables.sql (17.81ms)
2025/08/01 15:41:46 OK   000006_create_insights_table.sql (23.7ms)
2025/08/01 15:41:46 OK   000007_create_performance_indexes.sql (15.46ms)
2025/08/01 15:41:46 OK   000008_create_analytics_views.sql (27.61ms)
2025/08/01 15:41:46 OK   000009_development_data.sql (25.02ms)
2025/08/01 15:41:46 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36313/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36313/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:46 🐳 Stopping container: 55800f5edc57
2025/08/01 15:41:46 ✅ Container stopped: 55800f5edc57
2025/08/01 15:41:46 🐳 Terminating container: 55800f5edc57
2025/08/01 15:41:46 🚫 Container terminated: 55800f5edc57
=== RUN   TestLogEntryHandler_GetLogEntries/with_invalid_filters
2025/08/01 15:41:46 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:46 ✅ Container created: 35365a515730
2025/08/01 15:41:46 🐳 Starting container: 35365a515730
2025/08/01 15:41:47 ✅ Container started: 35365a515730
2025/08/01 15:41:47 ⏳ Waiting for container id 35365a515730 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003e1cf38 Strategies:[0xc004401860 0xc00389d890]}
2025/08/01 15:41:48 🔔 Container is ready: 35365a515730
DSN: postgres://testuser:testpass@localhost:36314/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36314/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:48 OK   000001_create_users_table.sql (10.45ms)
2025/08/01 15:41:48 OK   000002_create_projects_table.sql (13.68ms)
2025/08/01 15:41:48 OK   000003_create_log_entries_table.sql (20.06ms)
2025/08/01 15:41:48 OK   000004_create_tags_system.sql (17.21ms)
2025/08/01 15:41:48 OK   000005_create_auth_tables.sql (20.37ms)
2025/08/01 15:41:48 OK   000006_create_insights_table.sql (24.83ms)
2025/08/01 15:41:48 OK   000007_create_performance_indexes.sql (15.73ms)
2025/08/01 15:41:49 OK   000008_create_analytics_views.sql (28.51ms)
2025/08/01 15:41:49 OK   000009_development_data.sql (29.7ms)
2025/08/01 15:41:49 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36314/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36314/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:49 🐳 Stopping container: 35365a515730
2025/08/01 15:41:49 ✅ Container stopped: 35365a515730
2025/08/01 15:41:49 🐳 Terminating container: 35365a515730
2025/08/01 15:41:49 🚫 Container terminated: 35365a515730
=== RUN   TestLogEntryHandler_GetLogEntries/with_pagination
2025/08/01 15:41:49 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:49 ✅ Container created: 9e974ad8348f
2025/08/01 15:41:49 🐳 Starting container: 9e974ad8348f
2025/08/01 15:41:49 ✅ Container started: 9e974ad8348f
2025/08/01 15:41:49 ⏳ Waiting for container id 9e974ad8348f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0048b06a8 Strategies:[0xc00101acc0 0xc004404e40]}
2025/08/01 15:41:51 🔔 Container is ready: 9e974ad8348f
DSN: postgres://testuser:testpass@localhost:36315/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36315/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:51 OK   000001_create_users_table.sql (10.89ms)
2025/08/01 15:41:51 OK   000002_create_projects_table.sql (13.32ms)
2025/08/01 15:41:51 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:41:51 OK   000004_create_tags_system.sql (15.49ms)
2025/08/01 15:41:51 OK   000005_create_auth_tables.sql (18.02ms)
2025/08/01 15:41:51 OK   000006_create_insights_table.sql (20.71ms)
2025/08/01 15:41:51 OK   000007_create_performance_indexes.sql (12.91ms)
2025/08/01 15:41:51 OK   000008_create_analytics_views.sql (23.3ms)
2025/08/01 15:41:51 OK   000009_development_data.sql (23.92ms)
2025/08/01 15:41:51 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36315/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36315/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:52 🐳 Stopping container: 9e974ad8348f
2025/08/01 15:41:52 ✅ Container stopped: 9e974ad8348f
2025/08/01 15:41:52 🐳 Terminating container: 9e974ad8348f
2025/08/01 15:41:52 🚫 Container terminated: 9e974ad8348f
--- PASS: TestLogEntryHandler_GetLogEntries (14.00s)
    --- PASS: TestLogEntryHandler_GetLogEntries/successful_log_entries_retrieval (2.83s)
    --- PASS: TestLogEntryHandler_GetLogEntries/unauthorized_access (2.68s)
    --- PASS: TestLogEntryHandler_GetLogEntries/with_filters (2.87s)
    --- PASS: TestLogEntryHandler_GetLogEntries/with_invalid_filters (2.83s)
    --- PASS: TestLogEntryHandler_GetLogEntries/with_pagination (2.80s)
=== RUN   TestLogEntryHandler_UpdateLogEntry
=== RUN   TestLogEntryHandler_UpdateLogEntry/successful_log_entry_update
2025/08/01 15:41:52 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:52 ✅ Container created: ef8d11100851
2025/08/01 15:41:52 🐳 Starting container: ef8d11100851
2025/08/01 15:41:52 ✅ Container started: ef8d11100851
2025/08/01 15:41:52 ⏳ Waiting for container id ef8d11100851 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003a0dd40 Strategies:[0xc0045a9920 0xc003b44e40]}
2025/08/01 15:41:54 🔔 Container is ready: ef8d11100851
DSN: postgres://testuser:testpass@localhost:36316/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36316/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:54 OK   000001_create_users_table.sql (10.6ms)
2025/08/01 15:41:54 OK   000002_create_projects_table.sql (12.5ms)
2025/08/01 15:41:54 OK   000003_create_log_entries_table.sql (20.26ms)
2025/08/01 15:41:54 OK   000004_create_tags_system.sql (17.22ms)
2025/08/01 15:41:54 OK   000005_create_auth_tables.sql (21.2ms)
2025/08/01 15:41:54 OK   000006_create_insights_table.sql (25.48ms)
2025/08/01 15:41:54 OK   000007_create_performance_indexes.sql (16.15ms)
2025/08/01 15:41:54 OK   000008_create_analytics_views.sql (28.97ms)
2025/08/01 15:41:54 OK   000009_development_data.sql (29.4ms)
2025/08/01 15:41:54 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36316/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36316/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:54 🐳 Stopping container: ef8d11100851
2025/08/01 15:41:55 ✅ Container stopped: ef8d11100851
2025/08/01 15:41:55 🐳 Terminating container: ef8d11100851
2025/08/01 15:41:55 🚫 Container terminated: ef8d11100851
=== RUN   TestLogEntryHandler_UpdateLogEntry/unauthorized_access
2025/08/01 15:41:55 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:55 ✅ Container created: 4b09f89d6c2e
2025/08/01 15:41:55 🐳 Starting container: 4b09f89d6c2e
2025/08/01 15:41:55 ✅ Container started: 4b09f89d6c2e
2025/08/01 15:41:55 ⏳ Waiting for container id 4b09f89d6c2e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003e73d60 Strategies:[0xc0044f8060 0xc003df9bf0]}
2025/08/01 15:41:57 🔔 Container is ready: 4b09f89d6c2e
DSN: postgres://testuser:testpass@localhost:36317/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36317/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:57 OK   000001_create_users_table.sql (10.48ms)
2025/08/01 15:41:57 OK   000002_create_projects_table.sql (12.08ms)
2025/08/01 15:41:57 OK   000003_create_log_entries_table.sql (19.1ms)
2025/08/01 15:41:57 OK   000004_create_tags_system.sql (15.99ms)
2025/08/01 15:41:57 OK   000005_create_auth_tables.sql (19.88ms)
2025/08/01 15:41:57 OK   000006_create_insights_table.sql (23.78ms)
2025/08/01 15:41:57 OK   000007_create_performance_indexes.sql (15.06ms)
2025/08/01 15:41:57 OK   000008_create_analytics_views.sql (27.89ms)
2025/08/01 15:41:57 OK   000009_development_data.sql (28.1ms)
2025/08/01 15:41:57 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36317/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36317/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:41:57 🐳 Stopping container: 4b09f89d6c2e
2025/08/01 15:41:57 ✅ Container stopped: 4b09f89d6c2e
2025/08/01 15:41:57 🐳 Terminating container: 4b09f89d6c2e
2025/08/01 15:41:57 🚫 Container terminated: 4b09f89d6c2e
=== RUN   TestLogEntryHandler_UpdateLogEntry/log_entry_not_found
2025/08/01 15:41:57 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:41:57 ✅ Container created: 4c8435d40f81
2025/08/01 15:41:57 🐳 Starting container: 4c8435d40f81
2025/08/01 15:41:57 ✅ Container started: 4c8435d40f81
2025/08/01 15:41:57 ⏳ Waiting for container id 4c8435d40f81 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004458748 Strategies:[0xc0047da540 0xc006a86ed0]}
2025/08/01 15:41:59 🔔 Container is ready: 4c8435d40f81
DSN: postgres://testuser:testpass@localhost:36318/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36318/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:41:59 OK   000001_create_users_table.sql (10.72ms)
2025/08/01 15:41:59 OK   000002_create_projects_table.sql (12.86ms)
2025/08/01 15:41:59 OK   000003_create_log_entries_table.sql (19.99ms)
2025/08/01 15:41:59 OK   000004_create_tags_system.sql (17.56ms)
2025/08/01 15:41:59 OK   000005_create_auth_tables.sql (21.47ms)
2025/08/01 15:41:59 OK   000006_create_insights_table.sql (25.52ms)
2025/08/01 15:41:59 OK   000007_create_performance_indexes.sql (16.13ms)
2025/08/01 15:41:59 OK   000008_create_analytics_views.sql (30.13ms)
2025/08/01 15:41:59 OK   000009_development_data.sql (29.87ms)
2025/08/01 15:41:59 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36318/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36318/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:00 🐳 Stopping container: 4c8435d40f81
2025/08/01 15:42:00 ✅ Container stopped: 4c8435d40f81
2025/08/01 15:42:00 🐳 Terminating container: 4c8435d40f81
2025/08/01 15:42:00 🚫 Container terminated: 4c8435d40f81
=== RUN   TestLogEntryHandler_UpdateLogEntry/invalid_request_body
2025/08/01 15:42:00 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:00 ✅ Container created: f3045d566569
2025/08/01 15:42:00 🐳 Starting container: f3045d566569
2025/08/01 15:42:00 ✅ Container started: f3045d566569
2025/08/01 15:42:00 ⏳ Waiting for container id f3045d566569 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003d9b4b8 Strategies:[0xc006aa5860 0xc003ddd050]}
2025/08/01 15:42:02 🔔 Container is ready: f3045d566569
DSN: postgres://testuser:testpass@localhost:36319/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36319/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:02 OK   000001_create_users_table.sql (10.47ms)
2025/08/01 15:42:02 OK   000002_create_projects_table.sql (12.61ms)
2025/08/01 15:42:02 OK   000003_create_log_entries_table.sql (19.76ms)
2025/08/01 15:42:02 OK   000004_create_tags_system.sql (17.18ms)
2025/08/01 15:42:02 OK   000005_create_auth_tables.sql (21.46ms)
2025/08/01 15:42:02 OK   000006_create_insights_table.sql (25.5ms)
2025/08/01 15:42:02 OK   000007_create_performance_indexes.sql (16.62ms)
2025/08/01 15:42:02 OK   000008_create_analytics_views.sql (30.42ms)
2025/08/01 15:42:02 OK   000009_development_data.sql (29.61ms)
2025/08/01 15:42:02 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36319/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36319/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:03 🐳 Stopping container: f3045d566569
2025/08/01 15:42:03 ✅ Container stopped: f3045d566569
2025/08/01 15:42:03 🐳 Terminating container: f3045d566569
2025/08/01 15:42:03 🚫 Container terminated: f3045d566569
=== RUN   TestLogEntryHandler_UpdateLogEntry/missing_log_entry_ID
2025/08/01 15:42:03 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:03 ✅ Container created: 5201573638cd
2025/08/01 15:42:03 🐳 Starting container: 5201573638cd
2025/08/01 15:42:03 ✅ Container started: 5201573638cd
2025/08/01 15:42:03 ⏳ Waiting for container id 5201573638cd image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0004d4488 Strategies:[0xc003d0ad80 0xc003d8abd0]}
2025/08/01 15:42:05 🔔 Container is ready: 5201573638cd
DSN: postgres://testuser:testpass@localhost:36320/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36320/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:05 OK   000001_create_users_table.sql (10.72ms)
2025/08/01 15:42:05 OK   000002_create_projects_table.sql (12.52ms)
2025/08/01 15:42:05 OK   000003_create_log_entries_table.sql (19.24ms)
2025/08/01 15:42:05 OK   000004_create_tags_system.sql (16.19ms)
2025/08/01 15:42:05 OK   000005_create_auth_tables.sql (20.44ms)
2025/08/01 15:42:05 OK   000006_create_insights_table.sql (24.36ms)
2025/08/01 15:42:05 OK   000007_create_performance_indexes.sql (15.22ms)
2025/08/01 15:42:05 OK   000008_create_analytics_views.sql (29.45ms)
2025/08/01 15:42:05 OK   000009_development_data.sql (29.77ms)
2025/08/01 15:42:05 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36320/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36320/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:05 🐳 Stopping container: 5201573638cd
2025/08/01 15:42:06 ✅ Container stopped: 5201573638cd
2025/08/01 15:42:06 🐳 Terminating container: 5201573638cd
2025/08/01 15:42:06 🚫 Container terminated: 5201573638cd
--- PASS: TestLogEntryHandler_UpdateLogEntry (13.76s)
    --- PASS: TestLogEntryHandler_UpdateLogEntry/successful_log_entry_update (2.77s)
    --- PASS: TestLogEntryHandler_UpdateLogEntry/unauthorized_access (2.54s)
    --- PASS: TestLogEntryHandler_UpdateLogEntry/log_entry_not_found (2.79s)
    --- PASS: TestLogEntryHandler_UpdateLogEntry/invalid_request_body (2.88s)
    --- PASS: TestLogEntryHandler_UpdateLogEntry/missing_log_entry_ID (2.78s)
=== RUN   TestLogEntryHandler_DeleteLogEntry
=== RUN   TestLogEntryHandler_DeleteLogEntry/successful_log_entry_deletion
2025/08/01 15:42:06 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:06 ✅ Container created: 16f2a5fb6314
2025/08/01 15:42:06 🐳 Starting container: 16f2a5fb6314
2025/08/01 15:42:06 ✅ Container started: 16f2a5fb6314
2025/08/01 15:42:06 ⏳ Waiting for container id 16f2a5fb6314 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003a0d408 Strategies:[0xc0047da3c0 0xc003f11800]}
2025/08/01 15:42:08 🔔 Container is ready: 16f2a5fb6314
DSN: postgres://testuser:testpass@localhost:36321/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36321/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:08 OK   000001_create_users_table.sql (10.46ms)
2025/08/01 15:42:08 OK   000002_create_projects_table.sql (12.43ms)
2025/08/01 15:42:08 OK   000003_create_log_entries_table.sql (19.02ms)
2025/08/01 15:42:08 OK   000004_create_tags_system.sql (16.19ms)
2025/08/01 15:42:08 OK   000005_create_auth_tables.sql (20.08ms)
2025/08/01 15:42:08 OK   000006_create_insights_table.sql (24.3ms)
2025/08/01 15:42:08 OK   000007_create_performance_indexes.sql (15.48ms)
2025/08/01 15:42:08 OK   000008_create_analytics_views.sql (27.6ms)
2025/08/01 15:42:08 OK   000009_development_data.sql (29.46ms)
2025/08/01 15:42:08 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36321/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36321/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:08 🐳 Stopping container: 16f2a5fb6314
2025/08/01 15:42:09 ✅ Container stopped: 16f2a5fb6314
2025/08/01 15:42:09 🐳 Terminating container: 16f2a5fb6314
2025/08/01 15:42:09 🚫 Container terminated: 16f2a5fb6314
=== RUN   TestLogEntryHandler_DeleteLogEntry/unauthorized_access
2025/08/01 15:42:09 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:09 ✅ Container created: 460afb18887a
2025/08/01 15:42:09 🐳 Starting container: 460afb18887a
2025/08/01 15:42:09 ✅ Container started: 460afb18887a
2025/08/01 15:42:09 ⏳ Waiting for container id 460afb18887a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001f099a0 Strategies:[0xc002506b40 0xc003a03380]}
2025/08/01 15:42:11 🔔 Container is ready: 460afb18887a
DSN: postgres://testuser:testpass@localhost:36322/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36322/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:11 OK   000001_create_users_table.sql (10.28ms)
2025/08/01 15:42:11 OK   000002_create_projects_table.sql (12.34ms)
2025/08/01 15:42:11 OK   000003_create_log_entries_table.sql (20.33ms)
2025/08/01 15:42:11 OK   000004_create_tags_system.sql (16.76ms)
2025/08/01 15:42:11 OK   000005_create_auth_tables.sql (21.57ms)
2025/08/01 15:42:11 OK   000006_create_insights_table.sql (26.2ms)
2025/08/01 15:42:11 OK   000007_create_performance_indexes.sql (16.91ms)
2025/08/01 15:42:11 OK   000008_create_analytics_views.sql (30.19ms)
2025/08/01 15:42:11 OK   000009_development_data.sql (29.8ms)
2025/08/01 15:42:11 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36322/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36322/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:11 🐳 Stopping container: 460afb18887a
2025/08/01 15:42:11 ✅ Container stopped: 460afb18887a
2025/08/01 15:42:11 🐳 Terminating container: 460afb18887a
2025/08/01 15:42:11 🚫 Container terminated: 460afb18887a
=== RUN   TestLogEntryHandler_DeleteLogEntry/log_entry_not_found
2025/08/01 15:42:11 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:11 ✅ Container created: f13516aa0785
2025/08/01 15:42:11 🐳 Starting container: f13516aa0785
2025/08/01 15:42:11 ✅ Container started: f13516aa0785
2025/08/01 15:42:11 ⏳ Waiting for container id f13516aa0785 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004821858 Strategies:[0xc003cffce0 0xc0038fd5c0]}
2025/08/01 15:42:13 🔔 Container is ready: f13516aa0785
DSN: postgres://testuser:testpass@localhost:36323/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36323/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:13 OK   000001_create_users_table.sql (11.29ms)
2025/08/01 15:42:13 OK   000002_create_projects_table.sql (12.88ms)
2025/08/01 15:42:13 OK   000003_create_log_entries_table.sql (20.86ms)
2025/08/01 15:42:13 OK   000004_create_tags_system.sql (17.36ms)
2025/08/01 15:42:13 OK   000005_create_auth_tables.sql (21.36ms)
2025/08/01 15:42:13 OK   000006_create_insights_table.sql (26.01ms)
2025/08/01 15:42:13 OK   000007_create_performance_indexes.sql (15.62ms)
2025/08/01 15:42:13 OK   000008_create_analytics_views.sql (28.86ms)
2025/08/01 15:42:14 OK   000009_development_data.sql (29.21ms)
2025/08/01 15:42:14 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36323/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36323/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:14 🐳 Stopping container: f13516aa0785
2025/08/01 15:42:14 ✅ Container stopped: f13516aa0785
2025/08/01 15:42:14 🐳 Terminating container: f13516aa0785
2025/08/01 15:42:14 🚫 Container terminated: f13516aa0785
=== RUN   TestLogEntryHandler_DeleteLogEntry/missing_log_entry_ID
2025/08/01 15:42:14 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:14 ✅ Container created: 6b93afdfd1a2
2025/08/01 15:42:14 🐳 Starting container: 6b93afdfd1a2
2025/08/01 15:42:14 ✅ Container started: 6b93afdfd1a2
2025/08/01 15:42:14 ⏳ Waiting for container id 6b93afdfd1a2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00468c248 Strategies:[0xc0039488a0 0xc0048aab40]}
2025/08/01 15:42:16 🔔 Container is ready: 6b93afdfd1a2
DSN: postgres://testuser:testpass@localhost:36324/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36324/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:16 OK   000001_create_users_table.sql (10.32ms)
2025/08/01 15:42:16 OK   000002_create_projects_table.sql (12.15ms)
2025/08/01 15:42:16 OK   000003_create_log_entries_table.sql (18.92ms)
2025/08/01 15:42:16 OK   000004_create_tags_system.sql (16.11ms)
2025/08/01 15:42:16 OK   000005_create_auth_tables.sql (20.16ms)
2025/08/01 15:42:16 OK   000006_create_insights_table.sql (24.19ms)
2025/08/01 15:42:16 OK   000007_create_performance_indexes.sql (15.7ms)
2025/08/01 15:42:16 OK   000008_create_analytics_views.sql (27.83ms)
2025/08/01 15:42:16 OK   000009_development_data.sql (28.46ms)
2025/08/01 15:42:16 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36324/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36324/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:17 🐳 Stopping container: 6b93afdfd1a2
2025/08/01 15:42:17 ✅ Container stopped: 6b93afdfd1a2
2025/08/01 15:42:17 🐳 Terminating container: 6b93afdfd1a2
2025/08/01 15:42:17 🚫 Container terminated: 6b93afdfd1a2
--- PASS: TestLogEntryHandler_DeleteLogEntry (11.21s)
    --- PASS: TestLogEntryHandler_DeleteLogEntry/successful_log_entry_deletion (2.89s)
    --- PASS: TestLogEntryHandler_DeleteLogEntry/unauthorized_access (2.67s)
    --- PASS: TestLogEntryHandler_DeleteLogEntry/log_entry_not_found (2.83s)
    --- PASS: TestLogEntryHandler_DeleteLogEntry/missing_log_entry_ID (2.82s)
=== RUN   TestLogEntryHandler_BulkCreateLogEntries
=== RUN   TestLogEntryHandler_BulkCreateLogEntries/successful_bulk_creation
2025/08/01 15:42:17 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:17 ✅ Container created: 33ca164fd1c6
2025/08/01 15:42:17 🐳 Starting container: 33ca164fd1c6
2025/08/01 15:42:17 ✅ Container started: 33ca164fd1c6
2025/08/01 15:42:17 ⏳ Waiting for container id 33ca164fd1c6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0047178a0 Strategies:[0xc004730c60 0xc004713e90]}
2025/08/01 15:42:19 🔔 Container is ready: 33ca164fd1c6
DSN: postgres://testuser:testpass@localhost:36325/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36325/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:19 OK   000001_create_users_table.sql (11.11ms)
2025/08/01 15:42:19 OK   000002_create_projects_table.sql (14ms)
2025/08/01 15:42:19 OK   000003_create_log_entries_table.sql (20.08ms)
2025/08/01 15:42:19 OK   000004_create_tags_system.sql (17.15ms)
2025/08/01 15:42:19 OK   000005_create_auth_tables.sql (20.63ms)
2025/08/01 15:42:19 OK   000006_create_insights_table.sql (27.09ms)
2025/08/01 15:42:19 OK   000007_create_performance_indexes.sql (16.51ms)
2025/08/01 15:42:19 OK   000008_create_analytics_views.sql (29.74ms)
2025/08/01 15:42:19 OK   000009_development_data.sql (29.19ms)
2025/08/01 15:42:19 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36325/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36325/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:19 🐳 Stopping container: 33ca164fd1c6
2025/08/01 15:42:20 ✅ Container stopped: 33ca164fd1c6
2025/08/01 15:42:20 🐳 Terminating container: 33ca164fd1c6
2025/08/01 15:42:20 🚫 Container terminated: 33ca164fd1c6
=== RUN   TestLogEntryHandler_BulkCreateLogEntries/unauthorized_access
2025/08/01 15:42:20 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:20 ✅ Container created: 338c558a6cf9
2025/08/01 15:42:20 🐳 Starting container: 338c558a6cf9
2025/08/01 15:42:20 ✅ Container started: 338c558a6cf9
2025/08/01 15:42:20 ⏳ Waiting for container id 338c558a6cf9 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00469aca8 Strategies:[0xc0044f8d80 0xc0044d3170]}
2025/08/01 15:42:22 🔔 Container is ready: 338c558a6cf9
DSN: postgres://testuser:testpass@localhost:36326/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36326/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:22 OK   000001_create_users_table.sql (10.65ms)
2025/08/01 15:42:22 OK   000002_create_projects_table.sql (12.58ms)
2025/08/01 15:42:22 OK   000003_create_log_entries_table.sql (19.39ms)
2025/08/01 15:42:22 OK   000004_create_tags_system.sql (16.48ms)
2025/08/01 15:42:22 OK   000005_create_auth_tables.sql (20.43ms)
2025/08/01 15:42:22 OK   000006_create_insights_table.sql (24.51ms)
2025/08/01 15:42:22 OK   000007_create_performance_indexes.sql (15.79ms)
2025/08/01 15:42:22 OK   000008_create_analytics_views.sql (27.73ms)
2025/08/01 15:42:22 OK   000009_development_data.sql (28.9ms)
2025/08/01 15:42:22 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36326/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36326/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:22 🐳 Stopping container: 338c558a6cf9
2025/08/01 15:42:22 ✅ Container stopped: 338c558a6cf9
2025/08/01 15:42:22 🐳 Terminating container: 338c558a6cf9
2025/08/01 15:42:22 🚫 Container terminated: 338c558a6cf9
=== RUN   TestLogEntryHandler_BulkCreateLogEntries/invalid_request_body
2025/08/01 15:42:22 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:22 ✅ Container created: 2e44854bbcaf
2025/08/01 15:42:22 🐳 Starting container: 2e44854bbcaf
2025/08/01 15:42:23 ✅ Container started: 2e44854bbcaf
2025/08/01 15:42:23 ⏳ Waiting for container id 2e44854bbcaf image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003a0d798 Strategies:[0xc0048fdd40 0xc001fa2000]}
2025/08/01 15:42:24 🔔 Container is ready: 2e44854bbcaf
DSN: postgres://testuser:testpass@localhost:36327/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36327/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:24 OK   000001_create_users_table.sql (10.73ms)
2025/08/01 15:42:24 OK   000002_create_projects_table.sql (12.47ms)
2025/08/01 15:42:24 OK   000003_create_log_entries_table.sql (19.82ms)
2025/08/01 15:42:24 OK   000004_create_tags_system.sql (16.67ms)
2025/08/01 15:42:24 OK   000005_create_auth_tables.sql (20.6ms)
2025/08/01 15:42:25 OK   000006_create_insights_table.sql (24.48ms)
2025/08/01 15:42:25 OK   000007_create_performance_indexes.sql (15.63ms)
2025/08/01 15:42:25 OK   000008_create_analytics_views.sql (28.15ms)
2025/08/01 15:42:25 OK   000009_development_data.sql (29.02ms)
2025/08/01 15:42:25 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36327/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36327/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:25 🐳 Stopping container: 2e44854bbcaf
2025/08/01 15:42:25 ✅ Container stopped: 2e44854bbcaf
2025/08/01 15:42:25 🐳 Terminating container: 2e44854bbcaf
2025/08/01 15:42:25 🚫 Container terminated: 2e44854bbcaf
=== RUN   TestLogEntryHandler_BulkCreateLogEntries/empty_entries_array
2025/08/01 15:42:25 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:25 ✅ Container created: b652b7275457
2025/08/01 15:42:25 🐳 Starting container: b652b7275457
2025/08/01 15:42:25 ✅ Container started: b652b7275457
2025/08/01 15:42:25 ⏳ Waiting for container id b652b7275457 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0045aba78 Strategies:[0xc0047db680 0xc003b08540]}
2025/08/01 15:42:27 🔔 Container is ready: b652b7275457
DSN: postgres://testuser:testpass@localhost:36328/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36328/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:27 OK   000001_create_users_table.sql (10.48ms)
2025/08/01 15:42:27 OK   000002_create_projects_table.sql (12.42ms)
2025/08/01 15:42:27 OK   000003_create_log_entries_table.sql (20.47ms)
2025/08/01 15:42:27 OK   000004_create_tags_system.sql (17.03ms)
2025/08/01 15:42:27 OK   000005_create_auth_tables.sql (20.61ms)
2025/08/01 15:42:27 OK   000006_create_insights_table.sql (25.09ms)
2025/08/01 15:42:27 OK   000007_create_performance_indexes.sql (15.91ms)
2025/08/01 15:42:27 OK   000008_create_analytics_views.sql (28.29ms)
2025/08/01 15:42:27 OK   000009_development_data.sql (29.02ms)
2025/08/01 15:42:27 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36328/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36328/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:28 🐳 Stopping container: b652b7275457
2025/08/01 15:42:28 ✅ Container stopped: b652b7275457
2025/08/01 15:42:28 🐳 Terminating container: b652b7275457
2025/08/01 15:42:28 🚫 Container terminated: b652b7275457
=== RUN   TestLogEntryHandler_BulkCreateLogEntries/partial_success
2025/08/01 15:42:28 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:28 ✅ Container created: fa10cb919673
2025/08/01 15:42:28 🐳 Starting container: fa10cb919673
2025/08/01 15:42:28 ✅ Container started: fa10cb919673
2025/08/01 15:42:28 ⏳ Waiting for container id fa10cb919673 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003af61b0 Strategies:[0xc00003f4a0 0xc004604990]}
2025/08/01 15:42:30 🔔 Container is ready: fa10cb919673
DSN: postgres://testuser:testpass@localhost:36329/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36329/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:30 OK   000001_create_users_table.sql (11.22ms)
2025/08/01 15:42:30 OK   000002_create_projects_table.sql (12.45ms)
2025/08/01 15:42:30 OK   000003_create_log_entries_table.sql (20.49ms)
2025/08/01 15:42:30 OK   000004_create_tags_system.sql (17.77ms)
2025/08/01 15:42:30 OK   000005_create_auth_tables.sql (21.65ms)
2025/08/01 15:42:30 OK   000006_create_insights_table.sql (25.79ms)
2025/08/01 15:42:30 OK   000007_create_performance_indexes.sql (16.02ms)
2025/08/01 15:42:30 OK   000008_create_analytics_views.sql (29.84ms)
2025/08/01 15:42:30 OK   000009_development_data.sql (29.57ms)
2025/08/01 15:42:30 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36329/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36329/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:31 🐳 Stopping container: fa10cb919673
2025/08/01 15:42:31 ✅ Container stopped: fa10cb919673
2025/08/01 15:42:31 🐳 Terminating container: fa10cb919673
2025/08/01 15:42:31 🚫 Container terminated: fa10cb919673
--- PASS: TestLogEntryHandler_BulkCreateLogEntries (13.92s)
    --- PASS: TestLogEntryHandler_BulkCreateLogEntries/successful_bulk_creation (2.87s)
    --- PASS: TestLogEntryHandler_BulkCreateLogEntries/unauthorized_access (2.57s)
    --- PASS: TestLogEntryHandler_BulkCreateLogEntries/invalid_request_body (2.82s)
    --- PASS: TestLogEntryHandler_BulkCreateLogEntries/empty_entries_array (2.88s)
    --- PASS: TestLogEntryHandler_BulkCreateLogEntries/partial_success (2.78s)
=== RUN   TestParseLogEntryFilters
=== RUN   TestParseLogEntryFilters/valid_filters
=== RUN   TestParseLogEntryFilters/no_filters
=== RUN   TestParseLogEntryFilters/invalid_date_format
=== RUN   TestParseLogEntryFilters/invalid_activity_type
=== RUN   TestParseLogEntryFilters/invalid_value_rating
=== RUN   TestParseLogEntryFilters/invalid_impact_level
--- PASS: TestParseLogEntryFilters (0.00s)
    --- PASS: TestParseLogEntryFilters/valid_filters (0.00s)
    --- PASS: TestParseLogEntryFilters/no_filters (0.00s)
    --- PASS: TestParseLogEntryFilters/invalid_date_format (0.00s)
    --- PASS: TestParseLogEntryFilters/invalid_activity_type (0.00s)
    --- PASS: TestParseLogEntryFilters/invalid_value_rating (0.00s)
    --- PASS: TestParseLogEntryFilters/invalid_impact_level (0.00s)
=== RUN   TestParsePagination
=== RUN   TestParsePagination/default_values
=== RUN   TestParsePagination/custom_values
=== RUN   TestParsePagination/invalid_page
=== RUN   TestParsePagination/negative_page
=== RUN   TestParsePagination/limit_too_high
=== RUN   TestParsePagination/negative_limit
=== RUN   TestParsePagination/invalid_page_string
=== RUN   TestParsePagination/invalid_limit_string
=== RUN   TestParsePagination/max_valid_limit
=== RUN   TestParsePagination/boundary_page
--- PASS: TestParsePagination (0.00s)
    --- PASS: TestParsePagination/default_values (0.00s)
    --- PASS: TestParsePagination/custom_values (0.00s)
    --- PASS: TestParsePagination/invalid_page (0.00s)
    --- PASS: TestParsePagination/negative_page (0.00s)
    --- PASS: TestParsePagination/limit_too_high (0.00s)
    --- PASS: TestParsePagination/negative_limit (0.00s)
    --- PASS: TestParsePagination/invalid_page_string (0.00s)
    --- PASS: TestParsePagination/invalid_limit_string (0.00s)
    --- PASS: TestParsePagination/max_valid_limit (0.00s)
    --- PASS: TestParsePagination/boundary_page (0.00s)
=== RUN   TestPaginate
=== RUN   TestPaginate/normal_pagination
=== RUN   TestPaginate/first_page
=== RUN   TestPaginate/last_page
=== RUN   TestPaginate/page_beyond_available_data
=== RUN   TestPaginate/partial_last_page
=== RUN   TestPaginate/empty_entries
--- PASS: TestPaginate (0.00s)
    --- PASS: TestPaginate/normal_pagination (0.00s)
    --- PASS: TestPaginate/first_page (0.00s)
    --- PASS: TestPaginate/last_page (0.00s)
    --- PASS: TestPaginate/page_beyond_available_data (0.00s)
    --- PASS: TestPaginate/partial_last_page (0.00s)
    --- PASS: TestPaginate/empty_entries (0.00s)
=== RUN   TestProjectHandler_CreateProject_Comprehensive
=== RUN   TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios
2025/08/01 15:42:31 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:31 ✅ Container created: b919a3bb6609
2025/08/01 15:42:31 🐳 Starting container: b919a3bb6609
2025/08/01 15:42:31 ✅ Container started: b919a3bb6609
2025/08/01 15:42:31 ⏳ Waiting for container id b919a3bb6609 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0038bbdf0 Strategies:[0xc003f27bc0 0xc0019edc20]}
2025/08/01 15:42:33 🔔 Container is ready: b919a3bb6609
DSN: postgres://testuser:testpass@localhost:36330/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36330/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:33 OK   000001_create_users_table.sql (10.02ms)
2025/08/01 15:42:33 OK   000002_create_projects_table.sql (11.88ms)
2025/08/01 15:42:33 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:42:33 OK   000004_create_tags_system.sql (17.04ms)
2025/08/01 15:42:33 OK   000005_create_auth_tables.sql (21.18ms)
2025/08/01 15:42:33 OK   000006_create_insights_table.sql (25.17ms)
2025/08/01 15:42:33 OK   000007_create_performance_indexes.sql (15.74ms)
2025/08/01 15:42:33 OK   000008_create_analytics_views.sql (28.72ms)
2025/08/01 15:42:33 OK   000009_development_data.sql (29.35ms)
2025/08/01 15:42:33 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36330/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36330/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
=== RUN   TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_active
=== RUN   TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_completed
=== RUN   TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_on_hold
=== RUN   TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_cancelled
2025/08/01 15:42:33 🐳 Stopping container: b919a3bb6609
2025/08/01 15:42:34 ✅ Container stopped: b919a3bb6609
2025/08/01 15:42:34 🐳 Terminating container: b919a3bb6609
2025/08/01 15:42:34 🚫 Container terminated: b919a3bb6609
=== RUN   TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios
2025/08/01 15:42:34 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:34 ✅ Container created: 9a7765ff3e7c
2025/08/01 15:42:34 🐳 Starting container: 9a7765ff3e7c
2025/08/01 15:42:34 ✅ Container started: 9a7765ff3e7c
2025/08/01 15:42:34 ⏳ Waiting for container id 9a7765ff3e7c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006cd5210 Strategies:[0xc006ce2780 0xc006cee9c0]}
2025/08/01 15:42:36 🔔 Container is ready: 9a7765ff3e7c
DSN: postgres://testuser:testpass@localhost:36331/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36331/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:36 OK   000001_create_users_table.sql (10.36ms)
2025/08/01 15:42:36 OK   000002_create_projects_table.sql (12.2ms)
2025/08/01 15:42:36 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:42:36 OK   000004_create_tags_system.sql (16.15ms)
2025/08/01 15:42:36 OK   000005_create_auth_tables.sql (20.4ms)
2025/08/01 15:42:36 OK   000006_create_insights_table.sql (24.8ms)
2025/08/01 15:42:36 OK   000007_create_performance_indexes.sql (15.76ms)
2025/08/01 15:42:36 OK   000008_create_analytics_views.sql (27.99ms)
2025/08/01 15:42:36 OK   000009_development_data.sql (28.53ms)
2025/08/01 15:42:36 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36331/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36331/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
=== RUN   TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/empty_name
=== RUN   TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/invalid_color_format
=== RUN   TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/very_long_name
2025/08/01 15:42:36 🐳 Stopping container: 9a7765ff3e7c
2025/08/01 15:42:36 ✅ Container stopped: 9a7765ff3e7c
2025/08/01 15:42:36 🐳 Terminating container: 9a7765ff3e7c
2025/08/01 15:42:36 🚫 Container terminated: 9a7765ff3e7c
--- PASS: TestProjectHandler_CreateProject_Comprehensive (5.64s)
    --- PASS: TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios (2.85s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_active (0.01s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_completed (0.00s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_on_hold (0.00s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/comprehensive_project_creation_scenarios/create_project_with_status_cancelled (0.00s)
    --- PASS: TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios (2.79s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/empty_name (0.00s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/invalid_color_format (0.00s)
        --- PASS: TestProjectHandler_CreateProject_Comprehensive/project_validation_scenarios/very_long_name (0.00s)
=== RUN   TestProjectHandler_UpdateProject_Comprehensive
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios
2025/08/01 15:42:36 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:37 ✅ Container created: d08043b75c03
2025/08/01 15:42:37 🐳 Starting container: d08043b75c03
2025/08/01 15:42:37 ✅ Container started: d08043b75c03
2025/08/01 15:42:37 ⏳ Waiting for container id d08043b75c03 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007dba710 Strategies:[0xc006a5f560 0xc007f0f470]}
2025/08/01 15:42:38 🔔 Container is ready: d08043b75c03
DSN: postgres://testuser:testpass@localhost:36332/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36332/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:39 OK   000001_create_users_table.sql (9.94ms)
2025/08/01 15:42:39 OK   000002_create_projects_table.sql (12.11ms)
2025/08/01 15:42:39 OK   000003_create_log_entries_table.sql (19.71ms)
2025/08/01 15:42:39 OK   000004_create_tags_system.sql (16.84ms)
2025/08/01 15:42:39 OK   000005_create_auth_tables.sql (20.96ms)
2025/08/01 15:42:39 OK   000006_create_insights_table.sql (25.46ms)
2025/08/01 15:42:39 OK   000007_create_performance_indexes.sql (15.78ms)
2025/08/01 15:42:39 OK   000008_create_analytics_views.sql (29.8ms)
2025/08/01 15:42:39 OK   000009_development_data.sql (30.25ms)
2025/08/01 15:42:39 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36332/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36332/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_active_to_on_hold
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_on_hold_to_active
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_active_to_completed
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_completed_to_cancelled
2025/08/01 15:42:39 🐳 Stopping container: d08043b75c03
2025/08/01 15:42:39 ✅ Container stopped: d08043b75c03
2025/08/01 15:42:39 🐳 Terminating container: d08043b75c03
2025/08/01 15:42:39 🚫 Container terminated: d08043b75c03
=== RUN   TestProjectHandler_UpdateProject_Comprehensive/partial_update_scenarios
2025/08/01 15:42:39 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:39 ✅ Container created: 37f596cc5849
2025/08/01 15:42:39 🐳 Starting container: 37f596cc5849
2025/08/01 15:42:40 ✅ Container started: 37f596cc5849
2025/08/01 15:42:40 ⏳ Waiting for container id 37f596cc5849 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003a0dbb0 Strategies:[0xc006c12420 0xc003df94d0]}
2025/08/01 15:42:41 🔔 Container is ready: 37f596cc5849
DSN: postgres://testuser:testpass@localhost:36333/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36333/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:41 OK   000001_create_users_table.sql (10.28ms)
2025/08/01 15:42:41 OK   000002_create_projects_table.sql (12.05ms)
2025/08/01 15:42:41 OK   000003_create_log_entries_table.sql (19.1ms)
2025/08/01 15:42:41 OK   000004_create_tags_system.sql (16.15ms)
2025/08/01 15:42:41 OK   000005_create_auth_tables.sql (20.39ms)
2025/08/01 15:42:41 OK   000006_create_insights_table.sql (24.8ms)
2025/08/01 15:42:42 OK   000007_create_performance_indexes.sql (15.37ms)
2025/08/01 15:42:42 OK   000008_create_analytics_views.sql (28.3ms)
2025/08/01 15:42:42 OK   000009_development_data.sql (29.17ms)
2025/08/01 15:42:42 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36333/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36333/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:42 🐳 Stopping container: 37f596cc5849
2025/08/01 15:42:42 ✅ Container stopped: 37f596cc5849
2025/08/01 15:42:42 🐳 Terminating container: 37f596cc5849
2025/08/01 15:42:42 🚫 Container terminated: 37f596cc5849
--- PASS: TestProjectHandler_UpdateProject_Comprehensive (5.72s)
    --- PASS: TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios (2.87s)
        --- PASS: TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_active_to_on_hold (0.02s)
        --- PASS: TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_on_hold_to_active (0.01s)
        --- PASS: TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_active_to_completed (0.01s)
        --- PASS: TestProjectHandler_UpdateProject_Comprehensive/comprehensive_update_scenarios/transition_from_completed_to_cancelled (0.01s)
    --- PASS: TestProjectHandler_UpdateProject_Comprehensive/partial_update_scenarios (2.85s)
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive/authorization_edge_cases
2025/08/01 15:42:42 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:42 ✅ Container created: 4ec22992ec0f
2025/08/01 15:42:42 🐳 Starting container: 4ec22992ec0f
2025/08/01 15:42:42 ✅ Container started: 4ec22992ec0f
2025/08/01 15:42:42 ⏳ Waiting for container id 4ec22992ec0f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0069d4378 Strategies:[0xc0044006c0 0xc000465e60]}
2025/08/01 15:42:44 🔔 Container is ready: 4ec22992ec0f
DSN: postgres://testuser:testpass@localhost:36334/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36334/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:44 OK   000001_create_users_table.sql (11.68ms)
2025/08/01 15:42:44 OK   000002_create_projects_table.sql (13.26ms)
2025/08/01 15:42:44 OK   000003_create_log_entries_table.sql (20.23ms)
2025/08/01 15:42:44 OK   000004_create_tags_system.sql (16.42ms)
2025/08/01 15:42:44 OK   000005_create_auth_tables.sql (20.47ms)
2025/08/01 15:42:44 OK   000006_create_insights_table.sql (24.55ms)
2025/08/01 15:42:44 OK   000007_create_performance_indexes.sql (15.52ms)
2025/08/01 15:42:44 OK   000008_create_analytics_views.sql (27.98ms)
2025/08/01 15:42:44 OK   000009_development_data.sql (28.9ms)
2025/08/01 15:42:44 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36334/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36334/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:45 🐳 Stopping container: 4ec22992ec0f
2025/08/01 15:42:45 ✅ Container stopped: 4ec22992ec0f
2025/08/01 15:42:45 🐳 Terminating container: 4ec22992ec0f
2025/08/01 15:42:45 🚫 Container terminated: 4ec22992ec0f
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios
2025/08/01 15:42:45 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:45 ✅ Container created: 8679134dfffb
2025/08/01 15:42:45 🐳 Starting container: 8679134dfffb
2025/08/01 15:42:45 ✅ Container started: 8679134dfffb
2025/08/01 15:42:45 ⏳ Waiting for container id 8679134dfffb image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00445c5f8 Strategies:[0xc00480f2c0 0xc0032f4570]}
2025/08/01 15:42:47 🔔 Container is ready: 8679134dfffb
DSN: postgres://testuser:testpass@localhost:36335/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36335/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:47 OK   000001_create_users_table.sql (11.22ms)
2025/08/01 15:42:47 OK   000002_create_projects_table.sql (14.53ms)
2025/08/01 15:42:47 OK   000003_create_log_entries_table.sql (20.51ms)
2025/08/01 15:42:47 OK   000004_create_tags_system.sql (15.96ms)
2025/08/01 15:42:47 OK   000005_create_auth_tables.sql (20.11ms)
2025/08/01 15:42:48 OK   000006_create_insights_table.sql (24.49ms)
2025/08/01 15:42:48 OK   000007_create_performance_indexes.sql (15.62ms)
2025/08/01 15:42:48 OK   000008_create_analytics_views.sql (27.85ms)
2025/08/01 15:42:48 OK   000009_development_data.sql (29.09ms)
2025/08/01 15:42:48 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36335/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36335/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/malformed_token
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/missing_Bearer_prefix
=== RUN   TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/empty_header
2025/08/01 15:42:48 🐳 Stopping container: 8679134dfffb
2025/08/01 15:42:48 ✅ Container stopped: 8679134dfffb
2025/08/01 15:42:48 🐳 Terminating container: 8679134dfffb
2025/08/01 15:42:48 🚫 Container terminated: 8679134dfffb
--- PASS: TestProjectHandler_ErrorHandling_Comprehensive (5.71s)
    --- PASS: TestProjectHandler_ErrorHandling_Comprehensive/authorization_edge_cases (3.04s)
    --- PASS: TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios (2.67s)
        --- PASS: TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/malformed_token (0.00s)
        --- PASS: TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/missing_Bearer_prefix (0.00s)
        --- PASS: TestProjectHandler_ErrorHandling_Comprehensive/invalid_token_scenarios/empty_header (0.00s)
=== RUN   TestProjectHandler_DefaultProject_Comprehensive
=== RUN   TestProjectHandler_DefaultProject_Comprehensive/default_project_handling
2025/08/01 15:42:48 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:48 ✅ Container created: d2f8ca41386d
2025/08/01 15:42:48 🐳 Starting container: d2f8ca41386d
2025/08/01 15:42:48 ✅ Container started: d2f8ca41386d
2025/08/01 15:42:48 ⏳ Waiting for container id d2f8ca41386d image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003cc8b80 Strategies:[0xc000301920 0xc003cda300]}
2025/08/01 15:42:50 🔔 Container is ready: d2f8ca41386d
DSN: postgres://testuser:testpass@localhost:36336/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36336/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:50 OK   000001_create_users_table.sql (11.24ms)
2025/08/01 15:42:50 OK   000002_create_projects_table.sql (13ms)
2025/08/01 15:42:50 OK   000003_create_log_entries_table.sql (20.53ms)
2025/08/01 15:42:50 OK   000004_create_tags_system.sql (17.55ms)
2025/08/01 15:42:50 OK   000005_create_auth_tables.sql (20.72ms)
2025/08/01 15:42:50 OK   000006_create_insights_table.sql (24.78ms)
2025/08/01 15:42:50 OK   000007_create_performance_indexes.sql (15.79ms)
2025/08/01 15:42:50 OK   000008_create_analytics_views.sql (27.94ms)
2025/08/01 15:42:50 OK   000009_development_data.sql (28.8ms)
2025/08/01 15:42:50 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36336/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36336/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
    projects_comprehensive_test.go:327: Service correctly prevents multiple default projects
2025/08/01 15:42:50 🐳 Stopping container: d2f8ca41386d
2025/08/01 15:42:51 ✅ Container stopped: d2f8ca41386d
2025/08/01 15:42:51 🐳 Terminating container: d2f8ca41386d
2025/08/01 15:42:51 🚫 Container terminated: d2f8ca41386d
--- PASS: TestProjectHandler_DefaultProject_Comprehensive (2.79s)
    --- PASS: TestProjectHandler_DefaultProject_Comprehensive/default_project_handling (2.79s)
=== RUN   TestProjectHandler_CreateProject
=== RUN   TestProjectHandler_CreateProject/successful_project_creation
2025/08/01 15:42:51 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:51 ✅ Container created: d021efcda4d0
2025/08/01 15:42:51 🐳 Starting container: d021efcda4d0
2025/08/01 15:42:51 ✅ Container started: d021efcda4d0
2025/08/01 15:42:51 ⏳ Waiting for container id d021efcda4d0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007d442d8 Strategies:[0xc003edf140 0xc006a53ad0]}
2025/08/01 15:42:53 🔔 Container is ready: d021efcda4d0
DSN: postgres://testuser:testpass@localhost:36337/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36337/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:53 OK   000001_create_users_table.sql (11.44ms)
2025/08/01 15:42:53 OK   000002_create_projects_table.sql (13.89ms)
2025/08/01 15:42:53 OK   000003_create_log_entries_table.sql (21.91ms)
2025/08/01 15:42:53 OK   000004_create_tags_system.sql (17.54ms)
2025/08/01 15:42:53 OK   000005_create_auth_tables.sql (21.55ms)
2025/08/01 15:42:53 OK   000006_create_insights_table.sql (25.11ms)
2025/08/01 15:42:53 OK   000007_create_performance_indexes.sql (15.41ms)
2025/08/01 15:42:53 OK   000008_create_analytics_views.sql (28.28ms)
2025/08/01 15:42:53 OK   000009_development_data.sql (28.91ms)
2025/08/01 15:42:53 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36337/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36337/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:53 🐳 Stopping container: d021efcda4d0
2025/08/01 15:42:53 ✅ Container stopped: d021efcda4d0
2025/08/01 15:42:53 🐳 Terminating container: d021efcda4d0
2025/08/01 15:42:54 🚫 Container terminated: d021efcda4d0
=== RUN   TestProjectHandler_CreateProject/unauthorized_access
2025/08/01 15:42:54 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:54 ✅ Container created: b423e0b18ec9
2025/08/01 15:42:54 🐳 Starting container: b423e0b18ec9
2025/08/01 15:42:54 ✅ Container started: b423e0b18ec9
2025/08/01 15:42:54 ⏳ Waiting for container id b423e0b18ec9 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002480880 Strategies:[0xc0048fc900 0xc001e7bbc0]}
2025/08/01 15:42:56 🔔 Container is ready: b423e0b18ec9
DSN: postgres://testuser:testpass@localhost:36338/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36338/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:56 OK   000001_create_users_table.sql (11.31ms)
2025/08/01 15:42:56 OK   000002_create_projects_table.sql (13.7ms)
2025/08/01 15:42:56 OK   000003_create_log_entries_table.sql (20.67ms)
2025/08/01 15:42:56 OK   000004_create_tags_system.sql (16.96ms)
2025/08/01 15:42:56 OK   000005_create_auth_tables.sql (21.07ms)
2025/08/01 15:42:56 OK   000006_create_insights_table.sql (24.98ms)
2025/08/01 15:42:56 OK   000007_create_performance_indexes.sql (16.89ms)
2025/08/01 15:42:56 OK   000008_create_analytics_views.sql (27.9ms)
2025/08/01 15:42:56 OK   000009_development_data.sql (28.34ms)
2025/08/01 15:42:56 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36338/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36338/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:56 🐳 Stopping container: b423e0b18ec9
2025/08/01 15:42:56 ✅ Container stopped: b423e0b18ec9
2025/08/01 15:42:56 🐳 Terminating container: b423e0b18ec9
2025/08/01 15:42:56 🚫 Container terminated: b423e0b18ec9
=== RUN   TestProjectHandler_CreateProject/invalid_request_body
2025/08/01 15:42:56 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:56 ✅ Container created: 745e2ba30309
2025/08/01 15:42:56 🐳 Starting container: 745e2ba30309
2025/08/01 15:42:56 ✅ Container started: 745e2ba30309
2025/08/01 15:42:56 ⏳ Waiting for container id 745e2ba30309 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00469b4d8 Strategies:[0xc006c12ea0 0xc0047aee10]}
2025/08/01 15:42:58 🔔 Container is ready: 745e2ba30309
DSN: postgres://testuser:testpass@localhost:36339/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36339/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:42:58 OK   000001_create_users_table.sql (11.04ms)
2025/08/01 15:42:58 OK   000002_create_projects_table.sql (12.36ms)
2025/08/01 15:42:58 OK   000003_create_log_entries_table.sql (19.94ms)
2025/08/01 15:42:58 OK   000004_create_tags_system.sql (16.75ms)
2025/08/01 15:42:58 OK   000005_create_auth_tables.sql (20.94ms)
2025/08/01 15:42:58 OK   000006_create_insights_table.sql (25.01ms)
2025/08/01 15:42:58 OK   000007_create_performance_indexes.sql (16.07ms)
2025/08/01 15:42:58 OK   000008_create_analytics_views.sql (25.84ms)
2025/08/01 15:42:58 OK   000009_development_data.sql (23.81ms)
2025/08/01 15:42:58 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36339/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36339/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:42:59 🐳 Stopping container: 745e2ba30309
2025/08/01 15:42:59 ✅ Container stopped: 745e2ba30309
2025/08/01 15:42:59 🐳 Terminating container: 745e2ba30309
2025/08/01 15:42:59 🚫 Container terminated: 745e2ba30309
--- PASS: TestProjectHandler_CreateProject (8.15s)
    --- PASS: TestProjectHandler_CreateProject/successful_project_creation (2.86s)
    --- PASS: TestProjectHandler_CreateProject/unauthorized_access (2.55s)
    --- PASS: TestProjectHandler_CreateProject/invalid_request_body (2.74s)
=== RUN   TestProjectHandler_GetProject
=== RUN   TestProjectHandler_GetProject/successful_project_retrieval
2025/08/01 15:42:59 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:42:59 ✅ Container created: ffd322260f38
2025/08/01 15:42:59 🐳 Starting container: ffd322260f38
2025/08/01 15:42:59 ✅ Container started: ffd322260f38
2025/08/01 15:42:59 ⏳ Waiting for container id ffd322260f38 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0022476f0 Strategies:[0xc0047317a0 0xc003f39d70]}
2025/08/01 15:43:01 🔔 Container is ready: ffd322260f38
DSN: postgres://testuser:testpass@localhost:36340/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36340/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:01 OK   000001_create_users_table.sql (11.78ms)
2025/08/01 15:43:01 OK   000002_create_projects_table.sql (13.77ms)
2025/08/01 15:43:01 OK   000003_create_log_entries_table.sql (22.85ms)
2025/08/01 15:43:01 OK   000004_create_tags_system.sql (18.62ms)
2025/08/01 15:43:01 OK   000005_create_auth_tables.sql (21.96ms)
2025/08/01 15:43:01 OK   000006_create_insights_table.sql (26.24ms)
2025/08/01 15:43:01 OK   000007_create_performance_indexes.sql (16.52ms)
2025/08/01 15:43:01 OK   000008_create_analytics_views.sql (29.74ms)
2025/08/01 15:43:01 OK   000009_development_data.sql (30.28ms)
2025/08/01 15:43:01 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36340/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36340/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:01 🐳 Stopping container: ffd322260f38
2025/08/01 15:43:02 ✅ Container stopped: ffd322260f38
2025/08/01 15:43:02 🐳 Terminating container: ffd322260f38
2025/08/01 15:43:02 🚫 Container terminated: ffd322260f38
=== RUN   TestProjectHandler_GetProject/project_not_found
2025/08/01 15:43:02 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:02 ✅ Container created: 109c7ca794a2
2025/08/01 15:43:02 🐳 Starting container: 109c7ca794a2
2025/08/01 15:43:02 ✅ Container started: 109c7ca794a2
2025/08/01 15:43:02 ⏳ Waiting for container id 109c7ca794a2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006cd5870 Strategies:[0xc007cd0900 0xc006449bf0]}
2025/08/01 15:43:04 🔔 Container is ready: 109c7ca794a2
DSN: postgres://testuser:testpass@localhost:36341/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36341/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:04 OK   000001_create_users_table.sql (10.82ms)
2025/08/01 15:43:04 OK   000002_create_projects_table.sql (12.5ms)
2025/08/01 15:43:04 OK   000003_create_log_entries_table.sql (19.4ms)
2025/08/01 15:43:04 OK   000004_create_tags_system.sql (16.62ms)
2025/08/01 15:43:04 OK   000005_create_auth_tables.sql (20.54ms)
2025/08/01 15:43:04 OK   000006_create_insights_table.sql (24.71ms)
2025/08/01 15:43:04 OK   000007_create_performance_indexes.sql (15.44ms)
2025/08/01 15:43:04 OK   000008_create_analytics_views.sql (28.19ms)
2025/08/01 15:43:04 OK   000009_development_data.sql (28.96ms)
2025/08/01 15:43:04 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36341/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36341/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:04 🐳 Stopping container: 109c7ca794a2
2025/08/01 15:43:04 ✅ Container stopped: 109c7ca794a2
2025/08/01 15:43:04 🐳 Terminating container: 109c7ca794a2
2025/08/01 15:43:04 🚫 Container terminated: 109c7ca794a2
--- PASS: TestProjectHandler_GetProject (5.59s)
    --- PASS: TestProjectHandler_GetProject/successful_project_retrieval (2.83s)
    --- PASS: TestProjectHandler_GetProject/project_not_found (2.76s)
=== RUN   TestProjectHandler_GetProjects
=== RUN   TestProjectHandler_GetProjects/successful_projects_list
2025/08/01 15:43:04 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:04 ✅ Container created: 223103a41186
2025/08/01 15:43:04 🐳 Starting container: 223103a41186
2025/08/01 15:43:05 ✅ Container started: 223103a41186
2025/08/01 15:43:05 ⏳ Waiting for container id 223103a41186 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00446c550 Strategies:[0xc0046dd8c0 0xc007f0e870]}
2025/08/01 15:43:06 🔔 Container is ready: 223103a41186
DSN: postgres://testuser:testpass@localhost:36342/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36342/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:07 OK   000001_create_users_table.sql (11.39ms)
2025/08/01 15:43:07 OK   000002_create_projects_table.sql (12.87ms)
2025/08/01 15:43:07 OK   000003_create_log_entries_table.sql (20.4ms)
2025/08/01 15:43:07 OK   000004_create_tags_system.sql (17.48ms)
2025/08/01 15:43:07 OK   000005_create_auth_tables.sql (21.76ms)
2025/08/01 15:43:07 OK   000006_create_insights_table.sql (25.43ms)
2025/08/01 15:43:07 OK   000007_create_performance_indexes.sql (16.08ms)
2025/08/01 15:43:07 OK   000008_create_analytics_views.sql (28.5ms)
2025/08/01 15:43:07 OK   000009_development_data.sql (29.69ms)
2025/08/01 15:43:07 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36342/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36342/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:07 🐳 Stopping container: 223103a41186
2025/08/01 15:43:07 ✅ Container stopped: 223103a41186
2025/08/01 15:43:07 🐳 Terminating container: 223103a41186
2025/08/01 15:43:07 🚫 Container terminated: 223103a41186
--- PASS: TestProjectHandler_GetProjects (2.81s)
    --- PASS: TestProjectHandler_GetProjects/successful_projects_list (2.81s)
=== RUN   TestProjectHandler_UpdateProject
=== RUN   TestProjectHandler_UpdateProject/successful_project_update
2025/08/01 15:43:07 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:07 ✅ Container created: e220401bf2c2
2025/08/01 15:43:07 🐳 Starting container: e220401bf2c2
2025/08/01 15:43:07 ✅ Container started: e220401bf2c2
2025/08/01 15:43:07 ⏳ Waiting for container id e220401bf2c2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006bacfd0 Strategies:[0xc003948720 0xc005b594a0]}
2025/08/01 15:43:09 🔔 Container is ready: e220401bf2c2
DSN: postgres://testuser:testpass@localhost:36343/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36343/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:09 OK   000001_create_users_table.sql (10.12ms)
2025/08/01 15:43:09 OK   000002_create_projects_table.sql (11.7ms)
2025/08/01 15:43:09 OK   000003_create_log_entries_table.sql (18.74ms)
2025/08/01 15:43:09 OK   000004_create_tags_system.sql (16.45ms)
2025/08/01 15:43:09 OK   000005_create_auth_tables.sql (19.55ms)
2025/08/01 15:43:09 OK   000006_create_insights_table.sql (23.96ms)
2025/08/01 15:43:09 OK   000007_create_performance_indexes.sql (15.18ms)
2025/08/01 15:43:09 OK   000008_create_analytics_views.sql (27.9ms)
2025/08/01 15:43:09 OK   000009_development_data.sql (27.95ms)
2025/08/01 15:43:09 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36343/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36343/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:10 🐳 Stopping container: e220401bf2c2
2025/08/01 15:43:10 ✅ Container stopped: e220401bf2c2
2025/08/01 15:43:10 🐳 Terminating container: e220401bf2c2
2025/08/01 15:43:10 🚫 Container terminated: e220401bf2c2
--- PASS: TestProjectHandler_UpdateProject (2.77s)
    --- PASS: TestProjectHandler_UpdateProject/successful_project_update (2.77s)
=== RUN   TestProjectHandler_DeleteProject
=== RUN   TestProjectHandler_DeleteProject/successful_project_deletion
2025/08/01 15:43:10 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:10 ✅ Container created: 28681c59caae
2025/08/01 15:43:10 🐳 Starting container: 28681c59caae
2025/08/01 15:43:10 ✅ Container started: 28681c59caae
2025/08/01 15:43:10 ⏳ Waiting for container id 28681c59caae image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00469b6b0 Strategies:[0xc007e4acc0 0xc003246630]}
2025/08/01 15:43:12 🔔 Container is ready: 28681c59caae
DSN: postgres://testuser:testpass@localhost:36344/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36344/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:12 OK   000001_create_users_table.sql (9.96ms)
2025/08/01 15:43:12 OK   000002_create_projects_table.sql (12.17ms)
2025/08/01 15:43:12 OK   000003_create_log_entries_table.sql (18.87ms)
2025/08/01 15:43:12 OK   000004_create_tags_system.sql (16.12ms)
2025/08/01 15:43:12 OK   000005_create_auth_tables.sql (20.69ms)
2025/08/01 15:43:12 OK   000006_create_insights_table.sql (24.41ms)
2025/08/01 15:43:12 OK   000007_create_performance_indexes.sql (15.32ms)
2025/08/01 15:43:12 OK   000008_create_analytics_views.sql (28.3ms)
2025/08/01 15:43:12 OK   000009_development_data.sql (28.56ms)
2025/08/01 15:43:12 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36344/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36344/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:13 🐳 Stopping container: 28681c59caae
2025/08/01 15:43:13 ✅ Container stopped: 28681c59caae
2025/08/01 15:43:13 🐳 Terminating container: 28681c59caae
2025/08/01 15:43:13 🚫 Container terminated: 28681c59caae
--- PASS: TestProjectHandler_DeleteProject (2.80s)
    --- PASS: TestProjectHandler_DeleteProject/successful_project_deletion (2.80s)
=== RUN   TestRespondWithError
=== RUN   TestRespondWithError/error_without_details
=== RUN   TestRespondWithError/error_with_details
=== RUN   TestRespondWithError/error_with_multiple_details_(only_first_is_used)
=== RUN   TestRespondWithError/not_found_error
--- PASS: TestRespondWithError (0.00s)
    --- PASS: TestRespondWithError/error_without_details (0.00s)
    --- PASS: TestRespondWithError/error_with_details (0.00s)
    --- PASS: TestRespondWithError/error_with_multiple_details_(only_first_is_used) (0.00s)
    --- PASS: TestRespondWithError/not_found_error (0.00s)
=== RUN   TestRespondWithSuccess
=== RUN   TestRespondWithSuccess/success_without_message
=== RUN   TestRespondWithSuccess/success_with_message
=== RUN   TestRespondWithSuccess/success_with_multiple_messages_(only_first_is_used)
=== RUN   TestRespondWithSuccess/success_with_nil_data
=== RUN   TestRespondWithSuccess/success_with_complex_data_structure
--- PASS: TestRespondWithSuccess (0.00s)
    --- PASS: TestRespondWithSuccess/success_without_message (0.00s)
    --- PASS: TestRespondWithSuccess/success_with_message (0.00s)
    --- PASS: TestRespondWithSuccess/success_with_multiple_messages_(only_first_is_used) (0.00s)
    --- PASS: TestRespondWithSuccess/success_with_nil_data (0.00s)
    --- PASS: TestRespondWithSuccess/success_with_complex_data_structure (0.00s)
=== RUN   TestRespondWithPagination
=== RUN   TestRespondWithPagination/paginated_list_with_data
=== RUN   TestRespondWithPagination/empty_paginated_result
=== RUN   TestRespondWithPagination/last_page_pagination
--- PASS: TestRespondWithPagination (0.00s)
    --- PASS: TestRespondWithPagination/paginated_list_with_data (0.00s)
    --- PASS: TestRespondWithPagination/empty_paginated_result (0.00s)
    --- PASS: TestRespondWithPagination/last_page_pagination (0.00s)
=== RUN   TestGetUserIDFromContext
=== RUN   TestGetUserIDFromContext/valid_user_ID_exists
=== RUN   TestGetUserIDFromContext/user_ID_not_set
=== RUN   TestGetUserIDFromContext/user_ID_is_not_a_string
=== RUN   TestGetUserIDFromContext/user_ID_is_empty_string
=== RUN   TestGetUserIDFromContext/user_ID_is_nil
--- PASS: TestGetUserIDFromContext (0.00s)
    --- PASS: TestGetUserIDFromContext/valid_user_ID_exists (0.00s)
    --- PASS: TestGetUserIDFromContext/user_ID_not_set (0.00s)
    --- PASS: TestGetUserIDFromContext/user_ID_is_not_a_string (0.00s)
    --- PASS: TestGetUserIDFromContext/user_ID_is_empty_string (0.00s)
    --- PASS: TestGetUserIDFromContext/user_ID_is_nil (0.00s)
=== RUN   TestRequireUserID
=== RUN   TestRequireUserID/valid_user_ID_-_allows_request_to_continue
=== RUN   TestRequireUserID/no_user_ID_-_aborts_request
=== RUN   TestRequireUserID/invalid_user_ID_type_-_aborts_request
=== RUN   TestRequireUserID/empty_user_ID_-_allows_request_to_continue
--- PASS: TestRequireUserID (0.00s)
    --- PASS: TestRequireUserID/valid_user_ID_-_allows_request_to_continue (0.00s)
    --- PASS: TestRequireUserID/no_user_ID_-_aborts_request (0.00s)
    --- PASS: TestRequireUserID/invalid_user_ID_type_-_aborts_request (0.00s)
    --- PASS: TestRequireUserID/empty_user_ID_-_allows_request_to_continue (0.00s)
=== RUN   TestResponseStructures
=== RUN   TestResponseStructures/APIResponse_JSON_serialization
=== RUN   TestResponseStructures/PaginationResponse_JSON_serialization
=== RUN   TestResponseStructures/DataWithPagination_JSON_serialization
--- PASS: TestResponseStructures (0.00s)
    --- PASS: TestResponseStructures/APIResponse_JSON_serialization (0.00s)
    --- PASS: TestResponseStructures/PaginationResponse_JSON_serialization (0.00s)
    --- PASS: TestResponseStructures/DataWithPagination_JSON_serialization (0.00s)
=== RUN   TestTagHandler_Comprehensive_CreateTag
=== RUN   TestTagHandler_Comprehensive_CreateTag/valid_tag_with_all_fields
2025/08/01 15:43:13 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:13 ✅ Container created: 050a7dbb01e0
2025/08/01 15:43:13 🐳 Starting container: 050a7dbb01e0
2025/08/01 15:43:13 ✅ Container started: 050a7dbb01e0
2025/08/01 15:43:13 ⏳ Waiting for container id 050a7dbb01e0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003b73290 Strategies:[0xc00200b6e0 0xc003e460c0]}
2025/08/01 15:43:15 🔔 Container is ready: 050a7dbb01e0
DSN: postgres://testuser:testpass@localhost:36345/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36345/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:15 OK   000001_create_users_table.sql (10.51ms)
2025/08/01 15:43:15 OK   000002_create_projects_table.sql (12.2ms)
2025/08/01 15:43:15 OK   000003_create_log_entries_table.sql (18.74ms)
2025/08/01 15:43:15 OK   000004_create_tags_system.sql (16.38ms)
2025/08/01 15:43:15 OK   000005_create_auth_tables.sql (19.73ms)
2025/08/01 15:43:15 OK   000006_create_insights_table.sql (23.97ms)
2025/08/01 15:43:15 OK   000007_create_performance_indexes.sql (15.33ms)
2025/08/01 15:43:15 OK   000008_create_analytics_views.sql (27.31ms)
2025/08/01 15:43:15 OK   000009_development_data.sql (28.33ms)
2025/08/01 15:43:15 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36345/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36345/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:15 🐳 Stopping container: 050a7dbb01e0
2025/08/01 15:43:16 ✅ Container stopped: 050a7dbb01e0
2025/08/01 15:43:16 🐳 Terminating container: 050a7dbb01e0
2025/08/01 15:43:16 🚫 Container terminated: 050a7dbb01e0
=== RUN   TestTagHandler_Comprehensive_CreateTag/valid_tag_with_minimal_fields
2025/08/01 15:43:16 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:16 ✅ Container created: c1e7e87a95f1
2025/08/01 15:43:16 🐳 Starting container: c1e7e87a95f1
2025/08/01 15:43:16 ✅ Container started: c1e7e87a95f1
2025/08/01 15:43:16 ⏳ Waiting for container id c1e7e87a95f1 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0038bbd50 Strategies:[0xc0039a91a0 0xc00477e1b0]}
2025/08/01 15:43:18 🔔 Container is ready: c1e7e87a95f1
DSN: postgres://testuser:testpass@localhost:36346/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36346/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:18 OK   000001_create_users_table.sql (10.43ms)
2025/08/01 15:43:18 OK   000002_create_projects_table.sql (12.59ms)
2025/08/01 15:43:18 OK   000003_create_log_entries_table.sql (20.54ms)
2025/08/01 15:43:18 OK   000004_create_tags_system.sql (16.16ms)
2025/08/01 15:43:18 OK   000005_create_auth_tables.sql (20.5ms)
2025/08/01 15:43:18 OK   000006_create_insights_table.sql (24.78ms)
2025/08/01 15:43:18 OK   000007_create_performance_indexes.sql (15.98ms)
2025/08/01 15:43:18 OK   000008_create_analytics_views.sql (28.97ms)
2025/08/01 15:43:18 OK   000009_development_data.sql (28.86ms)
2025/08/01 15:43:18 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36346/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36346/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:18 🐳 Stopping container: c1e7e87a95f1
2025/08/01 15:43:18 ✅ Container stopped: c1e7e87a95f1
2025/08/01 15:43:18 🐳 Terminating container: c1e7e87a95f1
2025/08/01 15:43:18 🚫 Container terminated: c1e7e87a95f1
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_very_long_name
2025/08/01 15:43:18 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:18 ✅ Container created: 26c0893cee0a
2025/08/01 15:43:18 🐳 Starting container: 26c0893cee0a
2025/08/01 15:43:19 ✅ Container started: 26c0893cee0a
2025/08/01 15:43:19 ⏳ Waiting for container id 26c0893cee0a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003ecbbe8 Strategies:[0xc003d82e40 0xc004443e30]}
2025/08/01 15:43:20 🔔 Container is ready: 26c0893cee0a
DSN: postgres://testuser:testpass@localhost:36347/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36347/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:21 OK   000001_create_users_table.sql (10.27ms)
2025/08/01 15:43:21 OK   000002_create_projects_table.sql (11.73ms)
2025/08/01 15:43:21 OK   000003_create_log_entries_table.sql (19.02ms)
2025/08/01 15:43:21 OK   000004_create_tags_system.sql (16.23ms)
2025/08/01 15:43:21 OK   000005_create_auth_tables.sql (20.36ms)
2025/08/01 15:43:21 OK   000006_create_insights_table.sql (24.65ms)
2025/08/01 15:43:21 OK   000007_create_performance_indexes.sql (15.22ms)
2025/08/01 15:43:21 OK   000008_create_analytics_views.sql (27.76ms)
2025/08/01 15:43:21 OK   000009_development_data.sql (28.86ms)
2025/08/01 15:43:21 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36347/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36347/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:21 🐳 Stopping container: 26c0893cee0a
2025/08/01 15:43:21 ✅ Container stopped: 26c0893cee0a
2025/08/01 15:43:21 🐳 Terminating container: 26c0893cee0a
2025/08/01 15:43:21 🚫 Container terminated: 26c0893cee0a
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_extremely_long_name_(exceeds_limit)
2025/08/01 15:43:21 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:21 ✅ Container created: b5a2791646e5
2025/08/01 15:43:21 🐳 Starting container: b5a2791646e5
2025/08/01 15:43:21 ✅ Container started: b5a2791646e5
2025/08/01 15:43:21 ⏳ Waiting for container id b5a2791646e5 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0042fd278 Strategies:[0xc0045c6180 0xc0045a7320]}
2025/08/01 15:43:23 🔔 Container is ready: b5a2791646e5
DSN: postgres://testuser:testpass@localhost:36348/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36348/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:23 OK   000001_create_users_table.sql (11.4ms)
2025/08/01 15:43:23 OK   000002_create_projects_table.sql (12.53ms)
2025/08/01 15:43:23 OK   000003_create_log_entries_table.sql (19.75ms)
2025/08/01 15:43:23 OK   000004_create_tags_system.sql (17.12ms)
2025/08/01 15:43:23 OK   000005_create_auth_tables.sql (21.08ms)
2025/08/01 15:43:23 OK   000006_create_insights_table.sql (25.28ms)
2025/08/01 15:43:23 OK   000007_create_performance_indexes.sql (15.54ms)
2025/08/01 15:43:23 OK   000008_create_analytics_views.sql (28.11ms)
2025/08/01 15:43:23 OK   000009_development_data.sql (29.19ms)
2025/08/01 15:43:23 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36348/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36348/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:24 🐳 Stopping container: b5a2791646e5
2025/08/01 15:43:24 ✅ Container stopped: b5a2791646e5
2025/08/01 15:43:24 🐳 Terminating container: b5a2791646e5
2025/08/01 15:43:24 🚫 Container terminated: b5a2791646e5
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_empty_name
2025/08/01 15:43:24 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:24 ✅ Container created: a636a6cfa4c4
2025/08/01 15:43:24 🐳 Starting container: a636a6cfa4c4
2025/08/01 15:43:24 ✅ Container started: a636a6cfa4c4
2025/08/01 15:43:24 ⏳ Waiting for container id a636a6cfa4c4 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006bfd058 Strategies:[0xc006d04840 0xc007170090]}
2025/08/01 15:43:26 🔔 Container is ready: a636a6cfa4c4
DSN: postgres://testuser:testpass@localhost:36349/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36349/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:26 OK   000001_create_users_table.sql (10.08ms)
2025/08/01 15:43:26 OK   000002_create_projects_table.sql (11.69ms)
2025/08/01 15:43:26 OK   000003_create_log_entries_table.sql (18.28ms)
2025/08/01 15:43:26 OK   000004_create_tags_system.sql (15.43ms)
2025/08/01 15:43:26 OK   000005_create_auth_tables.sql (16.71ms)
2025/08/01 15:43:26 OK   000006_create_insights_table.sql (20.91ms)
2025/08/01 15:43:26 OK   000007_create_performance_indexes.sql (13.11ms)
2025/08/01 15:43:26 OK   000008_create_analytics_views.sql (23.15ms)
2025/08/01 15:43:26 OK   000009_development_data.sql (24.21ms)
2025/08/01 15:43:26 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36349/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36349/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:26 🐳 Stopping container: a636a6cfa4c4
2025/08/01 15:43:27 ✅ Container stopped: a636a6cfa4c4
2025/08/01 15:43:27 🐳 Terminating container: a636a6cfa4c4
2025/08/01 15:43:27 🚫 Container terminated: a636a6cfa4c4
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_only_whitespace_name
2025/08/01 15:43:27 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:27 ✅ Container created: c37df431d536
2025/08/01 15:43:27 🐳 Starting container: c37df431d536
2025/08/01 15:43:27 ✅ Container started: c37df431d536
2025/08/01 15:43:27 ⏳ Waiting for container id c37df431d536 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007e0b2a8 Strategies:[0xc00480e7e0 0xc007e96450]}
2025/08/01 15:43:29 🔔 Container is ready: c37df431d536
DSN: postgres://testuser:testpass@localhost:36350/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36350/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:29 OK   000001_create_users_table.sql (10.73ms)
2025/08/01 15:43:29 OK   000002_create_projects_table.sql (12.85ms)
2025/08/01 15:43:29 OK   000003_create_log_entries_table.sql (19.56ms)
2025/08/01 15:43:29 OK   000004_create_tags_system.sql (16.44ms)
2025/08/01 15:43:29 OK   000005_create_auth_tables.sql (20.75ms)
2025/08/01 15:43:29 OK   000006_create_insights_table.sql (24.04ms)
2025/08/01 15:43:29 OK   000007_create_performance_indexes.sql (15.5ms)
2025/08/01 15:43:29 OK   000008_create_analytics_views.sql (28.48ms)
2025/08/01 15:43:29 OK   000009_development_data.sql (28.57ms)
2025/08/01 15:43:29 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36350/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36350/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:29 🐳 Stopping container: c37df431d536
2025/08/01 15:43:29 ✅ Container stopped: c37df431d536
2025/08/01 15:43:29 🐳 Terminating container: c37df431d536
2025/08/01 15:43:30 🚫 Container terminated: c37df431d536
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_invalid_color_format
2025/08/01 15:43:30 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:30 ✅ Container created: 9c0dff8519ad
2025/08/01 15:43:30 🐳 Starting container: 9c0dff8519ad
2025/08/01 15:43:30 ✅ Container started: 9c0dff8519ad
2025/08/01 15:43:30 ⏳ Waiting for container id 9c0dff8519ad image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00590f060 Strategies:[0xc0067aea20 0xc003c4a060]}
2025/08/01 15:43:32 🔔 Container is ready: 9c0dff8519ad
DSN: postgres://testuser:testpass@localhost:36351/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36351/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:32 OK   000001_create_users_table.sql (10.38ms)
2025/08/01 15:43:32 OK   000002_create_projects_table.sql (11.52ms)
2025/08/01 15:43:32 OK   000003_create_log_entries_table.sql (18.67ms)
2025/08/01 15:43:32 OK   000004_create_tags_system.sql (16.33ms)
2025/08/01 15:43:32 OK   000005_create_auth_tables.sql (20.85ms)
2025/08/01 15:43:32 OK   000006_create_insights_table.sql (25.37ms)
2025/08/01 15:43:32 OK   000007_create_performance_indexes.sql (15.27ms)
2025/08/01 15:43:32 OK   000008_create_analytics_views.sql (27.9ms)
2025/08/01 15:43:32 OK   000009_development_data.sql (29.5ms)
2025/08/01 15:43:32 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36351/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36351/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:32 🐳 Stopping container: 9c0dff8519ad
2025/08/01 15:43:32 ✅ Container stopped: 9c0dff8519ad
2025/08/01 15:43:32 🐳 Terminating container: 9c0dff8519ad
2025/08/01 15:43:32 🚫 Container terminated: 9c0dff8519ad
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_special_characters_in_name
2025/08/01 15:43:32 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:32 ✅ Container created: e820cc6776bd
2025/08/01 15:43:32 🐳 Starting container: e820cc6776bd
2025/08/01 15:43:33 ✅ Container started: e820cc6776bd
2025/08/01 15:43:33 ⏳ Waiting for container id e820cc6776bd image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0047e4ae0 Strategies:[0xc005918840 0xc006b2b6b0]}
2025/08/01 15:43:34 🔔 Container is ready: e820cc6776bd
DSN: postgres://testuser:testpass@localhost:36352/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36352/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:34 OK   000001_create_users_table.sql (10.79ms)
2025/08/01 15:43:34 OK   000002_create_projects_table.sql (12.77ms)
2025/08/01 15:43:34 OK   000003_create_log_entries_table.sql (19.4ms)
2025/08/01 15:43:34 OK   000004_create_tags_system.sql (16.31ms)
2025/08/01 15:43:34 OK   000005_create_auth_tables.sql (20.49ms)
2025/08/01 15:43:34 OK   000006_create_insights_table.sql (24.85ms)
2025/08/01 15:43:34 OK   000007_create_performance_indexes.sql (15.55ms)
2025/08/01 15:43:35 OK   000008_create_analytics_views.sql (28.38ms)
2025/08/01 15:43:35 OK   000009_development_data.sql (28.81ms)
2025/08/01 15:43:35 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36352/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36352/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:35 🐳 Stopping container: e820cc6776bd
2025/08/01 15:43:35 ✅ Container stopped: e820cc6776bd
2025/08/01 15:43:35 🐳 Terminating container: e820cc6776bd
2025/08/01 15:43:35 🚫 Container terminated: e820cc6776bd
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_unicode_characters
2025/08/01 15:43:35 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:35 ✅ Container created: 306244a14424
2025/08/01 15:43:35 🐳 Starting container: 306244a14424
2025/08/01 15:43:35 ✅ Container started: 306244a14424
2025/08/01 15:43:35 ⏳ Waiting for container id 306244a14424 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0044589f8 Strategies:[0xc0045a8360 0xc0044d3a40]}
2025/08/01 15:43:37 🔔 Container is ready: 306244a14424
DSN: postgres://testuser:testpass@localhost:36353/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36353/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:37 OK   000001_create_users_table.sql (11.5ms)
2025/08/01 15:43:37 OK   000002_create_projects_table.sql (14.28ms)
2025/08/01 15:43:37 OK   000003_create_log_entries_table.sql (20.93ms)
2025/08/01 15:43:37 OK   000004_create_tags_system.sql (17.03ms)
2025/08/01 15:43:37 OK   000005_create_auth_tables.sql (21.04ms)
2025/08/01 15:43:37 OK   000006_create_insights_table.sql (24.11ms)
2025/08/01 15:43:37 OK   000007_create_performance_indexes.sql (15.44ms)
2025/08/01 15:43:37 OK   000008_create_analytics_views.sql (27.85ms)
2025/08/01 15:43:37 OK   000009_development_data.sql (27.64ms)
2025/08/01 15:43:37 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36353/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36353/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:38 🐳 Stopping container: 306244a14424
2025/08/01 15:43:38 ✅ Container stopped: 306244a14424
2025/08/01 15:43:38 🐳 Terminating container: 306244a14424
2025/08/01 15:43:38 🚫 Container terminated: 306244a14424
=== RUN   TestTagHandler_Comprehensive_CreateTag/tag_with_very_long_description
2025/08/01 15:43:38 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:38 ✅ Container created: 3e6c8c6e760f
2025/08/01 15:43:38 🐳 Starting container: 3e6c8c6e760f
2025/08/01 15:43:38 ✅ Container started: 3e6c8c6e760f
2025/08/01 15:43:38 ⏳ Waiting for container id 3e6c8c6e760f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001eb15f8 Strategies:[0xc001eb3980 0xc00447cd80]}
2025/08/01 15:43:40 🔔 Container is ready: 3e6c8c6e760f
DSN: postgres://testuser:testpass@localhost:36354/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36354/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:40 OK   000001_create_users_table.sql (9.15ms)
2025/08/01 15:43:40 OK   000002_create_projects_table.sql (11.21ms)
2025/08/01 15:43:40 OK   000003_create_log_entries_table.sql (18.1ms)
2025/08/01 15:43:40 OK   000004_create_tags_system.sql (15.6ms)
2025/08/01 15:43:40 OK   000005_create_auth_tables.sql (20.22ms)
2025/08/01 15:43:40 OK   000006_create_insights_table.sql (24.07ms)
2025/08/01 15:43:40 OK   000007_create_performance_indexes.sql (14.97ms)
2025/08/01 15:43:40 OK   000008_create_analytics_views.sql (27.58ms)
2025/08/01 15:43:40 OK   000009_development_data.sql (28.28ms)
2025/08/01 15:43:40 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36354/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36354/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:40 🐳 Stopping container: 3e6c8c6e760f
2025/08/01 15:43:41 ✅ Container stopped: 3e6c8c6e760f
2025/08/01 15:43:41 🐳 Terminating container: 3e6c8c6e760f
2025/08/01 15:43:41 🚫 Container terminated: 3e6c8c6e760f
--- PASS: TestTagHandler_Comprehensive_CreateTag (27.86s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/valid_tag_with_all_fields (2.80s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/valid_tag_with_minimal_fields (2.84s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_very_long_name (2.77s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_extremely_long_name_(exceeds_limit) (2.70s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_empty_name (2.79s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_only_whitespace_name (2.79s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_invalid_color_format (2.82s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_special_characters_in_name (2.78s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_unicode_characters (2.77s)
    --- PASS: TestTagHandler_Comprehensive_CreateTag/tag_with_very_long_description (2.80s)
=== RUN   TestTagHandler_Comprehensive_TagDuplication
=== RUN   TestTagHandler_Comprehensive_TagDuplication/duplicate_tag_name_should_fail
2025/08/01 15:43:41 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:41 ✅ Container created: 1eed4f0f3a13
2025/08/01 15:43:41 🐳 Starting container: 1eed4f0f3a13
2025/08/01 15:43:41 ✅ Container started: 1eed4f0f3a13
2025/08/01 15:43:41 ⏳ Waiting for container id 1eed4f0f3a13 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc005ab6cf8 Strategies:[0xc004784f60 0xc005aa79b0]}
2025/08/01 15:43:43 🔔 Container is ready: 1eed4f0f3a13
DSN: postgres://testuser:testpass@localhost:36355/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36355/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:43 OK   000001_create_users_table.sql (10.96ms)
2025/08/01 15:43:43 OK   000002_create_projects_table.sql (13.22ms)
2025/08/01 15:43:43 OK   000003_create_log_entries_table.sql (19.39ms)
2025/08/01 15:43:43 OK   000004_create_tags_system.sql (16.44ms)
2025/08/01 15:43:43 OK   000005_create_auth_tables.sql (20.07ms)
2025/08/01 15:43:43 OK   000006_create_insights_table.sql (24.26ms)
2025/08/01 15:43:43 OK   000007_create_performance_indexes.sql (15.73ms)
2025/08/01 15:43:43 OK   000008_create_analytics_views.sql (28.01ms)
2025/08/01 15:43:43 OK   000009_development_data.sql (28.35ms)
2025/08/01 15:43:43 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36355/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36355/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:43 🐳 Stopping container: 1eed4f0f3a13
2025/08/01 15:43:43 ✅ Container stopped: 1eed4f0f3a13
2025/08/01 15:43:43 🐳 Terminating container: 1eed4f0f3a13
2025/08/01 15:43:44 🚫 Container terminated: 1eed4f0f3a13
=== RUN   TestTagHandler_Comprehensive_TagDuplication/case_insensitive_tag_name_should_fail
2025/08/01 15:43:44 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:44 ✅ Container created: f067410429e7
2025/08/01 15:43:44 🐳 Starting container: f067410429e7
2025/08/01 15:43:44 ✅ Container started: f067410429e7
2025/08/01 15:43:44 ⏳ Waiting for container id f067410429e7 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00468d070 Strategies:[0xc0005a91a0 0xc006a861e0]}
2025/08/01 15:43:46 🔔 Container is ready: f067410429e7
DSN: postgres://testuser:testpass@localhost:36356/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36356/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:46 OK   000001_create_users_table.sql (10.1ms)
2025/08/01 15:43:46 OK   000002_create_projects_table.sql (12.03ms)
2025/08/01 15:43:46 OK   000003_create_log_entries_table.sql (18.96ms)
2025/08/01 15:43:46 OK   000004_create_tags_system.sql (16.18ms)
2025/08/01 15:43:46 OK   000005_create_auth_tables.sql (20.1ms)
2025/08/01 15:43:46 OK   000006_create_insights_table.sql (24.34ms)
2025/08/01 15:43:46 OK   000007_create_performance_indexes.sql (15.8ms)
2025/08/01 15:43:46 OK   000008_create_analytics_views.sql (28.4ms)
2025/08/01 15:43:46 OK   000009_development_data.sql (28.86ms)
2025/08/01 15:43:46 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36356/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36356/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:46 🐳 Stopping container: f067410429e7
2025/08/01 15:43:46 ✅ Container stopped: f067410429e7
2025/08/01 15:43:46 🐳 Terminating container: f067410429e7
2025/08/01 15:43:46 🚫 Container terminated: f067410429e7
--- PASS: TestTagHandler_Comprehensive_TagDuplication (5.65s)
    --- PASS: TestTagHandler_Comprehensive_TagDuplication/duplicate_tag_name_should_fail (2.83s)
    --- PASS: TestTagHandler_Comprehensive_TagDuplication/case_insensitive_tag_name_should_fail (2.82s)
=== RUN   TestTagHandler_Comprehensive_GetTags
=== RUN   TestTagHandler_Comprehensive_GetTags/get_all_tags_with_multiple_tags
2025/08/01 15:43:46 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:46 ✅ Container created: 8cd4191dfead
2025/08/01 15:43:46 🐳 Starting container: 8cd4191dfead
2025/08/01 15:43:47 ✅ Container started: 8cd4191dfead
2025/08/01 15:43:47 ⏳ Waiting for container id 8cd4191dfead image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0047d9660 Strategies:[0xc0046bd320 0xc003adb710]}
2025/08/01 15:43:48 🔔 Container is ready: 8cd4191dfead
DSN: postgres://testuser:testpass@localhost:36357/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36357/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:48 OK   000001_create_users_table.sql (9.18ms)
2025/08/01 15:43:48 OK   000002_create_projects_table.sql (11.16ms)
2025/08/01 15:43:48 OK   000003_create_log_entries_table.sql (18.39ms)
2025/08/01 15:43:48 OK   000004_create_tags_system.sql (15.54ms)
2025/08/01 15:43:48 OK   000005_create_auth_tables.sql (20.09ms)
2025/08/01 15:43:48 OK   000006_create_insights_table.sql (24.2ms)
2025/08/01 15:43:48 OK   000007_create_performance_indexes.sql (15.23ms)
2025/08/01 15:43:49 OK   000008_create_analytics_views.sql (28.19ms)
2025/08/01 15:43:49 OK   000009_development_data.sql (28.51ms)
2025/08/01 15:43:49 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36357/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36357/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:49 🐳 Stopping container: 8cd4191dfead
2025/08/01 15:43:49 ✅ Container stopped: 8cd4191dfead
2025/08/01 15:43:49 🐳 Terminating container: 8cd4191dfead
2025/08/01 15:43:49 🚫 Container terminated: 8cd4191dfead
=== RUN   TestTagHandler_Comprehensive_GetTags/get_tags_with_seeded_database
2025/08/01 15:43:49 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:49 ✅ Container created: 3890c3746989
2025/08/01 15:43:49 🐳 Starting container: 3890c3746989
2025/08/01 15:43:49 ✅ Container started: 3890c3746989
2025/08/01 15:43:49 ⏳ Waiting for container id 3890c3746989 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0024d9f30 Strategies:[0xc006a5ed20 0xc00641ac00]}
2025/08/01 15:43:51 🔔 Container is ready: 3890c3746989
DSN: postgres://testuser:testpass@localhost:36358/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36358/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:51 OK   000001_create_users_table.sql (10.86ms)
2025/08/01 15:43:51 OK   000002_create_projects_table.sql (13.09ms)
2025/08/01 15:43:51 OK   000003_create_log_entries_table.sql (20.52ms)
2025/08/01 15:43:51 OK   000004_create_tags_system.sql (16.36ms)
2025/08/01 15:43:51 OK   000005_create_auth_tables.sql (20.57ms)
2025/08/01 15:43:51 OK   000006_create_insights_table.sql (25.15ms)
2025/08/01 15:43:51 OK   000007_create_performance_indexes.sql (15.33ms)
2025/08/01 15:43:51 OK   000008_create_analytics_views.sql (28.02ms)
2025/08/01 15:43:51 OK   000009_development_data.sql (28.47ms)
2025/08/01 15:43:51 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36358/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36358/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:52 🐳 Stopping container: 3890c3746989
2025/08/01 15:43:52 ✅ Container stopped: 3890c3746989
2025/08/01 15:43:52 🐳 Terminating container: 3890c3746989
2025/08/01 15:43:52 🚫 Container terminated: 3890c3746989
--- PASS: TestTagHandler_Comprehensive_GetTags (5.61s)
    --- PASS: TestTagHandler_Comprehensive_GetTags/get_all_tags_with_multiple_tags (2.79s)
    --- PASS: TestTagHandler_Comprehensive_GetTags/get_tags_with_seeded_database (2.81s)
=== RUN   TestTagHandler_Comprehensive_TagSearch
=== RUN   TestTagHandler_Comprehensive_TagSearch/valid_search_query
2025/08/01 15:43:52 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:52 ✅ Container created: c7af4fa12611
2025/08/01 15:43:52 🐳 Starting container: c7af4fa12611
2025/08/01 15:43:52 ✅ Container started: c7af4fa12611
2025/08/01 15:43:52 ⏳ Waiting for container id c7af4fa12611 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004246208 Strategies:[0xc001e9c7e0 0xc0044433b0]}
2025/08/01 15:43:54 🔔 Container is ready: c7af4fa12611
DSN: postgres://testuser:testpass@localhost:36359/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36359/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:54 OK   000001_create_users_table.sql (10.73ms)
2025/08/01 15:43:54 OK   000002_create_projects_table.sql (12.41ms)
2025/08/01 15:43:54 OK   000003_create_log_entries_table.sql (19.53ms)
2025/08/01 15:43:54 OK   000004_create_tags_system.sql (16.5ms)
2025/08/01 15:43:54 OK   000005_create_auth_tables.sql (21.21ms)
2025/08/01 15:43:54 OK   000006_create_insights_table.sql (25.01ms)
2025/08/01 15:43:54 OK   000007_create_performance_indexes.sql (15.68ms)
2025/08/01 15:43:54 OK   000008_create_analytics_views.sql (28.5ms)
2025/08/01 15:43:54 OK   000009_development_data.sql (28.83ms)
2025/08/01 15:43:54 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36359/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36359/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:54 🐳 Stopping container: c7af4fa12611
2025/08/01 15:43:55 ✅ Container stopped: c7af4fa12611
2025/08/01 15:43:55 🐳 Terminating container: c7af4fa12611
2025/08/01 15:43:55 🚫 Container terminated: c7af4fa12611
=== RUN   TestTagHandler_Comprehensive_TagSearch/search_with_limit
2025/08/01 15:43:55 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:55 ✅ Container created: 9122f237799d
2025/08/01 15:43:55 🐳 Starting container: 9122f237799d
2025/08/01 15:43:55 ✅ Container started: 9122f237799d
2025/08/01 15:43:55 ⏳ Waiting for container id 9122f237799d image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007dbb100 Strategies:[0xc00548a300 0xc00466e4e0]}
2025/08/01 15:43:57 🔔 Container is ready: 9122f237799d
DSN: postgres://testuser:testpass@localhost:36360/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36360/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:43:57 OK   000001_create_users_table.sql (10.73ms)
2025/08/01 15:43:57 OK   000002_create_projects_table.sql (12.7ms)
2025/08/01 15:43:57 OK   000003_create_log_entries_table.sql (19.73ms)
2025/08/01 15:43:57 OK   000004_create_tags_system.sql (16.83ms)
2025/08/01 15:43:57 OK   000005_create_auth_tables.sql (21.03ms)
2025/08/01 15:43:57 OK   000006_create_insights_table.sql (25.3ms)
2025/08/01 15:43:57 OK   000007_create_performance_indexes.sql (16.17ms)
2025/08/01 15:43:57 OK   000008_create_analytics_views.sql (28.71ms)
2025/08/01 15:43:57 OK   000009_development_data.sql (29.49ms)
2025/08/01 15:43:57 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36360/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36360/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:43:57 🐳 Stopping container: 9122f237799d
2025/08/01 15:43:57 ✅ Container stopped: 9122f237799d
2025/08/01 15:43:57 🐳 Terminating container: 9122f237799d
2025/08/01 15:43:58 🚫 Container terminated: 9122f237799d
=== RUN   TestTagHandler_Comprehensive_TagSearch/search_with_maximum_limit
2025/08/01 15:43:58 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:43:58 ✅ Container created: 8b9932060a56
2025/08/01 15:43:58 🐳 Starting container: 8b9932060a56
2025/08/01 15:43:58 ✅ Container started: 8b9932060a56
2025/08/01 15:43:58 ⏳ Waiting for container id 8b9932060a56 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006ad22c8 Strategies:[0xc006b16fc0 0xc006485ad0]}
2025/08/01 15:44:00 🔔 Container is ready: 8b9932060a56
DSN: postgres://testuser:testpass@localhost:36361/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36361/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:00 OK   000001_create_users_table.sql (9.74ms)
2025/08/01 15:44:00 OK   000002_create_projects_table.sql (11.71ms)
2025/08/01 15:44:00 OK   000003_create_log_entries_table.sql (18.23ms)
2025/08/01 15:44:00 OK   000004_create_tags_system.sql (16.27ms)
2025/08/01 15:44:00 OK   000005_create_auth_tables.sql (20.36ms)
2025/08/01 15:44:00 OK   000006_create_insights_table.sql (24.15ms)
2025/08/01 15:44:00 OK   000007_create_performance_indexes.sql (15.2ms)
2025/08/01 15:44:00 OK   000008_create_analytics_views.sql (28.2ms)
2025/08/01 15:44:00 OK   000009_development_data.sql (28.43ms)
2025/08/01 15:44:00 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36361/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36361/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:00 🐳 Stopping container: 8b9932060a56
2025/08/01 15:44:00 ✅ Container stopped: 8b9932060a56
2025/08/01 15:44:00 🐳 Terminating container: 8b9932060a56
2025/08/01 15:44:00 🚫 Container terminated: 8b9932060a56
=== RUN   TestTagHandler_Comprehensive_TagSearch/search_with_excessive_limit
2025/08/01 15:44:00 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:00 ✅ Container created: ea6daf2a9102
2025/08/01 15:44:00 🐳 Starting container: ea6daf2a9102
2025/08/01 15:44:01 ✅ Container started: ea6daf2a9102
2025/08/01 15:44:01 ⏳ Waiting for container id ea6daf2a9102 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006ad2c98 Strategies:[0xc005a7eae0 0xc00454a330]}
2025/08/01 15:44:02 🔔 Container is ready: ea6daf2a9102
DSN: postgres://testuser:testpass@localhost:36362/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36362/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:02 OK   000001_create_users_table.sql (9.86ms)
2025/08/01 15:44:02 OK   000002_create_projects_table.sql (12.58ms)
2025/08/01 15:44:02 OK   000003_create_log_entries_table.sql (18.74ms)
2025/08/01 15:44:02 OK   000004_create_tags_system.sql (16.38ms)
2025/08/01 15:44:03 OK   000005_create_auth_tables.sql (20.18ms)
2025/08/01 15:44:03 OK   000006_create_insights_table.sql (24.26ms)
2025/08/01 15:44:03 OK   000007_create_performance_indexes.sql (15.47ms)
2025/08/01 15:44:03 OK   000008_create_analytics_views.sql (27.83ms)
2025/08/01 15:44:03 OK   000009_development_data.sql (28.31ms)
2025/08/01 15:44:03 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36362/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36362/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:03 🐳 Stopping container: ea6daf2a9102
2025/08/01 15:44:03 ✅ Container stopped: ea6daf2a9102
2025/08/01 15:44:03 🐳 Terminating container: ea6daf2a9102
2025/08/01 15:44:03 🚫 Container terminated: ea6daf2a9102
=== RUN   TestTagHandler_Comprehensive_TagSearch/search_with_invalid_limit
2025/08/01 15:44:03 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:03 ✅ Container created: 70776d193b67
2025/08/01 15:44:03 🐳 Starting container: 70776d193b67
2025/08/01 15:44:03 ✅ Container started: 70776d193b67
2025/08/01 15:44:03 ⏳ Waiting for container id 70776d193b67 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007be1790 Strategies:[0xc001eb34a0 0xc006d31ce0]}
2025/08/01 15:44:05 🔔 Container is ready: 70776d193b67
DSN: postgres://testuser:testpass@localhost:36363/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36363/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:05 OK   000001_create_users_table.sql (10.39ms)
2025/08/01 15:44:05 OK   000002_create_projects_table.sql (12.12ms)
2025/08/01 15:44:05 OK   000003_create_log_entries_table.sql (19.11ms)
2025/08/01 15:44:05 OK   000004_create_tags_system.sql (16.1ms)
2025/08/01 15:44:05 OK   000005_create_auth_tables.sql (20.02ms)
2025/08/01 15:44:05 OK   000006_create_insights_table.sql (25.15ms)
2025/08/01 15:44:05 OK   000007_create_performance_indexes.sql (13.99ms)
2025/08/01 15:44:05 OK   000008_create_analytics_views.sql (24.01ms)
2025/08/01 15:44:05 OK   000009_development_data.sql (24.25ms)
2025/08/01 15:44:05 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36363/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36363/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:06 🐳 Stopping container: 70776d193b67
2025/08/01 15:44:06 ✅ Container stopped: 70776d193b67
2025/08/01 15:44:06 🐳 Terminating container: 70776d193b67
2025/08/01 15:44:06 🚫 Container terminated: 70776d193b67
=== RUN   TestTagHandler_Comprehensive_TagSearch/search_with_negative_limit
2025/08/01 15:44:06 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:06 ✅ Container created: 1faa6fa2433c
2025/08/01 15:44:06 🐳 Starting container: 1faa6fa2433c
2025/08/01 15:44:06 ✅ Container started: 1faa6fa2433c
2025/08/01 15:44:06 ⏳ Waiting for container id 1faa6fa2433c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00478da88 Strategies:[0xc0069bafc0 0xc003c4ab40]}
2025/08/01 15:44:08 🔔 Container is ready: 1faa6fa2433c
DSN: postgres://testuser:testpass@localhost:36364/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36364/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:08 OK   000001_create_users_table.sql (10.52ms)
2025/08/01 15:44:08 OK   000002_create_projects_table.sql (12.18ms)
2025/08/01 15:44:08 OK   000003_create_log_entries_table.sql (20.29ms)
2025/08/01 15:44:08 OK   000004_create_tags_system.sql (17.97ms)
2025/08/01 15:44:08 OK   000005_create_auth_tables.sql (20.31ms)
2025/08/01 15:44:08 OK   000006_create_insights_table.sql (24.28ms)
2025/08/01 15:44:08 OK   000007_create_performance_indexes.sql (15.64ms)
2025/08/01 15:44:08 OK   000008_create_analytics_views.sql (28.2ms)
2025/08/01 15:44:08 OK   000009_development_data.sql (28.83ms)
2025/08/01 15:44:08 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36364/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36364/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:08 🐳 Stopping container: 1faa6fa2433c
2025/08/01 15:44:09 ✅ Container stopped: 1faa6fa2433c
2025/08/01 15:44:09 🐳 Terminating container: 1faa6fa2433c
2025/08/01 15:44:09 🚫 Container terminated: 1faa6fa2433c
=== RUN   TestTagHandler_Comprehensive_TagSearch/empty_search_query
2025/08/01 15:44:09 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:09 ✅ Container created: 366fa32f4334
2025/08/01 15:44:09 🐳 Starting container: 366fa32f4334
2025/08/01 15:44:09 ✅ Container started: 366fa32f4334
2025/08/01 15:44:09 ⏳ Waiting for container id 366fa32f4334 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007e0af60 Strategies:[0xc00200be60 0xc004447a10]}
2025/08/01 15:44:11 🔔 Container is ready: 366fa32f4334
DSN: postgres://testuser:testpass@localhost:36365/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36365/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:11 OK   000001_create_users_table.sql (10.93ms)
2025/08/01 15:44:11 OK   000002_create_projects_table.sql (13.13ms)
2025/08/01 15:44:11 OK   000003_create_log_entries_table.sql (20.41ms)
2025/08/01 15:44:11 OK   000004_create_tags_system.sql (17.35ms)
2025/08/01 15:44:11 OK   000005_create_auth_tables.sql (22.86ms)
2025/08/01 15:44:11 OK   000006_create_insights_table.sql (25.55ms)
2025/08/01 15:44:11 OK   000007_create_performance_indexes.sql (15.95ms)
2025/08/01 15:44:11 OK   000008_create_analytics_views.sql (29.59ms)
2025/08/01 15:44:11 OK   000009_development_data.sql (31.17ms)
2025/08/01 15:44:11 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36365/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36365/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:11 🐳 Stopping container: 366fa32f4334
2025/08/01 15:44:11 ✅ Container stopped: 366fa32f4334
2025/08/01 15:44:11 🐳 Terminating container: 366fa32f4334
2025/08/01 15:44:12 🚫 Container terminated: 366fa32f4334
=== RUN   TestTagHandler_Comprehensive_TagSearch/whitespace_search_query
2025/08/01 15:44:12 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:12 ✅ Container created: fe2715732524
2025/08/01 15:44:12 🐳 Starting container: fe2715732524
2025/08/01 15:44:12 ✅ Container started: fe2715732524
2025/08/01 15:44:12 ⏳ Waiting for container id fe2715732524 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006706e78 Strategies:[0xc000709c80 0xc001dbb860]}
2025/08/01 15:44:14 🔔 Container is ready: fe2715732524
DSN: postgres://testuser:testpass@localhost:36366/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36366/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:14 OK   000001_create_users_table.sql (10.05ms)
2025/08/01 15:44:14 OK   000002_create_projects_table.sql (11.87ms)
2025/08/01 15:44:14 OK   000003_create_log_entries_table.sql (19.64ms)
2025/08/01 15:44:14 OK   000004_create_tags_system.sql (16.28ms)
2025/08/01 15:44:14 OK   000005_create_auth_tables.sql (22.47ms)
2025/08/01 15:44:14 OK   000006_create_insights_table.sql (24.29ms)
2025/08/01 15:44:14 OK   000007_create_performance_indexes.sql (15.56ms)
2025/08/01 15:44:14 OK   000008_create_analytics_views.sql (28.37ms)
2025/08/01 15:44:14 OK   000009_development_data.sql (28.49ms)
2025/08/01 15:44:14 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36366/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36366/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:14 🐳 Stopping container: fe2715732524
2025/08/01 15:44:14 ✅ Container stopped: fe2715732524
2025/08/01 15:44:14 🐳 Terminating container: fe2715732524
2025/08/01 15:44:14 🚫 Container terminated: fe2715732524
=== RUN   TestTagHandler_Comprehensive_TagSearch/special_characters_search
2025/08/01 15:44:14 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:14 ✅ Container created: e96edafa9641
2025/08/01 15:44:14 🐳 Starting container: e96edafa9641
2025/08/01 15:44:15 ✅ Container started: e96edafa9641
2025/08/01 15:44:15 ⏳ Waiting for container id e96edafa9641 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc005ab73b8 Strategies:[0xc0069209c0 0xc0054dc720]}
2025/08/01 15:44:16 🔔 Container is ready: e96edafa9641
DSN: postgres://testuser:testpass@localhost:36367/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36367/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:17 OK   000001_create_users_table.sql (11.19ms)
2025/08/01 15:44:17 OK   000002_create_projects_table.sql (12.81ms)
2025/08/01 15:44:17 OK   000003_create_log_entries_table.sql (19.24ms)
2025/08/01 15:44:17 OK   000004_create_tags_system.sql (15.85ms)
2025/08/01 15:44:17 OK   000005_create_auth_tables.sql (19.77ms)
2025/08/01 15:44:17 OK   000006_create_insights_table.sql (24.05ms)
2025/08/01 15:44:17 OK   000007_create_performance_indexes.sql (15.26ms)
2025/08/01 15:44:17 OK   000008_create_analytics_views.sql (28.24ms)
2025/08/01 15:44:17 OK   000009_development_data.sql (28.85ms)
2025/08/01 15:44:17 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36367/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36367/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:17 🐳 Stopping container: e96edafa9641
2025/08/01 15:44:17 ✅ Container stopped: e96edafa9641
2025/08/01 15:44:17 🐳 Terminating container: e96edafa9641
2025/08/01 15:44:17 🚫 Container terminated: e96edafa9641
=== RUN   TestTagHandler_Comprehensive_TagSearch/unicode_search
2025/08/01 15:44:17 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:17 ✅ Container created: 63618d6d8af6
2025/08/01 15:44:17 🐳 Starting container: 63618d6d8af6
2025/08/01 15:44:17 ✅ Container started: 63618d6d8af6
2025/08/01 15:44:17 ⏳ Waiting for container id 63618d6d8af6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b8a9b8 Strategies:[0xc004401740 0xc00485eb10]}
2025/08/01 15:44:19 🔔 Container is ready: 63618d6d8af6
DSN: postgres://testuser:testpass@localhost:36368/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36368/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:19 OK   000001_create_users_table.sql (9.44ms)
2025/08/01 15:44:19 OK   000002_create_projects_table.sql (11.29ms)
2025/08/01 15:44:19 OK   000003_create_log_entries_table.sql (18.94ms)
2025/08/01 15:44:19 OK   000004_create_tags_system.sql (16.23ms)
2025/08/01 15:44:19 OK   000005_create_auth_tables.sql (20.43ms)
2025/08/01 15:44:19 OK   000006_create_insights_table.sql (24.71ms)
2025/08/01 15:44:19 OK   000007_create_performance_indexes.sql (16.11ms)
2025/08/01 15:44:19 OK   000008_create_analytics_views.sql (29.19ms)
2025/08/01 15:44:19 OK   000009_development_data.sql (29.86ms)
2025/08/01 15:44:19 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36368/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36368/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:20 🐳 Stopping container: 63618d6d8af6
2025/08/01 15:44:20 ✅ Container stopped: 63618d6d8af6
2025/08/01 15:44:20 🐳 Terminating container: 63618d6d8af6
2025/08/01 15:44:20 🚫 Container terminated: 63618d6d8af6
=== RUN   TestTagHandler_Comprehensive_TagSearch/very_long_search_query
2025/08/01 15:44:20 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:20 ✅ Container created: f89d94ad40b6
2025/08/01 15:44:20 🐳 Starting container: f89d94ad40b6
2025/08/01 15:44:20 ✅ Container started: f89d94ad40b6
2025/08/01 15:44:20 ⏳ Waiting for container id f89d94ad40b6 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0069e91e8 Strategies:[0xc006d04840 0xc003e6c630]}
2025/08/01 15:44:22 🔔 Container is ready: f89d94ad40b6
DSN: postgres://testuser:testpass@localhost:36369/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36369/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:22 OK   000001_create_users_table.sql (11.32ms)
2025/08/01 15:44:22 OK   000002_create_projects_table.sql (12.69ms)
2025/08/01 15:44:22 OK   000003_create_log_entries_table.sql (20.11ms)
2025/08/01 15:44:22 OK   000004_create_tags_system.sql (16.64ms)
2025/08/01 15:44:22 OK   000005_create_auth_tables.sql (20.18ms)
2025/08/01 15:44:22 OK   000006_create_insights_table.sql (23.79ms)
2025/08/01 15:44:22 OK   000007_create_performance_indexes.sql (14.86ms)
2025/08/01 15:44:22 OK   000008_create_analytics_views.sql (27.44ms)
2025/08/01 15:44:22 OK   000009_development_data.sql (27.42ms)
2025/08/01 15:44:22 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36369/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36369/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:22 🐳 Stopping container: f89d94ad40b6
2025/08/01 15:44:23 ✅ Container stopped: f89d94ad40b6
2025/08/01 15:44:23 🐳 Terminating container: f89d94ad40b6
2025/08/01 15:44:23 🚫 Container terminated: f89d94ad40b6
--- PASS: TestTagHandler_Comprehensive_TagSearch (30.86s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/valid_search_query (2.79s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/search_with_limit (2.85s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/search_with_maximum_limit (2.81s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/search_with_excessive_limit (2.85s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/search_with_invalid_limit (2.74s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/search_with_negative_limit (2.77s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/empty_search_query (2.79s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/whitespace_search_query (2.80s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/special_characters_search (2.86s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/unicode_search (2.78s)
    --- PASS: TestTagHandler_Comprehensive_TagSearch/very_long_search_query (2.82s)
=== RUN   TestTagHandler_Comprehensive_UpdateTag
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_all_fields
2025/08/01 15:44:23 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:23 ✅ Container created: ff68efb53f14
2025/08/01 15:44:23 🐳 Starting container: ff68efb53f14
2025/08/01 15:44:23 ✅ Container started: ff68efb53f14
2025/08/01 15:44:23 ⏳ Waiting for container id ff68efb53f14 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0071881e0 Strategies:[0xc008403b60 0xc00215e7e0]}
2025/08/01 15:44:25 🔔 Container is ready: ff68efb53f14
DSN: postgres://testuser:testpass@localhost:36370/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36370/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:25 OK   000001_create_users_table.sql (11.11ms)
2025/08/01 15:44:25 OK   000002_create_projects_table.sql (12.91ms)
2025/08/01 15:44:25 OK   000003_create_log_entries_table.sql (20.68ms)
2025/08/01 15:44:25 OK   000004_create_tags_system.sql (17.54ms)
2025/08/01 15:44:25 OK   000005_create_auth_tables.sql (21.26ms)
2025/08/01 15:44:25 OK   000006_create_insights_table.sql (25.06ms)
2025/08/01 15:44:25 OK   000007_create_performance_indexes.sql (15.85ms)
2025/08/01 15:44:25 OK   000008_create_analytics_views.sql (28.67ms)
2025/08/01 15:44:25 OK   000009_development_data.sql (28.57ms)
2025/08/01 15:44:25 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36370/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36370/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:25 🐳 Stopping container: ff68efb53f14
2025/08/01 15:44:26 ✅ Container stopped: ff68efb53f14
2025/08/01 15:44:26 🐳 Terminating container: ff68efb53f14
2025/08/01 15:44:26 🚫 Container terminated: ff68efb53f14
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_only_name
2025/08/01 15:44:26 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:26 ✅ Container created: 3328253a5cbf
2025/08/01 15:44:26 🐳 Starting container: 3328253a5cbf
2025/08/01 15:44:26 ✅ Container started: 3328253a5cbf
2025/08/01 15:44:26 ⏳ Waiting for container id 3328253a5cbf image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003e1d0d0 Strategies:[0xc005918960 0xc003bb1d10]}
2025/08/01 15:44:28 🔔 Container is ready: 3328253a5cbf
DSN: postgres://testuser:testpass@localhost:36371/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36371/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:28 OK   000001_create_users_table.sql (10.03ms)
2025/08/01 15:44:28 OK   000002_create_projects_table.sql (11.64ms)
2025/08/01 15:44:28 OK   000003_create_log_entries_table.sql (18.6ms)
2025/08/01 15:44:28 OK   000004_create_tags_system.sql (16.57ms)
2025/08/01 15:44:28 OK   000005_create_auth_tables.sql (20.48ms)
2025/08/01 15:44:28 OK   000006_create_insights_table.sql (24.59ms)
2025/08/01 15:44:28 OK   000007_create_performance_indexes.sql (15.47ms)
2025/08/01 15:44:28 OK   000008_create_analytics_views.sql (28.72ms)
2025/08/01 15:44:28 OK   000009_development_data.sql (28.63ms)
2025/08/01 15:44:28 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36371/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36371/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:28 🐳 Stopping container: 3328253a5cbf
2025/08/01 15:44:28 ✅ Container stopped: 3328253a5cbf
2025/08/01 15:44:28 🐳 Terminating container: 3328253a5cbf
2025/08/01 15:44:28 🚫 Container terminated: 3328253a5cbf
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_only_color
2025/08/01 15:44:28 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:29 ✅ Container created: 2c0d7ff44540
2025/08/01 15:44:29 🐳 Starting container: 2c0d7ff44540
2025/08/01 15:44:29 ✅ Container started: 2c0d7ff44540
2025/08/01 15:44:29 ⏳ Waiting for container id 2c0d7ff44540 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b14058 Strategies:[0xc007e4a120 0xc002760720]}
2025/08/01 15:44:30 🔔 Container is ready: 2c0d7ff44540
DSN: postgres://testuser:testpass@localhost:36372/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36372/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:31 OK   000001_create_users_table.sql (10.59ms)
2025/08/01 15:44:31 OK   000002_create_projects_table.sql (12.28ms)
2025/08/01 15:44:31 OK   000003_create_log_entries_table.sql (19.16ms)
2025/08/01 15:44:31 OK   000004_create_tags_system.sql (16.45ms)
2025/08/01 15:44:31 OK   000005_create_auth_tables.sql (20.37ms)
2025/08/01 15:44:31 OK   000006_create_insights_table.sql (24.33ms)
2025/08/01 15:44:31 OK   000007_create_performance_indexes.sql (15.52ms)
2025/08/01 15:44:31 OK   000008_create_analytics_views.sql (28.11ms)
2025/08/01 15:44:31 OK   000009_development_data.sql (28.73ms)
2025/08/01 15:44:31 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36372/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36372/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:31 🐳 Stopping container: 2c0d7ff44540
2025/08/01 15:44:31 ✅ Container stopped: 2c0d7ff44540
2025/08/01 15:44:31 🐳 Terminating container: 2c0d7ff44540
2025/08/01 15:44:31 🚫 Container terminated: 2c0d7ff44540
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_with_empty_name
2025/08/01 15:44:31 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:31 ✅ Container created: da475489d4e5
2025/08/01 15:44:31 🐳 Starting container: da475489d4e5
2025/08/01 15:44:32 ✅ Container started: da475489d4e5
2025/08/01 15:44:32 ⏳ Waiting for container id da475489d4e5 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004247a88 Strategies:[0xc0067af5c0 0xc00530d080]}
2025/08/01 15:44:33 🔔 Container is ready: da475489d4e5
DSN: postgres://testuser:testpass@localhost:36373/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36373/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:33 OK   000001_create_users_table.sql (10.63ms)
2025/08/01 15:44:33 OK   000002_create_projects_table.sql (12.76ms)
2025/08/01 15:44:33 OK   000003_create_log_entries_table.sql (19.59ms)
2025/08/01 15:44:33 OK   000004_create_tags_system.sql (17.14ms)
2025/08/01 15:44:33 OK   000005_create_auth_tables.sql (20.63ms)
2025/08/01 15:44:33 OK   000006_create_insights_table.sql (24.58ms)
2025/08/01 15:44:33 OK   000007_create_performance_indexes.sql (16.29ms)
2025/08/01 15:44:34 OK   000008_create_analytics_views.sql (28.46ms)
2025/08/01 15:44:34 OK   000009_development_data.sql (29.16ms)
2025/08/01 15:44:34 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36373/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36373/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:34 🐳 Stopping container: da475489d4e5
2025/08/01 15:44:34 ✅ Container stopped: da475489d4e5
2025/08/01 15:44:34 🐳 Terminating container: da475489d4e5
2025/08/01 15:44:34 🚫 Container terminated: da475489d4e5
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_with_invalid_color
2025/08/01 15:44:34 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:34 ✅ Container created: 38b573ccad62
2025/08/01 15:44:34 🐳 Starting container: 38b573ccad62
2025/08/01 15:44:34 ✅ Container started: 38b573ccad62
2025/08/01 15:44:34 ⏳ Waiting for container id 38b573ccad62 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0054cfda8 Strategies:[0xc004731380 0xc006441bc0]}
2025/08/01 15:44:36 🔔 Container is ready: 38b573ccad62
DSN: postgres://testuser:testpass@localhost:36374/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36374/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:36 OK   000001_create_users_table.sql (10.74ms)
2025/08/01 15:44:36 OK   000002_create_projects_table.sql (12.31ms)
2025/08/01 15:44:36 OK   000003_create_log_entries_table.sql (19.27ms)
2025/08/01 15:44:36 OK   000004_create_tags_system.sql (16.28ms)
2025/08/01 15:44:36 OK   000005_create_auth_tables.sql (20.63ms)
2025/08/01 15:44:36 OK   000006_create_insights_table.sql (24.81ms)
2025/08/01 15:44:36 OK   000007_create_performance_indexes.sql (15.9ms)
2025/08/01 15:44:36 OK   000008_create_analytics_views.sql (29.47ms)
2025/08/01 15:44:36 OK   000009_development_data.sql (29.09ms)
2025/08/01 15:44:36 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36374/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36374/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:37 🐳 Stopping container: 38b573ccad62
2025/08/01 15:44:37 ✅ Container stopped: 38b573ccad62
2025/08/01 15:44:37 🐳 Terminating container: 38b573ccad62
2025/08/01 15:44:37 🚫 Container terminated: 38b573ccad62
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_with_very_long_name
2025/08/01 15:44:37 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:37 ✅ Container created: 7e1a3094ca29
2025/08/01 15:44:37 🐳 Starting container: 7e1a3094ca29
2025/08/01 15:44:37 ✅ Container started: 7e1a3094ca29
2025/08/01 15:44:37 ⏳ Waiting for container id 7e1a3094ca29 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0083a4a18 Strategies:[0xc007cbb440 0xc007cff530]}
2025/08/01 15:44:39 🔔 Container is ready: 7e1a3094ca29
DSN: postgres://testuser:testpass@localhost:36375/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36375/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:39 OK   000001_create_users_table.sql (9.87ms)
2025/08/01 15:44:39 OK   000002_create_projects_table.sql (12.84ms)
2025/08/01 15:44:39 OK   000003_create_log_entries_table.sql (19.75ms)
2025/08/01 15:44:39 OK   000004_create_tags_system.sql (16.33ms)
2025/08/01 15:44:39 OK   000005_create_auth_tables.sql (20.24ms)
2025/08/01 15:44:39 OK   000006_create_insights_table.sql (24.27ms)
2025/08/01 15:44:39 OK   000007_create_performance_indexes.sql (15.23ms)
2025/08/01 15:44:39 OK   000008_create_analytics_views.sql (28.02ms)
2025/08/01 15:44:39 OK   000009_development_data.sql (29.45ms)
2025/08/01 15:44:39 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36375/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36375/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:39 🐳 Stopping container: 7e1a3094ca29
2025/08/01 15:44:40 ✅ Container stopped: 7e1a3094ca29
2025/08/01 15:44:40 🐳 Terminating container: 7e1a3094ca29
2025/08/01 15:44:40 🚫 Container terminated: 7e1a3094ca29
=== RUN   TestTagHandler_Comprehensive_UpdateTag/update_with_excessively_long_name
2025/08/01 15:44:40 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:40 ✅ Container created: f29212e2f4df
2025/08/01 15:44:40 🐳 Starting container: f29212e2f4df
2025/08/01 15:44:40 ✅ Container started: f29212e2f4df
2025/08/01 15:44:40 ⏳ Waiting for container id f29212e2f4df image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009b09498 Strategies:[0xc009b10b40 0xc009b23b90]}
2025/08/01 15:44:42 🔔 Container is ready: f29212e2f4df
DSN: postgres://testuser:testpass@localhost:36376/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36376/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:42 OK   000001_create_users_table.sql (10.64ms)
2025/08/01 15:44:42 OK   000002_create_projects_table.sql (12.31ms)
2025/08/01 15:44:42 OK   000003_create_log_entries_table.sql (19.4ms)
2025/08/01 15:44:42 OK   000004_create_tags_system.sql (16.69ms)
2025/08/01 15:44:42 OK   000005_create_auth_tables.sql (20.6ms)
2025/08/01 15:44:42 OK   000006_create_insights_table.sql (24.23ms)
2025/08/01 15:44:42 OK   000007_create_performance_indexes.sql (15.16ms)
2025/08/01 15:44:42 OK   000008_create_analytics_views.sql (28.11ms)
2025/08/01 15:44:42 OK   000009_development_data.sql (28.49ms)
2025/08/01 15:44:42 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36376/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36376/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:42 🐳 Stopping container: f29212e2f4df
2025/08/01 15:44:42 ✅ Container stopped: f29212e2f4df
2025/08/01 15:44:42 🐳 Terminating container: f29212e2f4df
2025/08/01 15:44:43 🚫 Container terminated: f29212e2f4df
--- PASS: TestTagHandler_Comprehensive_UpdateTag (19.75s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_all_fields (2.88s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_only_name (2.80s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_only_color (2.83s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_with_empty_name (2.80s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_with_invalid_color (2.78s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_with_very_long_name (2.79s)
    --- PASS: TestTagHandler_Comprehensive_UpdateTag/update_with_excessively_long_name (2.86s)
=== RUN   TestTagHandler_Comprehensive_DeleteTag
=== RUN   TestTagHandler_Comprehensive_DeleteTag/delete_existing_tag
2025/08/01 15:44:43 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:43 ✅ Container created: 1ee7a20ff593
2025/08/01 15:44:43 🐳 Starting container: 1ee7a20ff593
2025/08/01 15:44:43 ✅ Container started: 1ee7a20ff593
2025/08/01 15:44:43 ⏳ Waiting for container id 1ee7a20ff593 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0020ac560 Strategies:[0xc007b4e7e0 0xc005aa7170]}
2025/08/01 15:44:45 🔔 Container is ready: 1ee7a20ff593
DSN: postgres://testuser:testpass@localhost:36377/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36377/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:45 OK   000001_create_users_table.sql (10.07ms)
2025/08/01 15:44:45 OK   000002_create_projects_table.sql (11.82ms)
2025/08/01 15:44:45 OK   000003_create_log_entries_table.sql (18.49ms)
2025/08/01 15:44:45 OK   000004_create_tags_system.sql (16.11ms)
2025/08/01 15:44:45 OK   000005_create_auth_tables.sql (19.67ms)
2025/08/01 15:44:45 OK   000006_create_insights_table.sql (24.24ms)
2025/08/01 15:44:45 OK   000007_create_performance_indexes.sql (15.26ms)
2025/08/01 15:44:45 OK   000008_create_analytics_views.sql (27.77ms)
2025/08/01 15:44:45 OK   000009_development_data.sql (29.09ms)
2025/08/01 15:44:45 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36377/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36377/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:45 🐳 Stopping container: 1ee7a20ff593
2025/08/01 15:44:45 ✅ Container stopped: 1ee7a20ff593
2025/08/01 15:44:45 🐳 Terminating container: 1ee7a20ff593
2025/08/01 15:44:45 🚫 Container terminated: 1ee7a20ff593
=== RUN   TestTagHandler_Comprehensive_DeleteTag/delete_non-existent_tag
2025/08/01 15:44:45 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:45 ✅ Container created: 01bcec48c1f2
2025/08/01 15:44:45 🐳 Starting container: 01bcec48c1f2
2025/08/01 15:44:46 ✅ Container started: 01bcec48c1f2
2025/08/01 15:44:46 ⏳ Waiting for container id 01bcec48c1f2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009f92ae0 Strategies:[0xc009f706c0 0xc006682fc0]}
2025/08/01 15:44:47 🔔 Container is ready: 01bcec48c1f2
DSN: postgres://testuser:testpass@localhost:36378/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36378/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:47 OK   000001_create_users_table.sql (11.49ms)
2025/08/01 15:44:47 OK   000002_create_projects_table.sql (13.49ms)
2025/08/01 15:44:47 OK   000003_create_log_entries_table.sql (21.26ms)
2025/08/01 15:44:47 OK   000004_create_tags_system.sql (17.72ms)
2025/08/01 15:44:47 OK   000005_create_auth_tables.sql (21.92ms)
2025/08/01 15:44:48 OK   000006_create_insights_table.sql (26.09ms)
2025/08/01 15:44:48 OK   000007_create_performance_indexes.sql (16.29ms)
2025/08/01 15:44:48 OK   000008_create_analytics_views.sql (29.11ms)
2025/08/01 15:44:48 OK   000009_development_data.sql (29.48ms)
2025/08/01 15:44:48 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36378/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36378/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:48 🐳 Stopping container: 01bcec48c1f2
2025/08/01 15:44:48 ✅ Container stopped: 01bcec48c1f2
2025/08/01 15:44:48 🐳 Terminating container: 01bcec48c1f2
2025/08/01 15:44:48 🚫 Container terminated: 01bcec48c1f2
=== RUN   TestTagHandler_Comprehensive_DeleteTag/delete_with_invalid_tag_ID
2025/08/01 15:44:48 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:48 ✅ Container created: b1cfb1069baf
2025/08/01 15:44:48 🐳 Starting container: b1cfb1069baf
2025/08/01 15:44:48 ✅ Container started: b1cfb1069baf
2025/08/01 15:44:48 ⏳ Waiting for container id b1cfb1069baf image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b8a8f8 Strategies:[0xc005a7e0c0 0xc0039cf350]}
2025/08/01 15:44:50 🔔 Container is ready: b1cfb1069baf
DSN: postgres://testuser:testpass@localhost:36379/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36379/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:50 OK   000001_create_users_table.sql (11.02ms)
2025/08/01 15:44:50 OK   000002_create_projects_table.sql (12.92ms)
2025/08/01 15:44:50 OK   000003_create_log_entries_table.sql (19.81ms)
2025/08/01 15:44:50 OK   000004_create_tags_system.sql (16.55ms)
2025/08/01 15:44:50 OK   000005_create_auth_tables.sql (20.45ms)
2025/08/01 15:44:50 OK   000006_create_insights_table.sql (24.86ms)
2025/08/01 15:44:50 OK   000007_create_performance_indexes.sql (15.8ms)
2025/08/01 15:44:50 OK   000008_create_analytics_views.sql (28.57ms)
2025/08/01 15:44:50 OK   000009_development_data.sql (29.96ms)
2025/08/01 15:44:50 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36379/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36379/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:51 🐳 Stopping container: b1cfb1069baf
2025/08/01 15:44:51 ✅ Container stopped: b1cfb1069baf
2025/08/01 15:44:51 🐳 Terminating container: b1cfb1069baf
2025/08/01 15:44:51 🚫 Container terminated: b1cfb1069baf
--- PASS: TestTagHandler_Comprehensive_DeleteTag (8.43s)
    --- PASS: TestTagHandler_Comprehensive_DeleteTag/delete_existing_tag (2.80s)
    --- PASS: TestTagHandler_Comprehensive_DeleteTag/delete_non-existent_tag (2.85s)
    --- PASS: TestTagHandler_Comprehensive_DeleteTag/delete_with_invalid_tag_ID (2.79s)
=== RUN   TestTagHandler_Comprehensive_PopularTags
=== RUN   TestTagHandler_Comprehensive_PopularTags/default_limit
2025/08/01 15:44:51 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:51 ✅ Container created: e8f106097117
2025/08/01 15:44:51 🐳 Starting container: e8f106097117
2025/08/01 15:44:51 ✅ Container started: e8f106097117
2025/08/01 15:44:51 ⏳ Waiting for container id e8f106097117 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007d44370 Strategies:[0xc009947b60 0xc006aa1b90]}
2025/08/01 15:44:53 🔔 Container is ready: e8f106097117
DSN: postgres://testuser:testpass@localhost:36380/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36380/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:53 OK   000001_create_users_table.sql (10.85ms)
2025/08/01 15:44:53 OK   000002_create_projects_table.sql (12.66ms)
2025/08/01 15:44:53 OK   000003_create_log_entries_table.sql (20.41ms)
2025/08/01 15:44:53 OK   000004_create_tags_system.sql (16.87ms)
2025/08/01 15:44:53 OK   000005_create_auth_tables.sql (21.24ms)
2025/08/01 15:44:53 OK   000006_create_insights_table.sql (25.55ms)
2025/08/01 15:44:53 OK   000007_create_performance_indexes.sql (16.23ms)
2025/08/01 15:44:53 OK   000008_create_analytics_views.sql (29.32ms)
2025/08/01 15:44:53 OK   000009_development_data.sql (29.58ms)
2025/08/01 15:44:53 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36380/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36380/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:53 🐳 Stopping container: e8f106097117
2025/08/01 15:44:54 ✅ Container stopped: e8f106097117
2025/08/01 15:44:54 🐳 Terminating container: e8f106097117
2025/08/01 15:44:54 🚫 Container terminated: e8f106097117
=== RUN   TestTagHandler_Comprehensive_PopularTags/custom_limit
2025/08/01 15:44:54 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:54 ✅ Container created: 8af81a975a03
2025/08/01 15:44:54 🐳 Starting container: 8af81a975a03
2025/08/01 15:44:54 ✅ Container started: 8af81a975a03
2025/08/01 15:44:54 ⏳ Waiting for container id 8af81a975a03 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006421718 Strategies:[0xc006c07440 0xc006849aa0]}
2025/08/01 15:44:56 🔔 Container is ready: 8af81a975a03
DSN: postgres://testuser:testpass@localhost:36381/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36381/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:56 OK   000001_create_users_table.sql (10.18ms)
2025/08/01 15:44:56 OK   000002_create_projects_table.sql (11.78ms)
2025/08/01 15:44:56 OK   000003_create_log_entries_table.sql (18.82ms)
2025/08/01 15:44:56 OK   000004_create_tags_system.sql (16.36ms)
2025/08/01 15:44:56 OK   000005_create_auth_tables.sql (20.42ms)
2025/08/01 15:44:56 OK   000006_create_insights_table.sql (24.37ms)
2025/08/01 15:44:56 OK   000007_create_performance_indexes.sql (15.04ms)
2025/08/01 15:44:56 OK   000008_create_analytics_views.sql (28.16ms)
2025/08/01 15:44:56 OK   000009_development_data.sql (28.6ms)
2025/08/01 15:44:56 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36381/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36381/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:56 🐳 Stopping container: 8af81a975a03
2025/08/01 15:44:57 ✅ Container stopped: 8af81a975a03
2025/08/01 15:44:57 🐳 Terminating container: 8af81a975a03
2025/08/01 15:44:57 🚫 Container terminated: 8af81a975a03
=== RUN   TestTagHandler_Comprehensive_PopularTags/maximum_limit
2025/08/01 15:44:57 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:57 ✅ Container created: 62a16787afd5
2025/08/01 15:44:57 🐳 Starting container: 62a16787afd5
2025/08/01 15:44:57 ✅ Container started: 62a16787afd5
2025/08/01 15:44:57 ⏳ Waiting for container id 62a16787afd5 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0046295f8 Strategies:[0xc006921f20 0xc002760210]}
2025/08/01 15:44:59 🔔 Container is ready: 62a16787afd5
DSN: postgres://testuser:testpass@localhost:36382/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36382/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:44:59 OK   000001_create_users_table.sql (10.27ms)
2025/08/01 15:44:59 OK   000002_create_projects_table.sql (12.37ms)
2025/08/01 15:44:59 OK   000003_create_log_entries_table.sql (18.93ms)
2025/08/01 15:44:59 OK   000004_create_tags_system.sql (16.04ms)
2025/08/01 15:44:59 OK   000005_create_auth_tables.sql (20.19ms)
2025/08/01 15:44:59 OK   000006_create_insights_table.sql (24.62ms)
2025/08/01 15:44:59 OK   000007_create_performance_indexes.sql (15.42ms)
2025/08/01 15:44:59 OK   000008_create_analytics_views.sql (27.58ms)
2025/08/01 15:44:59 OK   000009_development_data.sql (26.84ms)
2025/08/01 15:44:59 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36382/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36382/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:44:59 🐳 Stopping container: 62a16787afd5
2025/08/01 15:44:59 ✅ Container stopped: 62a16787afd5
2025/08/01 15:44:59 🐳 Terminating container: 62a16787afd5
2025/08/01 15:44:59 🚫 Container terminated: 62a16787afd5
=== RUN   TestTagHandler_Comprehensive_PopularTags/excessive_limit
2025/08/01 15:44:59 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:44:59 ✅ Container created: 90fff4338a6e
2025/08/01 15:44:59 🐳 Starting container: 90fff4338a6e
2025/08/01 15:45:00 ✅ Container started: 90fff4338a6e
2025/08/01 15:45:00 ⏳ Waiting for container id 90fff4338a6e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b15ad8 Strategies:[0xc00185ad20 0xc0064407b0]}
2025/08/01 15:45:01 🔔 Container is ready: 90fff4338a6e
DSN: postgres://testuser:testpass@localhost:36383/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36383/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:01 OK   000001_create_users_table.sql (10.49ms)
2025/08/01 15:45:01 OK   000002_create_projects_table.sql (11.78ms)
2025/08/01 15:45:01 OK   000003_create_log_entries_table.sql (18.92ms)
2025/08/01 15:45:01 OK   000004_create_tags_system.sql (16.05ms)
2025/08/01 15:45:01 OK   000005_create_auth_tables.sql (20.2ms)
2025/08/01 15:45:01 OK   000006_create_insights_table.sql (24.91ms)
2025/08/01 15:45:01 OK   000007_create_performance_indexes.sql (15.48ms)
2025/08/01 15:45:01 OK   000008_create_analytics_views.sql (27.61ms)
2025/08/01 15:45:02 OK   000009_development_data.sql (28.62ms)
2025/08/01 15:45:02 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36383/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36383/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:02 🐳 Stopping container: 90fff4338a6e
2025/08/01 15:45:02 ✅ Container stopped: 90fff4338a6e
2025/08/01 15:45:02 🐳 Terminating container: 90fff4338a6e
2025/08/01 15:45:02 🚫 Container terminated: 90fff4338a6e
=== RUN   TestTagHandler_Comprehensive_PopularTags/invalid_limit
2025/08/01 15:45:02 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:02 ✅ Container created: d793ef52be7f
2025/08/01 15:45:02 🐳 Starting container: d793ef52be7f
2025/08/01 15:45:02 ✅ Container started: d793ef52be7f
2025/08/01 15:45:02 ⏳ Waiting for container id d793ef52be7f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0069e9af8 Strategies:[0xc006b16d80 0xc007170150]}
2025/08/01 15:45:04 🔔 Container is ready: d793ef52be7f
DSN: postgres://testuser:testpass@localhost:36384/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36384/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:04 OK   000001_create_users_table.sql (11.22ms)
2025/08/01 15:45:04 OK   000002_create_projects_table.sql (13.31ms)
2025/08/01 15:45:04 OK   000003_create_log_entries_table.sql (20.48ms)
2025/08/01 15:45:04 OK   000004_create_tags_system.sql (17.51ms)
2025/08/01 15:45:04 OK   000005_create_auth_tables.sql (21.14ms)
2025/08/01 15:45:04 OK   000006_create_insights_table.sql (24.77ms)
2025/08/01 15:45:04 OK   000007_create_performance_indexes.sql (15.4ms)
2025/08/01 15:45:04 OK   000008_create_analytics_views.sql (28.27ms)
2025/08/01 15:45:04 OK   000009_development_data.sql (29.01ms)
2025/08/01 15:45:04 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36384/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36384/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:05 🐳 Stopping container: d793ef52be7f
2025/08/01 15:45:05 ✅ Container stopped: d793ef52be7f
2025/08/01 15:45:05 🐳 Terminating container: d793ef52be7f
2025/08/01 15:45:05 🚫 Container terminated: d793ef52be7f
=== RUN   TestTagHandler_Comprehensive_PopularTags/negative_limit
2025/08/01 15:45:05 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:05 ✅ Container created: b753ab7cd410
2025/08/01 15:45:05 🐳 Starting container: b753ab7cd410
2025/08/01 15:45:05 ✅ Container started: b753ab7cd410
2025/08/01 15:45:05 ⏳ Waiting for container id b753ab7cd410 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007244628 Strategies:[0xc006c3a1e0 0xc0083ee930]}
2025/08/01 15:45:07 🔔 Container is ready: b753ab7cd410
DSN: postgres://testuser:testpass@localhost:36385/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36385/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:07 OK   000001_create_users_table.sql (11.16ms)
2025/08/01 15:45:07 OK   000002_create_projects_table.sql (13.5ms)
2025/08/01 15:45:07 OK   000003_create_log_entries_table.sql (20.3ms)
2025/08/01 15:45:07 OK   000004_create_tags_system.sql (18.09ms)
2025/08/01 15:45:07 OK   000005_create_auth_tables.sql (21.71ms)
2025/08/01 15:45:07 OK   000006_create_insights_table.sql (25.78ms)
2025/08/01 15:45:07 OK   000007_create_performance_indexes.sql (15.36ms)
2025/08/01 15:45:07 OK   000008_create_analytics_views.sql (28.41ms)
2025/08/01 15:45:07 OK   000009_development_data.sql (29.17ms)
2025/08/01 15:45:07 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36385/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36385/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:07 🐳 Stopping container: b753ab7cd410
2025/08/01 15:45:08 ✅ Container stopped: b753ab7cd410
2025/08/01 15:45:08 🐳 Terminating container: b753ab7cd410
2025/08/01 15:45:08 🚫 Container terminated: b753ab7cd410
=== RUN   TestTagHandler_Comprehensive_PopularTags/zero_limit
2025/08/01 15:45:08 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:08 ✅ Container created: f482ec2c2c22
2025/08/01 15:45:08 🐳 Starting container: f482ec2c2c22
2025/08/01 15:45:08 ✅ Container started: f482ec2c2c22
2025/08/01 15:45:08 ⏳ Waiting for container id f482ec2c2c22 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00668ddf8 Strategies:[0xc006d36540 0xc00a080510]}
2025/08/01 15:45:10 🔔 Container is ready: f482ec2c2c22
DSN: postgres://testuser:testpass@localhost:36386/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36386/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:10 OK   000001_create_users_table.sql (10.8ms)
2025/08/01 15:45:10 OK   000002_create_projects_table.sql (14.14ms)
2025/08/01 15:45:10 OK   000003_create_log_entries_table.sql (19.06ms)
2025/08/01 15:45:10 OK   000004_create_tags_system.sql (15.98ms)
2025/08/01 15:45:10 OK   000005_create_auth_tables.sql (20.16ms)
2025/08/01 15:45:10 OK   000006_create_insights_table.sql (25.08ms)
2025/08/01 15:45:10 OK   000007_create_performance_indexes.sql (15.53ms)
2025/08/01 15:45:10 OK   000008_create_analytics_views.sql (28.03ms)
2025/08/01 15:45:10 OK   000009_development_data.sql (28.79ms)
2025/08/01 15:45:10 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36386/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36386/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:10 🐳 Stopping container: f482ec2c2c22
2025/08/01 15:45:10 ✅ Container stopped: f482ec2c2c22
2025/08/01 15:45:10 🐳 Terminating container: f482ec2c2c22
2025/08/01 15:45:10 🚫 Container terminated: f482ec2c2c22
--- PASS: TestTagHandler_Comprehensive_PopularTags (19.51s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/default_limit (2.81s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/custom_limit (2.77s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/maximum_limit (2.74s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/excessive_limit (2.83s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/invalid_limit (2.82s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/negative_limit (2.79s)
    --- PASS: TestTagHandler_Comprehensive_PopularTags/zero_limit (2.74s)
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/POST_/v1/tags_without_authentication
2025/08/01 15:45:10 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:11 ✅ Container created: ad063d35c1e9
2025/08/01 15:45:11 🐳 Starting container: ad063d35c1e9
2025/08/01 15:45:11 ✅ Container started: ad063d35c1e9
2025/08/01 15:45:11 ⏳ Waiting for container id ad063d35c1e9 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007e55e60 Strategies:[0xc007cbac00 0xc0068c6570]}
2025/08/01 15:45:12 🔔 Container is ready: ad063d35c1e9
DSN: postgres://testuser:testpass@localhost:36387/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36387/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:13 OK   000001_create_users_table.sql (11.12ms)
2025/08/01 15:45:13 OK   000002_create_projects_table.sql (12.8ms)
2025/08/01 15:45:13 OK   000003_create_log_entries_table.sql (19.74ms)
2025/08/01 15:45:13 OK   000004_create_tags_system.sql (16.99ms)
2025/08/01 15:45:13 OK   000005_create_auth_tables.sql (20.84ms)
2025/08/01 15:45:13 OK   000006_create_insights_table.sql (24.37ms)
2025/08/01 15:45:13 OK   000007_create_performance_indexes.sql (15.56ms)
2025/08/01 15:45:13 OK   000008_create_analytics_views.sql (27.3ms)
2025/08/01 15:45:13 OK   000009_development_data.sql (28.05ms)
2025/08/01 15:45:13 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36387/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36387/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:13 🐳 Stopping container: ad063d35c1e9
2025/08/01 15:45:13 ✅ Container stopped: ad063d35c1e9
2025/08/01 15:45:13 🐳 Terminating container: ad063d35c1e9
2025/08/01 15:45:13 🚫 Container terminated: ad063d35c1e9
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/POST_/v1/tags_with_invalid_token
2025/08/01 15:45:13 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:13 ✅ Container created: f026c39beee0
2025/08/01 15:45:13 🐳 Starting container: f026c39beee0
2025/08/01 15:45:13 ✅ Container started: f026c39beee0
2025/08/01 15:45:13 ⏳ Waiting for container id f026c39beee0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009fc8bc0 Strategies:[0xc009fbeb40 0xc00721e5d0]}
2025/08/01 15:45:15 🔔 Container is ready: f026c39beee0
DSN: postgres://testuser:testpass@localhost:36388/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36388/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:15 OK   000001_create_users_table.sql (10.84ms)
2025/08/01 15:45:15 OK   000002_create_projects_table.sql (12.64ms)
2025/08/01 15:45:15 OK   000003_create_log_entries_table.sql (19.72ms)
2025/08/01 15:45:15 OK   000004_create_tags_system.sql (16.99ms)
2025/08/01 15:45:15 OK   000005_create_auth_tables.sql (20.34ms)
2025/08/01 15:45:15 OK   000006_create_insights_table.sql (24.35ms)
2025/08/01 15:45:15 OK   000007_create_performance_indexes.sql (15.54ms)
2025/08/01 15:45:15 OK   000008_create_analytics_views.sql (28.18ms)
2025/08/01 15:45:15 OK   000009_development_data.sql (28.95ms)
2025/08/01 15:45:15 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36388/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36388/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:15 🐳 Stopping container: f026c39beee0
2025/08/01 15:45:15 ✅ Container stopped: f026c39beee0
2025/08/01 15:45:15 🐳 Terminating container: f026c39beee0
2025/08/01 15:45:16 🚫 Container terminated: f026c39beee0
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags_without_authentication
2025/08/01 15:45:16 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:16 ✅ Container created: 85156e6f34a7
2025/08/01 15:45:16 🐳 Starting container: 85156e6f34a7
2025/08/01 15:45:16 ✅ Container started: 85156e6f34a7
2025/08/01 15:45:16 ⏳ Waiting for container id 85156e6f34a7 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009de04d8 Strategies:[0xc006a5fc80 0xc0054dd380]}
2025/08/01 15:45:17 🔔 Container is ready: 85156e6f34a7
DSN: postgres://testuser:testpass@localhost:36389/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36389/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:18 OK   000001_create_users_table.sql (10.57ms)
2025/08/01 15:45:18 OK   000002_create_projects_table.sql (11.71ms)
2025/08/01 15:45:18 OK   000003_create_log_entries_table.sql (18.81ms)
2025/08/01 15:45:18 OK   000004_create_tags_system.sql (15.77ms)
2025/08/01 15:45:18 OK   000005_create_auth_tables.sql (19.8ms)
2025/08/01 15:45:18 OK   000006_create_insights_table.sql (23.55ms)
2025/08/01 15:45:18 OK   000007_create_performance_indexes.sql (15.08ms)
2025/08/01 15:45:18 OK   000008_create_analytics_views.sql (28.01ms)
2025/08/01 15:45:18 OK   000009_development_data.sql (28.16ms)
2025/08/01 15:45:18 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36389/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36389/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:18 🐳 Stopping container: 85156e6f34a7
2025/08/01 15:45:18 ✅ Container stopped: 85156e6f34a7
2025/08/01 15:45:18 🐳 Terminating container: 85156e6f34a7
2025/08/01 15:45:18 🚫 Container terminated: 85156e6f34a7
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags_with_invalid_token
2025/08/01 15:45:18 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:18 ✅ Container created: fd5fd26fad4f
2025/08/01 15:45:18 🐳 Starting container: fd5fd26fad4f
2025/08/01 15:45:18 ✅ Container started: fd5fd26fad4f
2025/08/01 15:45:18 ⏳ Waiting for container id fd5fd26fad4f image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00ab8afc8 Strategies:[0xc001eb39e0 0xc003af2150]}
2025/08/01 15:45:20 🔔 Container is ready: fd5fd26fad4f
DSN: postgres://testuser:testpass@localhost:36390/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36390/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:20 OK   000001_create_users_table.sql (11.47ms)
2025/08/01 15:45:20 OK   000002_create_projects_table.sql (12.94ms)
2025/08/01 15:45:20 OK   000003_create_log_entries_table.sql (20.05ms)
2025/08/01 15:45:20 OK   000004_create_tags_system.sql (16.63ms)
2025/08/01 15:45:20 OK   000005_create_auth_tables.sql (20.95ms)
2025/08/01 15:45:20 OK   000006_create_insights_table.sql (24.93ms)
2025/08/01 15:45:20 OK   000007_create_performance_indexes.sql (15.67ms)
2025/08/01 15:45:20 OK   000008_create_analytics_views.sql (28.4ms)
2025/08/01 15:45:20 OK   000009_development_data.sql (25.8ms)
2025/08/01 15:45:20 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36390/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36390/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:20 🐳 Stopping container: fd5fd26fad4f
2025/08/01 15:45:20 ✅ Container stopped: fd5fd26fad4f
2025/08/01 15:45:20 🐳 Terminating container: fd5fd26fad4f
2025/08/01 15:45:21 🚫 Container terminated: fd5fd26fad4f
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/popular_without_authentication
2025/08/01 15:45:21 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:21 ✅ Container created: 9e275330e486
2025/08/01 15:45:21 🐳 Starting container: 9e275330e486
2025/08/01 15:45:21 ✅ Container started: 9e275330e486
2025/08/01 15:45:21 ⏳ Waiting for container id 9e275330e486 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00067ae48 Strategies:[0xc0001e32c0 0xc006449ef0]}
2025/08/01 15:45:23 🔔 Container is ready: 9e275330e486
DSN: postgres://testuser:testpass@localhost:36391/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36391/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:23 OK   000001_create_users_table.sql (10.56ms)
2025/08/01 15:45:23 OK   000002_create_projects_table.sql (12.63ms)
2025/08/01 15:45:23 OK   000003_create_log_entries_table.sql (20.76ms)
2025/08/01 15:45:23 OK   000004_create_tags_system.sql (16.47ms)
2025/08/01 15:45:23 OK   000005_create_auth_tables.sql (20.07ms)
2025/08/01 15:45:23 OK   000006_create_insights_table.sql (24.06ms)
2025/08/01 15:45:23 OK   000007_create_performance_indexes.sql (15.34ms)
2025/08/01 15:45:23 OK   000008_create_analytics_views.sql (27.62ms)
2025/08/01 15:45:23 OK   000009_development_data.sql (28.57ms)
2025/08/01 15:45:23 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36391/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36391/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:23 🐳 Stopping container: 9e275330e486
2025/08/01 15:45:23 ✅ Container stopped: 9e275330e486
2025/08/01 15:45:23 🐳 Terminating container: 9e275330e486
2025/08/01 15:45:23 🚫 Container terminated: 9e275330e486
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/popular_with_invalid_token
2025/08/01 15:45:23 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:23 ✅ Container created: 87633e5005f2
2025/08/01 15:45:23 🐳 Starting container: 87633e5005f2
2025/08/01 15:45:23 ✅ Container started: 87633e5005f2
2025/08/01 15:45:23 ⏳ Waiting for container id 87633e5005f2 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001f09698 Strategies:[0xc006ce2de0 0xc0027603c0]}
2025/08/01 15:45:25 🔔 Container is ready: 87633e5005f2
DSN: postgres://testuser:testpass@localhost:36393/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36393/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:25 OK   000001_create_users_table.sql (10.3ms)
2025/08/01 15:45:25 OK   000002_create_projects_table.sql (12.1ms)
2025/08/01 15:45:25 OK   000003_create_log_entries_table.sql (19.06ms)
2025/08/01 15:45:25 OK   000004_create_tags_system.sql (16.74ms)
2025/08/01 15:45:25 OK   000005_create_auth_tables.sql (20.42ms)
2025/08/01 15:45:25 OK   000006_create_insights_table.sql (24.44ms)
2025/08/01 15:45:25 OK   000007_create_performance_indexes.sql (15.39ms)
2025/08/01 15:45:25 OK   000008_create_analytics_views.sql (27.55ms)
2025/08/01 15:45:25 OK   000009_development_data.sql (27.98ms)
2025/08/01 15:45:25 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36393/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36393/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:25 🐳 Stopping container: 87633e5005f2
2025/08/01 15:45:26 ✅ Container stopped: 87633e5005f2
2025/08/01 15:45:26 🐳 Terminating container: 87633e5005f2
2025/08/01 15:45:26 🚫 Container terminated: 87633e5005f2
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/recent_without_authentication
2025/08/01 15:45:26 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:26 ✅ Container created: d66e696ffbf8
2025/08/01 15:45:26 🐳 Starting container: d66e696ffbf8
2025/08/01 15:45:26 ✅ Container started: d66e696ffbf8
2025/08/01 15:45:26 ⏳ Waiting for container id d66e696ffbf8 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b325d8 Strategies:[0xc0045c73e0 0xc005b333e0]}
2025/08/01 15:45:28 🔔 Container is ready: d66e696ffbf8
DSN: postgres://testuser:testpass@localhost:36395/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36395/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:28 OK   000001_create_users_table.sql (10.65ms)
2025/08/01 15:45:28 OK   000002_create_projects_table.sql (12.18ms)
2025/08/01 15:45:28 OK   000003_create_log_entries_table.sql (19.3ms)
2025/08/01 15:45:28 OK   000004_create_tags_system.sql (16.34ms)
2025/08/01 15:45:28 OK   000005_create_auth_tables.sql (20.99ms)
2025/08/01 15:45:28 OK   000006_create_insights_table.sql (24.72ms)
2025/08/01 15:45:28 OK   000007_create_performance_indexes.sql (15.49ms)
2025/08/01 15:45:28 OK   000008_create_analytics_views.sql (28.83ms)
2025/08/01 15:45:28 OK   000009_development_data.sql (29.92ms)
2025/08/01 15:45:28 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36395/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36395/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:28 🐳 Stopping container: d66e696ffbf8
2025/08/01 15:45:28 ✅ Container stopped: d66e696ffbf8
2025/08/01 15:45:28 🐳 Terminating container: d66e696ffbf8
2025/08/01 15:45:28 🚫 Container terminated: d66e696ffbf8
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/recent_with_invalid_token
2025/08/01 15:45:28 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:28 ✅ Container created: 7fc7508fe58a
2025/08/01 15:45:28 🐳 Starting container: 7fc7508fe58a
2025/08/01 15:45:28 ✅ Container started: 7fc7508fe58a
2025/08/01 15:45:28 ⏳ Waiting for container id 7fc7508fe58a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002bf6828 Strategies:[0xc0048fc2a0 0xc0083eedb0]}
2025/08/01 15:45:30 🔔 Container is ready: 7fc7508fe58a
DSN: postgres://testuser:testpass@localhost:36396/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36396/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:30 OK   000001_create_users_table.sql (10.4ms)
2025/08/01 15:45:30 OK   000002_create_projects_table.sql (12.13ms)
2025/08/01 15:45:30 OK   000003_create_log_entries_table.sql (19.37ms)
2025/08/01 15:45:30 OK   000004_create_tags_system.sql (16.22ms)
2025/08/01 15:45:30 OK   000005_create_auth_tables.sql (20.18ms)
2025/08/01 15:45:30 OK   000006_create_insights_table.sql (24.22ms)
2025/08/01 15:45:30 OK   000007_create_performance_indexes.sql (15.67ms)
2025/08/01 15:45:30 OK   000008_create_analytics_views.sql (28ms)
2025/08/01 15:45:31 OK   000009_development_data.sql (27.22ms)
2025/08/01 15:45:31 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36396/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36396/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:31 🐳 Stopping container: 7fc7508fe58a
2025/08/01 15:45:31 ✅ Container stopped: 7fc7508fe58a
2025/08/01 15:45:31 🐳 Terminating container: 7fc7508fe58a
2025/08/01 15:45:31 🚫 Container terminated: 7fc7508fe58a
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/search?q=test_without_authentication
2025/08/01 15:45:31 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:31 ✅ Container created: 26f1d9ba6049
2025/08/01 15:45:31 🐳 Starting container: 26f1d9ba6049
2025/08/01 15:45:31 ✅ Container started: 26f1d9ba6049
2025/08/01 15:45:31 ⏳ Waiting for container id 26f1d9ba6049 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004247698 Strategies:[0xc009946480 0xc00a1066f0]}
2025/08/01 15:45:33 🔔 Container is ready: 26f1d9ba6049
DSN: postgres://testuser:testpass@localhost:36397/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36397/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:33 OK   000001_create_users_table.sql (12.65ms)
2025/08/01 15:45:33 OK   000002_create_projects_table.sql (14.33ms)
2025/08/01 15:45:33 OK   000003_create_log_entries_table.sql (20.97ms)
2025/08/01 15:45:33 OK   000004_create_tags_system.sql (16.49ms)
2025/08/01 15:45:33 OK   000005_create_auth_tables.sql (21.1ms)
2025/08/01 15:45:33 OK   000006_create_insights_table.sql (24.88ms)
2025/08/01 15:45:33 OK   000007_create_performance_indexes.sql (15.38ms)
2025/08/01 15:45:33 OK   000008_create_analytics_views.sql (27.91ms)
2025/08/01 15:45:33 OK   000009_development_data.sql (28.89ms)
2025/08/01 15:45:33 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36397/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36397/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:33 🐳 Stopping container: 26f1d9ba6049
2025/08/01 15:45:33 ✅ Container stopped: 26f1d9ba6049
2025/08/01 15:45:33 🐳 Terminating container: 26f1d9ba6049
2025/08/01 15:45:33 🚫 Container terminated: 26f1d9ba6049
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/search?q=test_with_invalid_token
2025/08/01 15:45:33 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:33 ✅ Container created: ab381952d989
2025/08/01 15:45:33 🐳 Starting container: ab381952d989
2025/08/01 15:45:34 ✅ Container started: ab381952d989
2025/08/01 15:45:34 ⏳ Waiting for container id ab381952d989 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00483d0e8 Strategies:[0xc00434fa40 0xc006ca7320]}
2025/08/01 15:45:35 🔔 Container is ready: ab381952d989
DSN: postgres://testuser:testpass@localhost:36399/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36399/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:35 OK   000001_create_users_table.sql (10.31ms)
2025/08/01 15:45:36 OK   000002_create_projects_table.sql (12.19ms)
2025/08/01 15:45:36 OK   000003_create_log_entries_table.sql (19.22ms)
2025/08/01 15:45:36 OK   000004_create_tags_system.sql (16.42ms)
2025/08/01 15:45:36 OK   000005_create_auth_tables.sql (21.09ms)
2025/08/01 15:45:36 OK   000006_create_insights_table.sql (24.82ms)
2025/08/01 15:45:36 OK   000007_create_performance_indexes.sql (15.38ms)
2025/08/01 15:45:36 OK   000008_create_analytics_views.sql (27.86ms)
2025/08/01 15:45:36 OK   000009_development_data.sql (29.03ms)
2025/08/01 15:45:36 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36399/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36399/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:36 🐳 Stopping container: ab381952d989
2025/08/01 15:45:36 ✅ Container stopped: ab381952d989
2025/08/01 15:45:36 🐳 Terminating container: ab381952d989
2025/08/01 15:45:36 🚫 Container terminated: ab381952d989
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/usage_without_authentication
2025/08/01 15:45:36 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:36 ✅ Container created: 0efc82c6be74
2025/08/01 15:45:36 🐳 Starting container: 0efc82c6be74
2025/08/01 15:45:36 ✅ Container started: 0efc82c6be74
2025/08/01 15:45:36 ⏳ Waiting for container id 0efc82c6be74 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007e9cef8 Strategies:[0xc006a5fc80 0xc006a86d50]}
2025/08/01 15:45:38 🔔 Container is ready: 0efc82c6be74
DSN: postgres://testuser:testpass@localhost:36400/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36400/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:38 OK   000001_create_users_table.sql (10.2ms)
2025/08/01 15:45:38 OK   000002_create_projects_table.sql (12.32ms)
2025/08/01 15:45:38 OK   000003_create_log_entries_table.sql (18.84ms)
2025/08/01 15:45:38 OK   000004_create_tags_system.sql (16.5ms)
2025/08/01 15:45:38 OK   000005_create_auth_tables.sql (20ms)
2025/08/01 15:45:38 OK   000006_create_insights_table.sql (24.07ms)
2025/08/01 15:45:38 OK   000007_create_performance_indexes.sql (16.5ms)
2025/08/01 15:45:38 OK   000008_create_analytics_views.sql (30.53ms)
2025/08/01 15:45:38 OK   000009_development_data.sql (30.09ms)
2025/08/01 15:45:38 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36400/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36400/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:38 🐳 Stopping container: 0efc82c6be74
2025/08/01 15:45:39 ✅ Container stopped: 0efc82c6be74
2025/08/01 15:45:39 🐳 Terminating container: 0efc82c6be74
2025/08/01 15:45:39 🚫 Container terminated: 0efc82c6be74
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/usage_with_invalid_token
2025/08/01 15:45:39 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:39 ✅ Container created: 06fbefca0f7e
2025/08/01 15:45:39 🐳 Starting container: 06fbefca0f7e
2025/08/01 15:45:39 ✅ Container started: 06fbefca0f7e
2025/08/01 15:45:39 ⏳ Waiting for container id 06fbefca0f7e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003abeea8 Strategies:[0xc009b10d80 0xc00aadbda0]}
2025/08/01 15:45:41 🔔 Container is ready: 06fbefca0f7e
DSN: postgres://testuser:testpass@localhost:36401/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36401/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:41 OK   000001_create_users_table.sql (11.21ms)
2025/08/01 15:45:41 OK   000002_create_projects_table.sql (13.38ms)
2025/08/01 15:45:41 OK   000003_create_log_entries_table.sql (20.84ms)
2025/08/01 15:45:41 OK   000004_create_tags_system.sql (16.8ms)
2025/08/01 15:45:41 OK   000005_create_auth_tables.sql (20.83ms)
2025/08/01 15:45:41 OK   000006_create_insights_table.sql (25.39ms)
2025/08/01 15:45:41 OK   000007_create_performance_indexes.sql (15.89ms)
2025/08/01 15:45:41 OK   000008_create_analytics_views.sql (28.98ms)
2025/08/01 15:45:41 OK   000009_development_data.sql (29.21ms)
2025/08/01 15:45:41 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36401/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36401/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:41 🐳 Stopping container: 06fbefca0f7e
2025/08/01 15:45:41 ✅ Container stopped: 06fbefca0f7e
2025/08/01 15:45:41 🐳 Terminating container: 06fbefca0f7e
2025/08/01 15:45:41 🚫 Container terminated: 06fbefca0f7e
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/PUT_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication
2025/08/01 15:45:41 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:41 ✅ Container created: a2098d5ba3f0
2025/08/01 15:45:41 🐳 Starting container: a2098d5ba3f0
2025/08/01 15:45:41 ✅ Container started: a2098d5ba3f0
2025/08/01 15:45:41 ⏳ Waiting for container id a2098d5ba3f0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006a04e88 Strategies:[0xc0084029c0 0xc001172870]}
2025/08/01 15:45:43 🔔 Container is ready: a2098d5ba3f0
DSN: postgres://testuser:testpass@localhost:36402/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36402/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:43 OK   000001_create_users_table.sql (9.74ms)
2025/08/01 15:45:43 OK   000002_create_projects_table.sql (12.02ms)
2025/08/01 15:45:43 OK   000003_create_log_entries_table.sql (18.56ms)
2025/08/01 15:45:43 OK   000004_create_tags_system.sql (15.99ms)
2025/08/01 15:45:43 OK   000005_create_auth_tables.sql (20.39ms)
2025/08/01 15:45:43 OK   000006_create_insights_table.sql (23.33ms)
2025/08/01 15:45:43 OK   000007_create_performance_indexes.sql (14.73ms)
2025/08/01 15:45:43 OK   000008_create_analytics_views.sql (27.3ms)
2025/08/01 15:45:43 OK   000009_development_data.sql (28.9ms)
2025/08/01 15:45:43 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36402/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36402/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:43 🐳 Stopping container: a2098d5ba3f0
2025/08/01 15:45:44 ✅ Container stopped: a2098d5ba3f0
2025/08/01 15:45:44 🐳 Terminating container: a2098d5ba3f0
2025/08/01 15:45:44 🚫 Container terminated: a2098d5ba3f0
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/PUT_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token
2025/08/01 15:45:44 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:44 ✅ Container created: 58037d0c7c55
2025/08/01 15:45:44 🐳 Starting container: 58037d0c7c55
2025/08/01 15:45:44 ✅ Container started: 58037d0c7c55
2025/08/01 15:45:44 ⏳ Waiting for container id 58037d0c7c55 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0046287b0 Strategies:[0xc0045a92c0 0xc002d298c0]}
2025/08/01 15:45:46 🔔 Container is ready: 58037d0c7c55
DSN: postgres://testuser:testpass@localhost:36403/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36403/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:46 OK   000001_create_users_table.sql (10.56ms)
2025/08/01 15:45:46 OK   000002_create_projects_table.sql (13.64ms)
2025/08/01 15:45:46 OK   000003_create_log_entries_table.sql (20.21ms)
2025/08/01 15:45:46 OK   000004_create_tags_system.sql (16.6ms)
2025/08/01 15:45:46 OK   000005_create_auth_tables.sql (20.46ms)
2025/08/01 15:45:46 OK   000006_create_insights_table.sql (24.01ms)
2025/08/01 15:45:46 OK   000007_create_performance_indexes.sql (15.36ms)
2025/08/01 15:45:46 OK   000008_create_analytics_views.sql (29.02ms)
2025/08/01 15:45:46 OK   000009_development_data.sql (30.51ms)
2025/08/01 15:45:46 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36403/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36403/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:46 🐳 Stopping container: 58037d0c7c55
2025/08/01 15:45:46 ✅ Container stopped: 58037d0c7c55
2025/08/01 15:45:46 🐳 Terminating container: 58037d0c7c55
2025/08/01 15:45:46 🚫 Container terminated: 58037d0c7c55
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/DELETE_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication
2025/08/01 15:45:46 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:46 ✅ Container created: 2e337120cabe
2025/08/01 15:45:46 🐳 Starting container: 2e337120cabe
2025/08/01 15:45:47 ✅ Container started: 2e337120cabe
2025/08/01 15:45:47 ⏳ Waiting for container id 2e337120cabe image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00678d918 Strategies:[0xc00480e3c0 0xc009b22e10]}
2025/08/01 15:45:48 🔔 Container is ready: 2e337120cabe
DSN: postgres://testuser:testpass@localhost:36404/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36404/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:48 OK   000001_create_users_table.sql (10.3ms)
2025/08/01 15:45:48 OK   000002_create_projects_table.sql (12.33ms)
2025/08/01 15:45:48 OK   000003_create_log_entries_table.sql (19.04ms)
2025/08/01 15:45:48 OK   000004_create_tags_system.sql (16.25ms)
2025/08/01 15:45:48 OK   000005_create_auth_tables.sql (18.38ms)
2025/08/01 15:45:49 OK   000006_create_insights_table.sql (22.75ms)
2025/08/01 15:45:49 OK   000007_create_performance_indexes.sql (14.97ms)
2025/08/01 15:45:49 OK   000008_create_analytics_views.sql (27.51ms)
2025/08/01 15:45:49 OK   000009_development_data.sql (28.11ms)
2025/08/01 15:45:49 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36404/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36404/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:49 🐳 Stopping container: 2e337120cabe
2025/08/01 15:45:49 ✅ Container stopped: 2e337120cabe
2025/08/01 15:45:49 🐳 Terminating container: 2e337120cabe
2025/08/01 15:45:49 🚫 Container terminated: 2e337120cabe
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/DELETE_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token
2025/08/01 15:45:49 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:49 ✅ Container created: 86f8883f2ea1
2025/08/01 15:45:49 🐳 Starting container: 86f8883f2ea1
2025/08/01 15:45:49 ✅ Container started: 86f8883f2ea1
2025/08/01 15:45:49 ⏳ Waiting for container id 86f8883f2ea1 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006c403c8 Strategies:[0xc0046bce40 0xc009821dd0]}
2025/08/01 15:45:51 🔔 Container is ready: 86f8883f2ea1
DSN: postgres://testuser:testpass@localhost:36405/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36405/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:51 OK   000001_create_users_table.sql (11.25ms)
2025/08/01 15:45:51 OK   000002_create_projects_table.sql (12.85ms)
2025/08/01 15:45:51 OK   000003_create_log_entries_table.sql (19.72ms)
2025/08/01 15:45:51 OK   000004_create_tags_system.sql (16.46ms)
2025/08/01 15:45:51 OK   000005_create_auth_tables.sql (20.15ms)
2025/08/01 15:45:51 OK   000006_create_insights_table.sql (24.18ms)
2025/08/01 15:45:51 OK   000007_create_performance_indexes.sql (15.35ms)
2025/08/01 15:45:51 OK   000008_create_analytics_views.sql (31.92ms)
2025/08/01 15:45:51 OK   000009_development_data.sql (29.05ms)
2025/08/01 15:45:51 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36405/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36405/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:51 🐳 Stopping container: 86f8883f2ea1
2025/08/01 15:45:51 ✅ Container stopped: 86f8883f2ea1
2025/08/01 15:45:51 🐳 Terminating container: 86f8883f2ea1
2025/08/01 15:45:51 🚫 Container terminated: 86f8883f2ea1
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication
2025/08/01 15:45:51 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:51 ✅ Container created: c41ac9780a35
2025/08/01 15:45:51 🐳 Starting container: c41ac9780a35
2025/08/01 15:45:52 ✅ Container started: c41ac9780a35
2025/08/01 15:45:52 ⏳ Waiting for container id c41ac9780a35 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003af7e38 Strategies:[0xc00548a540 0xc00a081290]}
2025/08/01 15:45:53 🔔 Container is ready: c41ac9780a35
DSN: postgres://testuser:testpass@localhost:36406/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36406/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:54 OK   000001_create_users_table.sql (10.91ms)
2025/08/01 15:45:54 OK   000002_create_projects_table.sql (12.79ms)
2025/08/01 15:45:54 OK   000003_create_log_entries_table.sql (19.51ms)
2025/08/01 15:45:54 OK   000004_create_tags_system.sql (16.36ms)
2025/08/01 15:45:54 OK   000005_create_auth_tables.sql (20.28ms)
2025/08/01 15:45:54 OK   000006_create_insights_table.sql (24.39ms)
2025/08/01 15:45:54 OK   000007_create_performance_indexes.sql (16.23ms)
2025/08/01 15:45:54 OK   000008_create_analytics_views.sql (29.85ms)
2025/08/01 15:45:54 OK   000009_development_data.sql (29.91ms)
2025/08/01 15:45:54 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36406/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36406/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:54 🐳 Stopping container: c41ac9780a35
2025/08/01 15:45:54 ✅ Container stopped: c41ac9780a35
2025/08/01 15:45:54 🐳 Terminating container: c41ac9780a35
2025/08/01 15:45:54 🚫 Container terminated: c41ac9780a35
=== RUN   TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token
2025/08/01 15:45:54 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:54 ✅ Container created: 74e33687db63
2025/08/01 15:45:54 🐳 Starting container: 74e33687db63
2025/08/01 15:45:54 ✅ Container started: 74e33687db63
2025/08/01 15:45:54 ⏳ Waiting for container id 74e33687db63 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007b5b6d8 Strategies:[0xc0048fc5a0 0xc005aa6f30]}
2025/08/01 15:45:56 🔔 Container is ready: 74e33687db63
DSN: postgres://testuser:testpass@localhost:36407/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36407/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:56 OK   000001_create_users_table.sql (10.32ms)
2025/08/01 15:45:56 OK   000002_create_projects_table.sql (11.87ms)
2025/08/01 15:45:56 OK   000003_create_log_entries_table.sql (18.99ms)
2025/08/01 15:45:56 OK   000004_create_tags_system.sql (16.22ms)
2025/08/01 15:45:56 OK   000005_create_auth_tables.sql (19.82ms)
2025/08/01 15:45:56 OK   000006_create_insights_table.sql (24.57ms)
2025/08/01 15:45:56 OK   000007_create_performance_indexes.sql (15.14ms)
2025/08/01 15:45:56 OK   000008_create_analytics_views.sql (27.31ms)
2025/08/01 15:45:56 OK   000009_development_data.sql (28.07ms)
2025/08/01 15:45:56 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36407/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36407/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:56 🐳 Stopping container: 74e33687db63
2025/08/01 15:45:56 ✅ Container stopped: 74e33687db63
2025/08/01 15:45:56 🐳 Terminating container: 74e33687db63
2025/08/01 15:45:57 🚫 Container terminated: 74e33687db63
--- PASS: TestTagHandler_Comprehensive_AuthenticationRequired (46.06s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/POST_/v1/tags_without_authentication (2.54s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/POST_/v1/tags_with_invalid_token (2.55s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags_without_authentication (2.43s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags_with_invalid_token (2.53s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/popular_without_authentication (2.56s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/popular_with_invalid_token (2.50s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/recent_without_authentication (2.66s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/recent_with_invalid_token (2.58s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/search?q=test_without_authentication (2.60s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/search?q=test_with_invalid_token (2.55s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/usage_without_authentication (2.66s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/usage_with_invalid_token (2.59s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/PUT_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication (2.56s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/PUT_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token (2.57s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/DELETE_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication (2.50s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/DELETE_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token (2.56s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/550e8400-e29b-41d4-a716-446655440000_without_authentication (2.57s)
    --- PASS: TestTagHandler_Comprehensive_AuthenticationRequired/GET_/v1/tags/550e8400-e29b-41d4-a716-446655440000_with_invalid_token (2.54s)
=== RUN   TestTagHandler_CreateTag
=== RUN   TestTagHandler_CreateTag/successful_tag_creation
2025/08/01 15:45:57 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:57 ✅ Container created: d62678f3bbfc
2025/08/01 15:45:57 🐳 Starting container: d62678f3bbfc
2025/08/01 15:45:57 ✅ Container started: d62678f3bbfc
2025/08/01 15:45:57 ⏳ Waiting for container id d62678f3bbfc image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002bf7820 Strategies:[0xc009bdb920 0xc007f924b0]}
2025/08/01 15:45:59 🔔 Container is ready: d62678f3bbfc
DSN: postgres://testuser:testpass@localhost:36408/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36408/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:45:59 OK   000001_create_users_table.sql (10.89ms)
2025/08/01 15:45:59 OK   000002_create_projects_table.sql (12.7ms)
2025/08/01 15:45:59 OK   000003_create_log_entries_table.sql (19.52ms)
2025/08/01 15:45:59 OK   000004_create_tags_system.sql (16.24ms)
2025/08/01 15:45:59 OK   000005_create_auth_tables.sql (20.65ms)
2025/08/01 15:45:59 OK   000006_create_insights_table.sql (24.32ms)
2025/08/01 15:45:59 OK   000007_create_performance_indexes.sql (15.29ms)
2025/08/01 15:45:59 OK   000008_create_analytics_views.sql (27.75ms)
2025/08/01 15:45:59 OK   000009_development_data.sql (29.52ms)
2025/08/01 15:45:59 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36408/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36408/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:45:59 🐳 Stopping container: d62678f3bbfc
2025/08/01 15:45:59 ✅ Container stopped: d62678f3bbfc
2025/08/01 15:45:59 🐳 Terminating container: d62678f3bbfc
2025/08/01 15:45:59 🚫 Container terminated: d62678f3bbfc
=== RUN   TestTagHandler_CreateTag/invalid_request_body
2025/08/01 15:45:59 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:45:59 ✅ Container created: d2adcb1ba9aa
2025/08/01 15:45:59 🐳 Starting container: d2adcb1ba9aa
2025/08/01 15:46:00 ✅ Container started: d2adcb1ba9aa
2025/08/01 15:46:00 ⏳ Waiting for container id d2adcb1ba9aa image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001e7ca98 Strategies:[0xc006ac9980 0xc005a60ed0]}
2025/08/01 15:46:01 🔔 Container is ready: d2adcb1ba9aa
DSN: postgres://testuser:testpass@localhost:36409/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36409/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:01 OK   000001_create_users_table.sql (12.28ms)
2025/08/01 15:46:01 OK   000002_create_projects_table.sql (14.07ms)
2025/08/01 15:46:01 OK   000003_create_log_entries_table.sql (20.89ms)
2025/08/01 15:46:01 OK   000004_create_tags_system.sql (17.06ms)
2025/08/01 15:46:01 OK   000005_create_auth_tables.sql (20.46ms)
2025/08/01 15:46:02 OK   000006_create_insights_table.sql (24.11ms)
2025/08/01 15:46:02 OK   000007_create_performance_indexes.sql (15.16ms)
2025/08/01 15:46:02 OK   000008_create_analytics_views.sql (27.17ms)
2025/08/01 15:46:02 OK   000009_development_data.sql (28.53ms)
2025/08/01 15:46:02 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36409/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36409/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:02 🐳 Stopping container: d2adcb1ba9aa
2025/08/01 15:46:02 ✅ Container stopped: d2adcb1ba9aa
2025/08/01 15:46:02 🐳 Terminating container: d2adcb1ba9aa
2025/08/01 15:46:02 🚫 Container terminated: d2adcb1ba9aa
--- PASS: TestTagHandler_CreateTag (5.62s)
    --- PASS: TestTagHandler_CreateTag/successful_tag_creation (2.81s)
    --- PASS: TestTagHandler_CreateTag/invalid_request_body (2.81s)
=== RUN   TestTagHandler_GetTag
=== RUN   TestTagHandler_GetTag/successful_tag_retrieval
2025/08/01 15:46:02 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:02 ✅ Container created: 4137057a7a88
2025/08/01 15:46:02 🐳 Starting container: 4137057a7a88
2025/08/01 15:46:02 ✅ Container started: 4137057a7a88
2025/08/01 15:46:02 ⏳ Waiting for container id 4137057a7a88 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007f05a10 Strategies:[0xc006c079e0 0xc003e4abd0]}
2025/08/01 15:46:04 🔔 Container is ready: 4137057a7a88
DSN: postgres://testuser:testpass@localhost:36411/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36411/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:04 OK   000001_create_users_table.sql (11.07ms)
2025/08/01 15:46:04 OK   000002_create_projects_table.sql (12.71ms)
2025/08/01 15:46:04 OK   000003_create_log_entries_table.sql (20.18ms)
2025/08/01 15:46:04 OK   000004_create_tags_system.sql (17.08ms)
2025/08/01 15:46:04 OK   000005_create_auth_tables.sql (20.79ms)
2025/08/01 15:46:04 OK   000006_create_insights_table.sql (25.01ms)
2025/08/01 15:46:04 OK   000007_create_performance_indexes.sql (15.67ms)
2025/08/01 15:46:04 OK   000008_create_analytics_views.sql (28.59ms)
2025/08/01 15:46:04 OK   000009_development_data.sql (29.4ms)
2025/08/01 15:46:04 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36411/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36411/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:05 🐳 Stopping container: 4137057a7a88
2025/08/01 15:46:05 ✅ Container stopped: 4137057a7a88
2025/08/01 15:46:05 🐳 Terminating container: 4137057a7a88
2025/08/01 15:46:05 🚫 Container terminated: 4137057a7a88
=== RUN   TestTagHandler_GetTag/tag_not_found
2025/08/01 15:46:05 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:05 ✅ Container created: cf070858dd95
2025/08/01 15:46:05 🐳 Starting container: cf070858dd95
2025/08/01 15:46:05 ✅ Container started: cf070858dd95
2025/08/01 15:46:05 ⏳ Waiting for container id cf070858dd95 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0069e1520 Strategies:[0xc008402ae0 0xc0037f0390]}
2025/08/01 15:46:07 🔔 Container is ready: cf070858dd95
DSN: postgres://testuser:testpass@localhost:36412/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36412/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:07 OK   000001_create_users_table.sql (10.48ms)
2025/08/01 15:46:07 OK   000002_create_projects_table.sql (12.25ms)
2025/08/01 15:46:07 OK   000003_create_log_entries_table.sql (20.24ms)
2025/08/01 15:46:07 OK   000004_create_tags_system.sql (16.38ms)
2025/08/01 15:46:07 OK   000005_create_auth_tables.sql (20.52ms)
2025/08/01 15:46:07 OK   000006_create_insights_table.sql (25.02ms)
2025/08/01 15:46:07 OK   000007_create_performance_indexes.sql (15.69ms)
2025/08/01 15:46:07 OK   000008_create_analytics_views.sql (28.97ms)
2025/08/01 15:46:07 OK   000009_development_data.sql (29.25ms)
2025/08/01 15:46:07 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36412/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36412/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:08 🐳 Stopping container: cf070858dd95
2025/08/01 15:46:08 ✅ Container stopped: cf070858dd95
2025/08/01 15:46:08 🐳 Terminating container: cf070858dd95
2025/08/01 15:46:08 🚫 Container terminated: cf070858dd95
--- PASS: TestTagHandler_GetTag (5.65s)
    --- PASS: TestTagHandler_GetTag/successful_tag_retrieval (2.89s)
    --- PASS: TestTagHandler_GetTag/tag_not_found (2.76s)
=== RUN   TestTagHandler_GetTags
=== RUN   TestTagHandler_GetTags/successful_tags_retrieval
2025/08/01 15:46:08 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:08 ✅ Container created: 34a8cb3386e5
2025/08/01 15:46:08 🐳 Starting container: 34a8cb3386e5
2025/08/01 15:46:08 ✅ Container started: 34a8cb3386e5
2025/08/01 15:46:08 ⏳ Waiting for container id 34a8cb3386e5 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001185790 Strategies:[0xc002d6b0e0 0xc009e85650]}
2025/08/01 15:46:10 🔔 Container is ready: 34a8cb3386e5
DSN: postgres://testuser:testpass@localhost:36413/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36413/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:10 OK   000001_create_users_table.sql (10.66ms)
2025/08/01 15:46:10 OK   000002_create_projects_table.sql (12.57ms)
2025/08/01 15:46:10 OK   000003_create_log_entries_table.sql (19.56ms)
2025/08/01 15:46:10 OK   000004_create_tags_system.sql (16.26ms)
2025/08/01 15:46:10 OK   000005_create_auth_tables.sql (20.43ms)
2025/08/01 15:46:10 OK   000006_create_insights_table.sql (24.41ms)
2025/08/01 15:46:10 OK   000007_create_performance_indexes.sql (15.4ms)
2025/08/01 15:46:10 OK   000008_create_analytics_views.sql (28.12ms)
2025/08/01 15:46:10 OK   000009_development_data.sql (28.86ms)
2025/08/01 15:46:10 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36413/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36413/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:10 🐳 Stopping container: 34a8cb3386e5
2025/08/01 15:46:10 ✅ Container stopped: 34a8cb3386e5
2025/08/01 15:46:10 🐳 Terminating container: 34a8cb3386e5
2025/08/01 15:46:11 🚫 Container terminated: 34a8cb3386e5
--- PASS: TestTagHandler_GetTags (2.72s)
    --- PASS: TestTagHandler_GetTags/successful_tags_retrieval (2.72s)
=== RUN   TestTagHandler_GetPopularTags
=== RUN   TestTagHandler_GetPopularTags/successful_popular_tags_retrieval
2025/08/01 15:46:11 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:11 ✅ Container created: f040b7e775f0
2025/08/01 15:46:11 🐳 Starting container: f040b7e775f0
2025/08/01 15:46:11 ✅ Container started: f040b7e775f0
2025/08/01 15:46:11 ⏳ Waiting for container id f040b7e775f0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006b32988 Strategies:[0xc006ce2600 0xc006c50cf0]}
2025/08/01 15:46:13 🔔 Container is ready: f040b7e775f0
DSN: postgres://testuser:testpass@localhost:36414/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36414/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:13 OK   000001_create_users_table.sql (10.41ms)
2025/08/01 15:46:13 OK   000002_create_projects_table.sql (12.65ms)
2025/08/01 15:46:13 OK   000003_create_log_entries_table.sql (19.39ms)
2025/08/01 15:46:13 OK   000004_create_tags_system.sql (16.17ms)
2025/08/01 15:46:13 OK   000005_create_auth_tables.sql (20.06ms)
2025/08/01 15:46:13 OK   000006_create_insights_table.sql (23.74ms)
2025/08/01 15:46:13 OK   000007_create_performance_indexes.sql (16.15ms)
2025/08/01 15:46:13 OK   000008_create_analytics_views.sql (27.85ms)
2025/08/01 15:46:13 OK   000009_development_data.sql (27.72ms)
2025/08/01 15:46:13 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36414/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36414/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:13 🐳 Stopping container: f040b7e775f0
2025/08/01 15:46:13 ✅ Container stopped: f040b7e775f0
2025/08/01 15:46:13 🐳 Terminating container: f040b7e775f0
2025/08/01 15:46:13 🚫 Container terminated: f040b7e775f0
=== RUN   TestTagHandler_GetPopularTags/popular_tags_with_limit_parameter
2025/08/01 15:46:13 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:13 ✅ Container created: af4c6c1f66eb
2025/08/01 15:46:13 🐳 Starting container: af4c6c1f66eb
2025/08/01 15:46:14 ✅ Container started: af4c6c1f66eb
2025/08/01 15:46:14 ⏳ Waiting for container id af4c6c1f66eb image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0083a5570 Strategies:[0xc00480fbc0 0xc009918540]}
2025/08/01 15:46:15 🔔 Container is ready: af4c6c1f66eb
DSN: postgres://testuser:testpass@localhost:36415/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36415/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:15 OK   000001_create_users_table.sql (11ms)
2025/08/01 15:46:15 OK   000002_create_projects_table.sql (12.77ms)
2025/08/01 15:46:15 OK   000003_create_log_entries_table.sql (19.58ms)
2025/08/01 15:46:15 OK   000004_create_tags_system.sql (16.78ms)
2025/08/01 15:46:15 OK   000005_create_auth_tables.sql (20.93ms)
2025/08/01 15:46:15 OK   000006_create_insights_table.sql (25.49ms)
2025/08/01 15:46:16 OK   000007_create_performance_indexes.sql (15.65ms)
2025/08/01 15:46:16 OK   000008_create_analytics_views.sql (28.59ms)
2025/08/01 15:46:16 OK   000009_development_data.sql (29.88ms)
2025/08/01 15:46:16 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36415/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36415/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:16 🐳 Stopping container: af4c6c1f66eb
2025/08/01 15:46:16 ✅ Container stopped: af4c6c1f66eb
2025/08/01 15:46:16 🐳 Terminating container: af4c6c1f66eb
2025/08/01 15:46:16 🚫 Container terminated: af4c6c1f66eb
--- PASS: TestTagHandler_GetPopularTags (5.61s)
    --- PASS: TestTagHandler_GetPopularTags/successful_popular_tags_retrieval (2.80s)
    --- PASS: TestTagHandler_GetPopularTags/popular_tags_with_limit_parameter (2.81s)
=== RUN   TestTagHandler_GetRecentlyUsedTags
=== RUN   TestTagHandler_GetRecentlyUsedTags/successful_recently_used_tags_retrieval
2025/08/01 15:46:16 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:16 ✅ Container created: 8af631a2f6ab
2025/08/01 15:46:16 🐳 Starting container: 8af631a2f6ab
2025/08/01 15:46:16 ✅ Container started: 8af631a2f6ab
2025/08/01 15:46:16 ⏳ Waiting for container id 8af631a2f6ab image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006a83338 Strategies:[0xc007b47140 0xc007c6d140]}
2025/08/01 15:46:18 🔔 Container is ready: 8af631a2f6ab
DSN: postgres://testuser:testpass@localhost:36416/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36416/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:18 OK   000001_create_users_table.sql (10.51ms)
2025/08/01 15:46:18 OK   000002_create_projects_table.sql (12.1ms)
2025/08/01 15:46:18 OK   000003_create_log_entries_table.sql (18.48ms)
2025/08/01 15:46:18 OK   000004_create_tags_system.sql (14.5ms)
2025/08/01 15:46:18 OK   000005_create_auth_tables.sql (17.42ms)
2025/08/01 15:46:18 OK   000006_create_insights_table.sql (20.42ms)
2025/08/01 15:46:18 OK   000007_create_performance_indexes.sql (12.7ms)
2025/08/01 15:46:18 OK   000008_create_analytics_views.sql (23.37ms)
2025/08/01 15:46:18 OK   000009_development_data.sql (24.5ms)
2025/08/01 15:46:18 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36416/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36416/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:19 🐳 Stopping container: 8af631a2f6ab
2025/08/01 15:46:19 ✅ Container stopped: 8af631a2f6ab
2025/08/01 15:46:19 🐳 Terminating container: 8af631a2f6ab
2025/08/01 15:46:19 🚫 Container terminated: 8af631a2f6ab
--- PASS: TestTagHandler_GetRecentlyUsedTags (2.77s)
    --- PASS: TestTagHandler_GetRecentlyUsedTags/successful_recently_used_tags_retrieval (2.77s)
=== RUN   TestTagHandler_SearchTags
=== RUN   TestTagHandler_SearchTags/successful_tag_search
2025/08/01 15:46:19 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:19 ✅ Container created: 18e7b73d4aa7
2025/08/01 15:46:19 🐳 Starting container: 18e7b73d4aa7
2025/08/01 15:46:19 ✅ Container started: 18e7b73d4aa7
2025/08/01 15:46:19 ⏳ Waiting for container id 18e7b73d4aa7 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009e696e8 Strategies:[0xc00434fc20 0xc009e39890]}
2025/08/01 15:46:21 🔔 Container is ready: 18e7b73d4aa7
DSN: postgres://testuser:testpass@localhost:36417/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36417/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:21 OK   000001_create_users_table.sql (10.12ms)
2025/08/01 15:46:21 OK   000002_create_projects_table.sql (11.32ms)
2025/08/01 15:46:21 OK   000003_create_log_entries_table.sql (18.5ms)
2025/08/01 15:46:21 OK   000004_create_tags_system.sql (16.1ms)
2025/08/01 15:46:21 OK   000005_create_auth_tables.sql (20.14ms)
2025/08/01 15:46:21 OK   000006_create_insights_table.sql (24.07ms)
2025/08/01 15:46:21 OK   000007_create_performance_indexes.sql (15.17ms)
2025/08/01 15:46:21 OK   000008_create_analytics_views.sql (27.42ms)
2025/08/01 15:46:21 OK   000009_development_data.sql (28.34ms)
2025/08/01 15:46:21 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36417/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36417/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:21 🐳 Stopping container: 18e7b73d4aa7
2025/08/01 15:46:22 ✅ Container stopped: 18e7b73d4aa7
2025/08/01 15:46:22 🐳 Terminating container: 18e7b73d4aa7
2025/08/01 15:46:22 🚫 Container terminated: 18e7b73d4aa7
=== RUN   TestTagHandler_SearchTags/search_with_missing_query_parameter
2025/08/01 15:46:22 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:22 ✅ Container created: 4b1bbdd96af3
2025/08/01 15:46:22 🐳 Starting container: 4b1bbdd96af3
2025/08/01 15:46:22 ✅ Container started: 4b1bbdd96af3
2025/08/01 15:46:22 ⏳ Waiting for container id 4b1bbdd96af3 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009fc8a40 Strategies:[0xc001815800 0xc001f96420]}
2025/08/01 15:46:24 🔔 Container is ready: 4b1bbdd96af3
DSN: postgres://testuser:testpass@localhost:36419/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36419/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:24 OK   000001_create_users_table.sql (11.32ms)
2025/08/01 15:46:24 OK   000002_create_projects_table.sql (13.27ms)
2025/08/01 15:46:24 OK   000003_create_log_entries_table.sql (20.8ms)
2025/08/01 15:46:24 OK   000004_create_tags_system.sql (17.78ms)
2025/08/01 15:46:24 OK   000005_create_auth_tables.sql (21.17ms)
2025/08/01 15:46:24 OK   000006_create_insights_table.sql (25.67ms)
2025/08/01 15:46:24 OK   000007_create_performance_indexes.sql (15.93ms)
2025/08/01 15:46:24 OK   000008_create_analytics_views.sql (28.94ms)
2025/08/01 15:46:24 OK   000009_development_data.sql (30.5ms)
2025/08/01 15:46:24 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36419/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36419/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:24 🐳 Stopping container: 4b1bbdd96af3
2025/08/01 15:46:24 ✅ Container stopped: 4b1bbdd96af3
2025/08/01 15:46:24 🐳 Terminating container: 4b1bbdd96af3
2025/08/01 15:46:25 🚫 Container terminated: 4b1bbdd96af3
--- PASS: TestTagHandler_SearchTags (5.61s)
    --- PASS: TestTagHandler_SearchTags/successful_tag_search (2.79s)
    --- PASS: TestTagHandler_SearchTags/search_with_missing_query_parameter (2.82s)
=== RUN   TestTagHandler_UpdateTag
=== RUN   TestTagHandler_UpdateTag/successful_tag_update
2025/08/01 15:46:25 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:25 ✅ Container created: c7cfcef3e2ed
2025/08/01 15:46:25 🐳 Starting container: c7cfcef3e2ed
2025/08/01 15:46:25 ✅ Container started: c7cfcef3e2ed
2025/08/01 15:46:25 ⏳ Waiting for container id c7cfcef3e2ed image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc005435170 Strategies:[0xc006ac8ea0 0xc009eca990]}
2025/08/01 15:46:27 🔔 Container is ready: c7cfcef3e2ed
DSN: postgres://testuser:testpass@localhost:36420/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36420/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:27 OK   000001_create_users_table.sql (10.92ms)
2025/08/01 15:46:27 OK   000002_create_projects_table.sql (13.18ms)
2025/08/01 15:46:27 OK   000003_create_log_entries_table.sql (20ms)
2025/08/01 15:46:27 OK   000004_create_tags_system.sql (16.76ms)
2025/08/01 15:46:27 OK   000005_create_auth_tables.sql (20.55ms)
2025/08/01 15:46:27 OK   000006_create_insights_table.sql (24.52ms)
2025/08/01 15:46:27 OK   000007_create_performance_indexes.sql (15.15ms)
2025/08/01 15:46:27 OK   000008_create_analytics_views.sql (27.65ms)
2025/08/01 15:46:27 OK   000009_development_data.sql (28.6ms)
2025/08/01 15:46:27 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36420/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36420/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:27 🐳 Stopping container: c7cfcef3e2ed
2025/08/01 15:46:27 ✅ Container stopped: c7cfcef3e2ed
2025/08/01 15:46:27 🐳 Terminating container: c7cfcef3e2ed
2025/08/01 15:46:27 🚫 Container terminated: c7cfcef3e2ed
=== RUN   TestTagHandler_UpdateTag/update_non-existent_tag
2025/08/01 15:46:27 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:27 ✅ Container created: b0f7faf53b40
2025/08/01 15:46:27 🐳 Starting container: b0f7faf53b40
2025/08/01 15:46:28 ✅ Container started: b0f7faf53b40
2025/08/01 15:46:28 ⏳ Waiting for container id b0f7faf53b40 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc003abff10 Strategies:[0xc006d05980 0xc004618300]}
2025/08/01 15:46:29 🔔 Container is ready: b0f7faf53b40
DSN: postgres://testuser:testpass@localhost:36421/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36421/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:29 OK   000001_create_users_table.sql (11.73ms)
2025/08/01 15:46:29 OK   000002_create_projects_table.sql (13.86ms)
2025/08/01 15:46:30 OK   000003_create_log_entries_table.sql (21.42ms)
2025/08/01 15:46:30 OK   000004_create_tags_system.sql (19.09ms)
2025/08/01 15:46:30 OK   000005_create_auth_tables.sql (21.39ms)
2025/08/01 15:46:30 OK   000006_create_insights_table.sql (26.04ms)
2025/08/01 15:46:30 OK   000007_create_performance_indexes.sql (16.12ms)
2025/08/01 15:46:30 OK   000008_create_analytics_views.sql (29.35ms)
2025/08/01 15:46:30 OK   000009_development_data.sql (29.75ms)
2025/08/01 15:46:30 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36421/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36421/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:30 🐳 Stopping container: b0f7faf53b40
2025/08/01 15:46:30 ✅ Container stopped: b0f7faf53b40
2025/08/01 15:46:30 🐳 Terminating container: b0f7faf53b40
2025/08/01 15:46:30 🚫 Container terminated: b0f7faf53b40
--- PASS: TestTagHandler_UpdateTag (5.72s)
    --- PASS: TestTagHandler_UpdateTag/successful_tag_update (2.89s)
    --- PASS: TestTagHandler_UpdateTag/update_non-existent_tag (2.83s)
=== RUN   TestTagHandler_DeleteTag
=== RUN   TestTagHandler_DeleteTag/successful_tag_deletion
2025/08/01 15:46:30 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:30 ✅ Container created: 5af373b0bbbc
2025/08/01 15:46:30 🐳 Starting container: 5af373b0bbbc
2025/08/01 15:46:30 ✅ Container started: 5af373b0bbbc
2025/08/01 15:46:30 ⏳ Waiting for container id 5af373b0bbbc image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc005b7d470 Strategies:[0xc009f710e0 0xc0065091a0]}
2025/08/01 15:46:32 🔔 Container is ready: 5af373b0bbbc
DSN: postgres://testuser:testpass@localhost:36422/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36422/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:32 OK   000001_create_users_table.sql (9.62ms)
2025/08/01 15:46:32 OK   000002_create_projects_table.sql (11.61ms)
2025/08/01 15:46:32 OK   000003_create_log_entries_table.sql (18.59ms)
2025/08/01 15:46:32 OK   000004_create_tags_system.sql (16.18ms)
2025/08/01 15:46:32 OK   000005_create_auth_tables.sql (20.2ms)
2025/08/01 15:46:32 OK   000006_create_insights_table.sql (24.56ms)
2025/08/01 15:46:32 OK   000007_create_performance_indexes.sql (15.95ms)
2025/08/01 15:46:32 OK   000008_create_analytics_views.sql (28.59ms)
2025/08/01 15:46:32 OK   000009_development_data.sql (28.83ms)
2025/08/01 15:46:32 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36422/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36422/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:33 🐳 Stopping container: 5af373b0bbbc
2025/08/01 15:46:33 ✅ Container stopped: 5af373b0bbbc
2025/08/01 15:46:33 🐳 Terminating container: 5af373b0bbbc
2025/08/01 15:46:33 🚫 Container terminated: 5af373b0bbbc
=== RUN   TestTagHandler_DeleteTag/delete_non-existent_tag
2025/08/01 15:46:33 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:33 ✅ Container created: 375f1b1b5ccd
2025/08/01 15:46:33 🐳 Starting container: 375f1b1b5ccd
2025/08/01 15:46:33 ✅ Container started: 375f1b1b5ccd
2025/08/01 15:46:33 ⏳ Waiting for container id 375f1b1b5ccd image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006834e90 Strategies:[0xc00a15d260 0xc007d3b200]}
2025/08/01 15:46:35 🔔 Container is ready: 375f1b1b5ccd
DSN: postgres://testuser:testpass@localhost:36423/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36423/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:35 OK   000001_create_users_table.sql (11.01ms)
2025/08/01 15:46:35 OK   000002_create_projects_table.sql (12.65ms)
2025/08/01 15:46:35 OK   000003_create_log_entries_table.sql (19.68ms)
2025/08/01 15:46:35 OK   000004_create_tags_system.sql (16.54ms)
2025/08/01 15:46:35 OK   000005_create_auth_tables.sql (20.31ms)
2025/08/01 15:46:35 OK   000006_create_insights_table.sql (24.15ms)
2025/08/01 15:46:35 OK   000007_create_performance_indexes.sql (15.4ms)
2025/08/01 15:46:35 OK   000008_create_analytics_views.sql (27.69ms)
2025/08/01 15:46:35 OK   000009_development_data.sql (27.92ms)
2025/08/01 15:46:35 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36423/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36423/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:36 🐳 Stopping container: 375f1b1b5ccd
2025/08/01 15:46:36 ✅ Container stopped: 375f1b1b5ccd
2025/08/01 15:46:36 🐳 Terminating container: 375f1b1b5ccd
2025/08/01 15:46:36 🚫 Container terminated: 375f1b1b5ccd
--- PASS: TestTagHandler_DeleteTag (5.66s)
    --- PASS: TestTagHandler_DeleteTag/successful_tag_deletion (2.86s)
    --- PASS: TestTagHandler_DeleteTag/delete_non-existent_tag (2.80s)
=== RUN   TestTagHandler_GetUserTagUsage
=== RUN   TestTagHandler_GetUserTagUsage/successful_user_tag_usage_retrieval
2025/08/01 15:46:36 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:36 ✅ Container created: ee31d0b22573
2025/08/01 15:46:36 🐳 Starting container: ee31d0b22573
2025/08/01 15:46:36 ✅ Container started: ee31d0b22573
2025/08/01 15:46:36 ⏳ Waiting for container id ee31d0b22573 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006ad30f0 Strategies:[0xc003f27260 0xc009bb4c60]}
2025/08/01 15:46:38 🔔 Container is ready: ee31d0b22573
DSN: postgres://testuser:testpass@localhost:36425/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36425/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:38 OK   000001_create_users_table.sql (10.03ms)
2025/08/01 15:46:38 OK   000002_create_projects_table.sql (11.75ms)
2025/08/01 15:46:38 OK   000003_create_log_entries_table.sql (18.07ms)
2025/08/01 15:46:38 OK   000004_create_tags_system.sql (14.45ms)
2025/08/01 15:46:38 OK   000005_create_auth_tables.sql (16.94ms)
2025/08/01 15:46:38 OK   000006_create_insights_table.sql (20.42ms)
2025/08/01 15:46:38 OK   000007_create_performance_indexes.sql (12.91ms)
2025/08/01 15:46:38 OK   000008_create_analytics_views.sql (23.31ms)
2025/08/01 15:46:38 OK   000009_development_data.sql (24.48ms)
2025/08/01 15:46:38 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36425/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36425/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:38 🐳 Stopping container: ee31d0b22573
2025/08/01 15:46:39 ✅ Container stopped: ee31d0b22573
2025/08/01 15:46:39 🐳 Terminating container: ee31d0b22573
2025/08/01 15:46:39 🚫 Container terminated: ee31d0b22573
--- PASS: TestTagHandler_GetUserTagUsage (2.77s)
    --- PASS: TestTagHandler_GetUserTagUsage/successful_user_tag_usage_retrieval (2.77s)
=== RUN   TestTagHandler_ErrorHandling
=== RUN   TestTagHandler_ErrorHandling/unauthorized_create_tag_request
2025/08/01 15:46:39 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:39 ✅ Container created: 783a424f8d83
2025/08/01 15:46:39 🐳 Starting container: 783a424f8d83
2025/08/01 15:46:39 ✅ Container started: 783a424f8d83
2025/08/01 15:46:39 ⏳ Waiting for container id 783a424f8d83 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004256198 Strategies:[0xc003ede660 0xc005fd87b0]}
2025/08/01 15:46:41 🔔 Container is ready: 783a424f8d83
DSN: postgres://testuser:testpass@localhost:36426/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36426/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:41 OK   000001_create_users_table.sql (10.7ms)
2025/08/01 15:46:41 OK   000002_create_projects_table.sql (12.32ms)
2025/08/01 15:46:41 OK   000003_create_log_entries_table.sql (19.07ms)
2025/08/01 15:46:41 OK   000004_create_tags_system.sql (15.79ms)
2025/08/01 15:46:41 OK   000005_create_auth_tables.sql (19.84ms)
2025/08/01 15:46:41 OK   000006_create_insights_table.sql (24.2ms)
2025/08/01 15:46:41 OK   000007_create_performance_indexes.sql (15.36ms)
2025/08/01 15:46:41 OK   000008_create_analytics_views.sql (28.06ms)
2025/08/01 15:46:41 OK   000009_development_data.sql (28.55ms)
2025/08/01 15:46:41 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36426/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36426/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:41 🐳 Stopping container: 783a424f8d83
2025/08/01 15:46:41 ✅ Container stopped: 783a424f8d83
2025/08/01 15:46:41 🐳 Terminating container: 783a424f8d83
2025/08/01 15:46:41 🚫 Container terminated: 783a424f8d83
=== RUN   TestTagHandler_ErrorHandling/unauthorized_get_tag_request
2025/08/01 15:46:41 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:41 ✅ Container created: e692ef611f1c
2025/08/01 15:46:41 🐳 Starting container: e692ef611f1c
2025/08/01 15:46:41 ✅ Container started: e692ef611f1c
2025/08/01 15:46:41 ⏳ Waiting for container id e692ef611f1c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0098ec298 Strategies:[0xc006c127e0 0xc00989c2d0]}
2025/08/01 15:46:43 🔔 Container is ready: e692ef611f1c
DSN: postgres://testuser:testpass@localhost:36427/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36427/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:43 OK   000001_create_users_table.sql (11.13ms)
2025/08/01 15:46:43 OK   000002_create_projects_table.sql (12.24ms)
2025/08/01 15:46:43 OK   000003_create_log_entries_table.sql (18.45ms)
2025/08/01 15:46:43 OK   000004_create_tags_system.sql (16.11ms)
2025/08/01 15:46:43 OK   000005_create_auth_tables.sql (19.87ms)
2025/08/01 15:46:43 OK   000006_create_insights_table.sql (24.23ms)
2025/08/01 15:46:43 OK   000007_create_performance_indexes.sql (15.07ms)
2025/08/01 15:46:43 OK   000008_create_analytics_views.sql (28.03ms)
2025/08/01 15:46:43 OK   000009_development_data.sql (28.56ms)
2025/08/01 15:46:43 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36427/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36427/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:43 🐳 Stopping container: e692ef611f1c
2025/08/01 15:46:44 ✅ Container stopped: e692ef611f1c
2025/08/01 15:46:44 🐳 Terminating container: e692ef611f1c
2025/08/01 15:46:44 🚫 Container terminated: e692ef611f1c
=== RUN   TestTagHandler_ErrorHandling/invalid_tag_ID_format
2025/08/01 15:46:44 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:44 ✅ Container created: 845a73214e67
2025/08/01 15:46:44 🐳 Starting container: 845a73214e67
2025/08/01 15:46:44 ✅ Container started: 845a73214e67
2025/08/01 15:46:44 ⏳ Waiting for container id 845a73214e67 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00b18a2b8 Strategies:[0xc0069baf00 0xc006728900]}
2025/08/01 15:46:46 🔔 Container is ready: 845a73214e67
DSN: postgres://testuser:testpass@localhost:36428/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36428/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:46 OK   000001_create_users_table.sql (11.2ms)
2025/08/01 15:46:46 OK   000002_create_projects_table.sql (13.42ms)
2025/08/01 15:46:46 OK   000003_create_log_entries_table.sql (20.49ms)
2025/08/01 15:46:46 OK   000004_create_tags_system.sql (17.14ms)
2025/08/01 15:46:46 OK   000005_create_auth_tables.sql (20.6ms)
2025/08/01 15:46:46 OK   000006_create_insights_table.sql (25.3ms)
2025/08/01 15:46:46 OK   000007_create_performance_indexes.sql (16.34ms)
2025/08/01 15:46:46 OK   000008_create_analytics_views.sql (29.19ms)
2025/08/01 15:46:46 OK   000009_development_data.sql (29.39ms)
2025/08/01 15:46:46 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36428/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36428/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:46 🐳 Stopping container: 845a73214e67
2025/08/01 15:46:46 ✅ Container stopped: 845a73214e67
2025/08/01 15:46:46 🐳 Terminating container: 845a73214e67
2025/08/01 15:46:47 🚫 Container terminated: 845a73214e67
--- PASS: TestTagHandler_ErrorHandling (7.87s)
    --- PASS: TestTagHandler_ErrorHandling/unauthorized_create_tag_request (2.54s)
    --- PASS: TestTagHandler_ErrorHandling/unauthorized_get_tag_request (2.55s)
    --- PASS: TestTagHandler_ErrorHandling/invalid_tag_ID_format (2.77s)
=== RUN   TestUserHandler_ComprehensiveScenarios
=== RUN   TestUserHandler_ComprehensiveScenarios/complete_user_profile_lifecycle
2025/08/01 15:46:47 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:47 ✅ Container created: 6277dea2d96e
2025/08/01 15:46:47 🐳 Starting container: 6277dea2d96e
2025/08/01 15:46:47 ✅ Container started: 6277dea2d96e
2025/08/01 15:46:47 ⏳ Waiting for container id 6277dea2d96e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0069d4ae0 Strategies:[0xc006aa48a0 0xc0041e2ae0]}
2025/08/01 15:46:49 🔔 Container is ready: 6277dea2d96e
DSN: postgres://testuser:testpass@localhost:36429/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36429/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:49 OK   000001_create_users_table.sql (10.97ms)
2025/08/01 15:46:49 OK   000002_create_projects_table.sql (12.64ms)
2025/08/01 15:46:49 OK   000003_create_log_entries_table.sql (20.26ms)
2025/08/01 15:46:49 OK   000004_create_tags_system.sql (17.09ms)
2025/08/01 15:46:49 OK   000005_create_auth_tables.sql (21.12ms)
2025/08/01 15:46:49 OK   000006_create_insights_table.sql (25.27ms)
2025/08/01 15:46:49 OK   000007_create_performance_indexes.sql (15.84ms)
2025/08/01 15:46:49 OK   000008_create_analytics_views.sql (28.98ms)
2025/08/01 15:46:49 OK   000009_development_data.sql (29.95ms)
2025/08/01 15:46:49 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36429/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36429/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:49 🐳 Stopping container: 6277dea2d96e
2025/08/01 15:46:50 ✅ Container stopped: 6277dea2d96e
2025/08/01 15:46:50 🐳 Terminating container: 6277dea2d96e
2025/08/01 15:46:50 🚫 Container terminated: 6277dea2d96e
=== RUN   TestUserHandler_ComprehensiveScenarios/profile_update_validation
2025/08/01 15:46:50 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:50 ✅ Container created: e27cdb53f962
2025/08/01 15:46:50 🐳 Starting container: e27cdb53f962
2025/08/01 15:46:50 ✅ Container started: e27cdb53f962
2025/08/01 15:46:50 ⏳ Waiting for container id e27cdb53f962 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009f216e0 Strategies:[0xc00542a480 0xc0054dca20]}
2025/08/01 15:46:52 🔔 Container is ready: e27cdb53f962
DSN: postgres://testuser:testpass@localhost:36430/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36430/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:52 OK   000001_create_users_table.sql (10.63ms)
2025/08/01 15:46:52 OK   000002_create_projects_table.sql (12.65ms)
2025/08/01 15:46:52 OK   000003_create_log_entries_table.sql (20.14ms)
2025/08/01 15:46:52 OK   000004_create_tags_system.sql (17.05ms)
2025/08/01 15:46:52 OK   000005_create_auth_tables.sql (20.42ms)
2025/08/01 15:46:52 OK   000006_create_insights_table.sql (24.64ms)
2025/08/01 15:46:52 OK   000007_create_performance_indexes.sql (15.32ms)
2025/08/01 15:46:52 OK   000008_create_analytics_views.sql (28.52ms)
2025/08/01 15:46:52 OK   000009_development_data.sql (28.37ms)
2025/08/01 15:46:52 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36430/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36430/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:52 🐳 Stopping container: e27cdb53f962
2025/08/01 15:46:52 ✅ Container stopped: e27cdb53f962
2025/08/01 15:46:52 🐳 Terminating container: e27cdb53f962
2025/08/01 15:46:52 🚫 Container terminated: e27cdb53f962
=== RUN   TestUserHandler_ComprehensiveScenarios/password_change_security
2025/08/01 15:46:52 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:53 ✅ Container created: 64dde413207b
2025/08/01 15:46:53 🐳 Starting container: 64dde413207b
2025/08/01 15:46:53 ✅ Container started: 64dde413207b
2025/08/01 15:46:53 ⏳ Waiting for container id 64dde413207b image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0071edb98 Strategies:[0xc00a07dda0 0xc009ca3e90]}
2025/08/01 15:46:54 🔔 Container is ready: 64dde413207b
DSN: postgres://testuser:testpass@localhost:36431/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36431/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:55 OK   000001_create_users_table.sql (11.38ms)
2025/08/01 15:46:55 OK   000002_create_projects_table.sql (13.13ms)
2025/08/01 15:46:55 OK   000003_create_log_entries_table.sql (20.71ms)
2025/08/01 15:46:55 OK   000004_create_tags_system.sql (16.81ms)
2025/08/01 15:46:55 OK   000005_create_auth_tables.sql (21.06ms)
2025/08/01 15:46:55 OK   000006_create_insights_table.sql (25.11ms)
2025/08/01 15:46:55 OK   000007_create_performance_indexes.sql (15.8ms)
2025/08/01 15:46:55 OK   000008_create_analytics_views.sql (28.72ms)
2025/08/01 15:46:55 OK   000009_development_data.sql (28.84ms)
2025/08/01 15:46:55 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36431/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36431/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:55 🐳 Stopping container: 64dde413207b
2025/08/01 15:46:56 ✅ Container stopped: 64dde413207b
2025/08/01 15:46:56 🐳 Terminating container: 64dde413207b
2025/08/01 15:46:56 🚫 Container terminated: 64dde413207b
=== RUN   TestUserHandler_ComprehensiveScenarios/concurrent_operations
2025/08/01 15:46:56 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:56 ✅ Container created: 222527ef3e35
2025/08/01 15:46:56 🐳 Starting container: 222527ef3e35
2025/08/01 15:46:56 ✅ Container started: 222527ef3e35
2025/08/01 15:46:56 ⏳ Waiting for container id 222527ef3e35 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007e9c668 Strategies:[0xc00434b7a0 0xc0053792f0]}
2025/08/01 15:46:58 🔔 Container is ready: 222527ef3e35
DSN: postgres://testuser:testpass@localhost:36432/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36432/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:46:58 OK   000001_create_users_table.sql (10.48ms)
2025/08/01 15:46:58 OK   000002_create_projects_table.sql (12.47ms)
2025/08/01 15:46:58 OK   000003_create_log_entries_table.sql (19.2ms)
2025/08/01 15:46:58 OK   000004_create_tags_system.sql (16.85ms)
2025/08/01 15:46:58 OK   000005_create_auth_tables.sql (20.48ms)
2025/08/01 15:46:58 OK   000006_create_insights_table.sql (24.53ms)
2025/08/01 15:46:58 OK   000007_create_performance_indexes.sql (15.86ms)
2025/08/01 15:46:58 OK   000008_create_analytics_views.sql (28.19ms)
2025/08/01 15:46:58 OK   000009_development_data.sql (28.25ms)
2025/08/01 15:46:58 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36432/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36432/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:46:58 🐳 Stopping container: 222527ef3e35
2025/08/01 15:46:58 ✅ Container stopped: 222527ef3e35
2025/08/01 15:46:58 🐳 Terminating container: 222527ef3e35
2025/08/01 15:46:58 🚫 Container terminated: 222527ef3e35
=== RUN   TestUserHandler_ComprehensiveScenarios/malicious_input_handling
2025/08/01 15:46:58 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:46:59 ✅ Container created: abb53594b12c
2025/08/01 15:46:59 🐳 Starting container: abb53594b12c
2025/08/01 15:46:59 ✅ Container started: abb53594b12c
2025/08/01 15:46:59 ⏳ Waiting for container id abb53594b12c image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009cd8538 Strategies:[0xc00680e3c0 0xc005abf6e0]}
2025/08/01 15:47:00 🔔 Container is ready: abb53594b12c
DSN: postgres://testuser:testpass@localhost:36433/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36433/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:01 OK   000001_create_users_table.sql (10.42ms)
2025/08/01 15:47:01 OK   000002_create_projects_table.sql (12.01ms)
2025/08/01 15:47:01 OK   000003_create_log_entries_table.sql (19.5ms)
2025/08/01 15:47:01 OK   000004_create_tags_system.sql (15.9ms)
2025/08/01 15:47:01 OK   000005_create_auth_tables.sql (20.08ms)
2025/08/01 15:47:01 OK   000006_create_insights_table.sql (23.99ms)
2025/08/01 15:47:01 OK   000007_create_performance_indexes.sql (15.41ms)
2025/08/01 15:47:01 OK   000008_create_analytics_views.sql (28.08ms)
2025/08/01 15:47:01 OK   000009_development_data.sql (27.51ms)
2025/08/01 15:47:01 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36433/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36433/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:01 🐳 Stopping container: abb53594b12c
2025/08/01 15:47:01 ✅ Container stopped: abb53594b12c
2025/08/01 15:47:01 🐳 Terminating container: abb53594b12c
2025/08/01 15:47:02 🚫 Container terminated: abb53594b12c
--- PASS: TestUserHandler_ComprehensiveScenarios (15.01s)
    --- PASS: TestUserHandler_ComprehensiveScenarios/complete_user_profile_lifecycle (3.26s)
    --- PASS: TestUserHandler_ComprehensiveScenarios/profile_update_validation (2.68s)
    --- PASS: TestUserHandler_ComprehensiveScenarios/password_change_security (3.23s)
    --- PASS: TestUserHandler_ComprehensiveScenarios/concurrent_operations (2.77s)
    --- PASS: TestUserHandler_ComprehensiveScenarios/malicious_input_handling (3.07s)
=== RUN   TestUserHandler_EdgeCases
=== RUN   TestUserHandler_EdgeCases/token_expiration_scenarios
2025/08/01 15:47:02 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:02 ✅ Container created: 80086dad8172
2025/08/01 15:47:02 🐳 Starting container: 80086dad8172
2025/08/01 15:47:02 ✅ Container started: 80086dad8172
2025/08/01 15:47:02 ⏳ Waiting for container id 80086dad8172 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc002247130 Strategies:[0xc0006570e0 0xc009bb5560]}
2025/08/01 15:47:04 🔔 Container is ready: 80086dad8172
DSN: postgres://testuser:testpass@localhost:36434/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36434/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:04 OK   000001_create_users_table.sql (10.41ms)
2025/08/01 15:47:04 OK   000002_create_projects_table.sql (12.11ms)
2025/08/01 15:47:04 OK   000003_create_log_entries_table.sql (18.65ms)
2025/08/01 15:47:04 OK   000004_create_tags_system.sql (16.02ms)
2025/08/01 15:47:04 OK   000005_create_auth_tables.sql (20.13ms)
2025/08/01 15:47:04 OK   000006_create_insights_table.sql (23.95ms)
2025/08/01 15:47:04 OK   000007_create_performance_indexes.sql (15.01ms)
2025/08/01 15:47:04 OK   000008_create_analytics_views.sql (28.16ms)
2025/08/01 15:47:04 OK   000009_development_data.sql (29.33ms)
2025/08/01 15:47:04 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36434/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36434/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:04 🐳 Stopping container: 80086dad8172
2025/08/01 15:47:04 ✅ Container stopped: 80086dad8172
2025/08/01 15:47:04 🐳 Terminating container: 80086dad8172
2025/08/01 15:47:04 🚫 Container terminated: 80086dad8172
=== RUN   TestUserHandler_EdgeCases/unicode_and_special_characters
2025/08/01 15:47:04 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:04 ✅ Container created: 73877cbfea69
2025/08/01 15:47:04 🐳 Starting container: 73877cbfea69
2025/08/01 15:47:05 ✅ Container started: 73877cbfea69
2025/08/01 15:47:05 ⏳ Waiting for container id 73877cbfea69 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00ca43770 Strategies:[0xc00ac44fc0 0xc006bf2e10]}
2025/08/01 15:47:06 🔔 Container is ready: 73877cbfea69
DSN: postgres://testuser:testpass@localhost:36435/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36435/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:06 OK   000001_create_users_table.sql (10.63ms)
2025/08/01 15:47:06 OK   000002_create_projects_table.sql (12.32ms)
2025/08/01 15:47:06 OK   000003_create_log_entries_table.sql (19.68ms)
2025/08/01 15:47:06 OK   000004_create_tags_system.sql (16.59ms)
2025/08/01 15:47:06 OK   000005_create_auth_tables.sql (20.43ms)
2025/08/01 15:47:06 OK   000006_create_insights_table.sql (24.23ms)
2025/08/01 15:47:06 OK   000007_create_performance_indexes.sql (15.27ms)
2025/08/01 15:47:06 OK   000008_create_analytics_views.sql (27.81ms)
2025/08/01 15:47:07 OK   000009_development_data.sql (28.41ms)
2025/08/01 15:47:07 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36435/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36435/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:07 🐳 Stopping container: 73877cbfea69
2025/08/01 15:47:07 ✅ Container stopped: 73877cbfea69
2025/08/01 15:47:07 🐳 Terminating container: 73877cbfea69
2025/08/01 15:47:07 🚫 Container terminated: 73877cbfea69
=== RUN   TestUserHandler_EdgeCases/large_request_payloads
2025/08/01 15:47:07 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:07 ✅ Container created: c12560444922
2025/08/01 15:47:07 🐳 Starting container: c12560444922
2025/08/01 15:47:07 ✅ Container started: c12560444922
2025/08/01 15:47:07 ⏳ Waiting for container id c12560444922 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00135ab18 Strategies:[0xc006d05020 0xc009b0ee10]}
2025/08/01 15:47:09 🔔 Container is ready: c12560444922
DSN: postgres://testuser:testpass@localhost:36436/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36436/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:09 OK   000001_create_users_table.sql (10.36ms)
2025/08/01 15:47:09 OK   000002_create_projects_table.sql (12.02ms)
2025/08/01 15:47:09 OK   000003_create_log_entries_table.sql (19.04ms)
2025/08/01 15:47:09 OK   000004_create_tags_system.sql (16.32ms)
2025/08/01 15:47:09 OK   000005_create_auth_tables.sql (20.55ms)
2025/08/01 15:47:09 OK   000006_create_insights_table.sql (24.66ms)
2025/08/01 15:47:09 OK   000007_create_performance_indexes.sql (15.16ms)
2025/08/01 15:47:09 OK   000008_create_analytics_views.sql (27.86ms)
2025/08/01 15:47:09 OK   000009_development_data.sql (28.77ms)
2025/08/01 15:47:09 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36436/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36436/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:10 🐳 Stopping container: c12560444922
2025/08/01 15:47:10 ✅ Container stopped: c12560444922
2025/08/01 15:47:10 🐳 Terminating container: c12560444922
2025/08/01 15:47:10 🚫 Container terminated: c12560444922
=== RUN   TestUserHandler_EdgeCases/null_and_empty_values
2025/08/01 15:47:10 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:10 ✅ Container created: 6658b513d1fb
2025/08/01 15:47:10 🐳 Starting container: 6658b513d1fb
2025/08/01 15:47:10 ✅ Container started: 6658b513d1fb
2025/08/01 15:47:10 ⏳ Waiting for container id 6658b513d1fb image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00aba3d48 Strategies:[0xc00a14de00 0xc009c6f050]}
2025/08/01 15:47:12 🔔 Container is ready: 6658b513d1fb
DSN: postgres://testuser:testpass@localhost:36437/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36437/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:12 OK   000001_create_users_table.sql (10.78ms)
2025/08/01 15:47:12 OK   000002_create_projects_table.sql (12.43ms)
2025/08/01 15:47:12 OK   000003_create_log_entries_table.sql (18.48ms)
2025/08/01 15:47:12 OK   000004_create_tags_system.sql (16.21ms)
2025/08/01 15:47:12 OK   000005_create_auth_tables.sql (19.94ms)
2025/08/01 15:47:12 OK   000006_create_insights_table.sql (24.23ms)
2025/08/01 15:47:12 OK   000007_create_performance_indexes.sql (15.66ms)
2025/08/01 15:47:12 OK   000008_create_analytics_views.sql (27.66ms)
2025/08/01 15:47:12 OK   000009_development_data.sql (28.53ms)
2025/08/01 15:47:12 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36437/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36437/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:12 🐳 Stopping container: 6658b513d1fb
2025/08/01 15:47:13 ✅ Container stopped: 6658b513d1fb
2025/08/01 15:47:13 🐳 Terminating container: 6658b513d1fb
2025/08/01 15:47:13 🚫 Container terminated: 6658b513d1fb
--- PASS: TestUserHandler_EdgeCases (11.14s)
    --- PASS: TestUserHandler_EdgeCases/token_expiration_scenarios (2.82s)
    --- PASS: TestUserHandler_EdgeCases/unicode_and_special_characters (2.71s)
    --- PASS: TestUserHandler_EdgeCases/large_request_payloads (2.77s)
    --- PASS: TestUserHandler_EdgeCases/null_and_empty_values (2.84s)
=== RUN   TestUserHandler_PerformanceScenarios
=== RUN   TestUserHandler_PerformanceScenarios/rapid_successive_requests
2025/08/01 15:47:13 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:13 ✅ Container created: 2b797c66d3a0
2025/08/01 15:47:13 🐳 Starting container: 2b797c66d3a0
2025/08/01 15:47:13 ✅ Container started: 2b797c66d3a0
2025/08/01 15:47:13 ⏳ Waiting for container id 2b797c66d3a0 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00ca83ea0 Strategies:[0xc006a5f500 0xc00693ca80]}
2025/08/01 15:47:15 🔔 Container is ready: 2b797c66d3a0
DSN: postgres://testuser:testpass@localhost:36438/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36438/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:15 OK   000001_create_users_table.sql (10.48ms)
2025/08/01 15:47:15 OK   000002_create_projects_table.sql (11.89ms)
2025/08/01 15:47:15 OK   000003_create_log_entries_table.sql (19.47ms)
2025/08/01 15:47:15 OK   000004_create_tags_system.sql (16.32ms)
2025/08/01 15:47:15 OK   000005_create_auth_tables.sql (20.22ms)
2025/08/01 15:47:15 OK   000006_create_insights_table.sql (24.02ms)
2025/08/01 15:47:15 OK   000007_create_performance_indexes.sql (15.18ms)
2025/08/01 15:47:15 OK   000008_create_analytics_views.sql (27.92ms)
2025/08/01 15:47:15 OK   000009_development_data.sql (28.58ms)
2025/08/01 15:47:15 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36438/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36438/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:15 🐳 Stopping container: 2b797c66d3a0
2025/08/01 15:47:15 ✅ Container stopped: 2b797c66d3a0
2025/08/01 15:47:15 🐳 Terminating container: 2b797c66d3a0
2025/08/01 15:47:15 🚫 Container terminated: 2b797c66d3a0
=== RUN   TestUserHandler_PerformanceScenarios/memory_efficiency
2025/08/01 15:47:15 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:15 ✅ Container created: 33d0150a0069
2025/08/01 15:47:15 🐳 Starting container: 33d0150a0069
2025/08/01 15:47:16 ✅ Container started: 33d0150a0069
2025/08/01 15:47:16 ⏳ Waiting for container id 33d0150a0069 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00aaf5480 Strategies:[0xc00b191020 0xc003c63560]}
2025/08/01 15:47:17 🔔 Container is ready: 33d0150a0069
DSN: postgres://testuser:testpass@localhost:36439/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36439/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:17 OK   000001_create_users_table.sql (10.66ms)
2025/08/01 15:47:17 OK   000002_create_projects_table.sql (12.56ms)
2025/08/01 15:47:17 OK   000003_create_log_entries_table.sql (19.29ms)
2025/08/01 15:47:17 OK   000004_create_tags_system.sql (16.6ms)
2025/08/01 15:47:17 OK   000005_create_auth_tables.sql (20.76ms)
2025/08/01 15:47:17 OK   000006_create_insights_table.sql (25.3ms)
2025/08/01 15:47:17 OK   000007_create_performance_indexes.sql (15.76ms)
2025/08/01 15:47:18 OK   000008_create_analytics_views.sql (28.68ms)
2025/08/01 15:47:18 OK   000009_development_data.sql (28.41ms)
2025/08/01 15:47:18 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36439/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36439/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:18 🐳 Stopping container: 33d0150a0069
2025/08/01 15:47:18 ✅ Container stopped: 33d0150a0069
2025/08/01 15:47:18 🐳 Terminating container: 33d0150a0069
2025/08/01 15:47:18 🚫 Container terminated: 33d0150a0069
--- PASS: TestUserHandler_PerformanceScenarios (5.56s)
    --- PASS: TestUserHandler_PerformanceScenarios/rapid_successive_requests (2.64s)
    --- PASS: TestUserHandler_PerformanceScenarios/memory_efficiency (2.92s)
=== RUN   TestUserHandler_GetProfile
=== RUN   TestUserHandler_GetProfile/successful_profile_retrieval
2025/08/01 15:47:18 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:18 ✅ Container created: 8e408f4f2a02
2025/08/01 15:47:18 🐳 Starting container: 8e408f4f2a02
2025/08/01 15:47:18 ✅ Container started: 8e408f4f2a02
2025/08/01 15:47:18 ⏳ Waiting for container id 8e408f4f2a02 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00ab8ae50 Strategies:[0xc009bda5a0 0xc001b57980]}
2025/08/01 15:47:20 🔔 Container is ready: 8e408f4f2a02
DSN: postgres://testuser:testpass@localhost:36440/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36440/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:20 OK   000001_create_users_table.sql (10.8ms)
2025/08/01 15:47:20 OK   000002_create_projects_table.sql (12.62ms)
2025/08/01 15:47:20 OK   000003_create_log_entries_table.sql (19.88ms)
2025/08/01 15:47:20 OK   000004_create_tags_system.sql (17.06ms)
2025/08/01 15:47:20 OK   000005_create_auth_tables.sql (20.41ms)
2025/08/01 15:47:20 OK   000006_create_insights_table.sql (24.93ms)
2025/08/01 15:47:20 OK   000007_create_performance_indexes.sql (16.11ms)
2025/08/01 15:47:20 OK   000008_create_analytics_views.sql (28.87ms)
2025/08/01 15:47:20 OK   000009_development_data.sql (29.47ms)
2025/08/01 15:47:20 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36440/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36440/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:21 🐳 Stopping container: 8e408f4f2a02
2025/08/01 15:47:21 ✅ Container stopped: 8e408f4f2a02
2025/08/01 15:47:21 🐳 Terminating container: 8e408f4f2a02
2025/08/01 15:47:21 🚫 Container terminated: 8e408f4f2a02
=== RUN   TestUserHandler_GetProfile/unauthorized_access
2025/08/01 15:47:21 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:21 ✅ Container created: 2281f43328b3
2025/08/01 15:47:21 🐳 Starting container: 2281f43328b3
2025/08/01 15:47:21 ✅ Container started: 2281f43328b3
2025/08/01 15:47:21 ⏳ Waiting for container id 2281f43328b3 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007f05dc0 Strategies:[0xc0006565a0 0xc0099b7ec0]}
2025/08/01 15:47:23 🔔 Container is ready: 2281f43328b3
DSN: postgres://testuser:testpass@localhost:36441/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36441/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:23 OK   000001_create_users_table.sql (10.45ms)
2025/08/01 15:47:23 OK   000002_create_projects_table.sql (12.19ms)
2025/08/01 15:47:23 OK   000003_create_log_entries_table.sql (18.51ms)
2025/08/01 15:47:23 OK   000004_create_tags_system.sql (16ms)
2025/08/01 15:47:23 OK   000005_create_auth_tables.sql (19.56ms)
2025/08/01 15:47:23 OK   000006_create_insights_table.sql (23.8ms)
2025/08/01 15:47:23 OK   000007_create_performance_indexes.sql (15.12ms)
2025/08/01 15:47:23 OK   000008_create_analytics_views.sql (27.75ms)
2025/08/01 15:47:23 OK   000009_development_data.sql (28.4ms)
2025/08/01 15:47:23 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36441/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36441/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:23 🐳 Stopping container: 2281f43328b3
2025/08/01 15:47:24 ✅ Container stopped: 2281f43328b3
2025/08/01 15:47:24 🐳 Terminating container: 2281f43328b3
2025/08/01 15:47:24 🚫 Container terminated: 2281f43328b3
--- PASS: TestUserHandler_GetProfile (5.36s)
    --- PASS: TestUserHandler_GetProfile/successful_profile_retrieval (2.82s)
    --- PASS: TestUserHandler_GetProfile/unauthorized_access (2.54s)
=== RUN   TestUserHandler_UpdateProfile
=== RUN   TestUserHandler_UpdateProfile/successful_profile_update
2025/08/01 15:47:24 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:24 ✅ Container created: 1f90096ddce7
2025/08/01 15:47:24 🐳 Starting container: 1f90096ddce7
2025/08/01 15:47:24 ✅ Container started: 1f90096ddce7
2025/08/01 15:47:24 ⏳ Waiting for container id 1f90096ddce7 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc007bf01d0 Strategies:[0xc001201c80 0xc007bf8540]}
2025/08/01 15:47:26 🔔 Container is ready: 1f90096ddce7
DSN: postgres://testuser:testpass@localhost:36442/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36442/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:26 OK   000001_create_users_table.sql (9.19ms)
2025/08/01 15:47:26 OK   000002_create_projects_table.sql (10.82ms)
2025/08/01 15:47:26 OK   000003_create_log_entries_table.sql (16.78ms)
2025/08/01 15:47:26 OK   000004_create_tags_system.sql (13.54ms)
2025/08/01 15:47:26 OK   000005_create_auth_tables.sql (16.12ms)
2025/08/01 15:47:26 OK   000006_create_insights_table.sql (19.71ms)
2025/08/01 15:47:26 OK   000007_create_performance_indexes.sql (12.47ms)
2025/08/01 15:47:26 OK   000008_create_analytics_views.sql (22.85ms)
2025/08/01 15:47:26 OK   000009_development_data.sql (23.6ms)
2025/08/01 15:47:26 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36442/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36442/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:26 🐳 Stopping container: 1f90096ddce7
2025/08/01 15:47:26 ✅ Container stopped: 1f90096ddce7
2025/08/01 15:47:26 🐳 Terminating container: 1f90096ddce7
2025/08/01 15:47:26 🚫 Container terminated: 1f90096ddce7
=== RUN   TestUserHandler_UpdateProfile/invalid_request_body
2025/08/01 15:47:26 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:26 ✅ Container created: 4be77f49d27a
2025/08/01 15:47:26 🐳 Starting container: 4be77f49d27a
2025/08/01 15:47:27 ✅ Container started: 4be77f49d27a
2025/08/01 15:47:27 ⏳ Waiting for container id 4be77f49d27a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc004256af0 Strategies:[0xc004730c00 0xc0099a81e0]}
2025/08/01 15:47:29 🔔 Container is ready: 4be77f49d27a
DSN: postgres://testuser:testpass@localhost:36443/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36443/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:29 OK   000001_create_users_table.sql (11.18ms)
2025/08/01 15:47:29 OK   000002_create_projects_table.sql (12.92ms)
2025/08/01 15:47:29 OK   000003_create_log_entries_table.sql (20.29ms)
2025/08/01 15:47:29 OK   000004_create_tags_system.sql (16.8ms)
2025/08/01 15:47:29 OK   000005_create_auth_tables.sql (20.82ms)
2025/08/01 15:47:29 OK   000006_create_insights_table.sql (24.81ms)
2025/08/01 15:47:29 OK   000007_create_performance_indexes.sql (16.1ms)
2025/08/01 15:47:29 OK   000008_create_analytics_views.sql (28.77ms)
2025/08/01 15:47:29 OK   000009_development_data.sql (28.51ms)
2025/08/01 15:47:29 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36443/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36443/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:29 🐳 Stopping container: 4be77f49d27a
2025/08/01 15:47:29 ✅ Container stopped: 4be77f49d27a
2025/08/01 15:47:29 🐳 Terminating container: 4be77f49d27a
2025/08/01 15:47:29 🚫 Container terminated: 4be77f49d27a
=== RUN   TestUserHandler_UpdateProfile/unauthorized_update
2025/08/01 15:47:29 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:29 ✅ Container created: e98c50339225
2025/08/01 15:47:29 🐳 Starting container: e98c50339225
2025/08/01 15:47:30 ✅ Container started: e98c50339225
2025/08/01 15:47:30 ⏳ Waiting for container id e98c50339225 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00bea50d8 Strategies:[0xc00be94de0 0xc00bec69f0]}
2025/08/01 15:47:31 🔔 Container is ready: e98c50339225
DSN: postgres://testuser:testpass@localhost:36444/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36444/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:32 OK   000001_create_users_table.sql (11.82ms)
2025/08/01 15:47:32 OK   000002_create_projects_table.sql (17.32ms)
2025/08/01 15:47:32 OK   000003_create_log_entries_table.sql (30.24ms)
2025/08/01 15:47:32 OK   000004_create_tags_system.sql (25.71ms)
2025/08/01 15:47:32 OK   000005_create_auth_tables.sql (27.08ms)
2025/08/01 15:47:32 OK   000006_create_insights_table.sql (20.6ms)
2025/08/01 15:47:32 OK   000007_create_performance_indexes.sql (12.85ms)
2025/08/01 15:47:32 OK   000008_create_analytics_views.sql (29.11ms)
2025/08/01 15:47:32 OK   000009_development_data.sql (28.03ms)
2025/08/01 15:47:32 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36444/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36444/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:32 🐳 Stopping container: e98c50339225
2025/08/01 15:47:32 ✅ Container stopped: e98c50339225
2025/08/01 15:47:32 🐳 Terminating container: e98c50339225
2025/08/01 15:47:32 🚫 Container terminated: e98c50339225
--- PASS: TestUserHandler_UpdateProfile (8.43s)
    --- PASS: TestUserHandler_UpdateProfile/successful_profile_update (2.81s)
    --- PASS: TestUserHandler_UpdateProfile/invalid_request_body (2.93s)
    --- PASS: TestUserHandler_UpdateProfile/unauthorized_update (2.69s)
=== RUN   TestUserHandler_ChangePassword
=== RUN   TestUserHandler_ChangePassword/successful_password_change
2025/08/01 15:47:32 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:32 ✅ Container created: 16aac8b25d56
2025/08/01 15:47:32 🐳 Starting container: 16aac8b25d56
2025/08/01 15:47:32 ✅ Container started: 16aac8b25d56
2025/08/01 15:47:32 ⏳ Waiting for container id 16aac8b25d56 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00c7d3770 Strategies:[0xc00c98aa20 0xc00c989770]}
2025/08/01 15:47:34 🔔 Container is ready: 16aac8b25d56
DSN: postgres://testuser:testpass@localhost:36445/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36445/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:34 OK   000001_create_users_table.sql (11.51ms)
2025/08/01 15:47:34 OK   000002_create_projects_table.sql (12.93ms)
2025/08/01 15:47:34 OK   000003_create_log_entries_table.sql (20.23ms)
2025/08/01 15:47:34 OK   000004_create_tags_system.sql (18.06ms)
2025/08/01 15:47:34 OK   000005_create_auth_tables.sql (21.75ms)
2025/08/01 15:47:34 OK   000006_create_insights_table.sql (25.52ms)
2025/08/01 15:47:34 OK   000007_create_performance_indexes.sql (15.88ms)
2025/08/01 15:47:34 OK   000008_create_analytics_views.sql (28.28ms)
2025/08/01 15:47:34 OK   000009_development_data.sql (29.54ms)
2025/08/01 15:47:34 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36445/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36445/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:35 🐳 Stopping container: 16aac8b25d56
2025/08/01 15:47:35 ✅ Container stopped: 16aac8b25d56
2025/08/01 15:47:35 🐳 Terminating container: 16aac8b25d56
2025/08/01 15:47:35 🚫 Container terminated: 16aac8b25d56
=== RUN   TestUserHandler_ChangePassword/invalid_current_password
2025/08/01 15:47:35 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:35 ✅ Container created: 6edd1ffc4085
2025/08/01 15:47:35 🐳 Starting container: 6edd1ffc4085
2025/08/01 15:47:35 ✅ Container started: 6edd1ffc4085
2025/08/01 15:47:35 ⏳ Waiting for container id 6edd1ffc4085 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00ca43af0 Strategies:[0xc00cea8420 0xc009e0f560]}
2025/08/01 15:47:37 🔔 Container is ready: 6edd1ffc4085
DSN: postgres://testuser:testpass@localhost:36446/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36446/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:37 OK   000001_create_users_table.sql (11.06ms)
2025/08/01 15:47:37 OK   000002_create_projects_table.sql (13.2ms)
2025/08/01 15:47:37 OK   000003_create_log_entries_table.sql (19.84ms)
2025/08/01 15:47:37 OK   000004_create_tags_system.sql (17.17ms)
2025/08/01 15:47:37 OK   000005_create_auth_tables.sql (21.56ms)
2025/08/01 15:47:37 OK   000006_create_insights_table.sql (25.28ms)
2025/08/01 15:47:37 OK   000007_create_performance_indexes.sql (15.5ms)
2025/08/01 15:47:37 OK   000008_create_analytics_views.sql (28.91ms)
2025/08/01 15:47:37 OK   000009_development_data.sql (28.4ms)
2025/08/01 15:47:37 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36446/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36446/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:38 🐳 Stopping container: 6edd1ffc4085
2025/08/01 15:47:38 ✅ Container stopped: 6edd1ffc4085
2025/08/01 15:47:38 🐳 Terminating container: 6edd1ffc4085
2025/08/01 15:47:38 🚫 Container terminated: 6edd1ffc4085
=== RUN   TestUserHandler_ChangePassword/invalid_request_body
2025/08/01 15:47:38 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:38 ✅ Container created: e0c6e6f49db1
2025/08/01 15:47:38 🐳 Starting container: e0c6e6f49db1
2025/08/01 15:47:38 ✅ Container started: e0c6e6f49db1
2025/08/01 15:47:38 ⏳ Waiting for container id e0c6e6f49db1 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009eb2a60 Strategies:[0xc006c12ba0 0xc009c10de0]}
2025/08/01 15:47:40 🔔 Container is ready: e0c6e6f49db1
DSN: postgres://testuser:testpass@localhost:36447/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36447/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:40 OK   000001_create_users_table.sql (10.81ms)
2025/08/01 15:47:40 OK   000002_create_projects_table.sql (12.37ms)
2025/08/01 15:47:40 OK   000003_create_log_entries_table.sql (19.66ms)
2025/08/01 15:47:40 OK   000004_create_tags_system.sql (16.4ms)
2025/08/01 15:47:40 OK   000005_create_auth_tables.sql (19.84ms)
2025/08/01 15:47:40 OK   000006_create_insights_table.sql (24.17ms)
2025/08/01 15:47:40 OK   000007_create_performance_indexes.sql (15.52ms)
2025/08/01 15:47:40 OK   000008_create_analytics_views.sql (27.94ms)
2025/08/01 15:47:40 OK   000009_development_data.sql (28.66ms)
2025/08/01 15:47:40 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36447/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36447/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:41 🐳 Stopping container: e0c6e6f49db1
2025/08/01 15:47:41 ✅ Container stopped: e0c6e6f49db1
2025/08/01 15:47:41 🐳 Terminating container: e0c6e6f49db1
2025/08/01 15:47:41 🚫 Container terminated: e0c6e6f49db1
=== RUN   TestUserHandler_ChangePassword/unauthorized_password_change
2025/08/01 15:47:41 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:41 ✅ Container created: bbeb7c5bb331
2025/08/01 15:47:41 🐳 Starting container: bbeb7c5bb331
2025/08/01 15:47:41 ✅ Container started: bbeb7c5bb331
2025/08/01 15:47:41 ⏳ Waiting for container id bbeb7c5bb331 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc00590fa68 Strategies:[0xc009fbeae0 0xc009eab8c0]}
2025/08/01 15:47:43 🔔 Container is ready: bbeb7c5bb331
DSN: postgres://testuser:testpass@localhost:36448/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36448/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:43 OK   000001_create_users_table.sql (9.81ms)
2025/08/01 15:47:43 OK   000002_create_projects_table.sql (12.28ms)
2025/08/01 15:47:43 OK   000003_create_log_entries_table.sql (19.93ms)
2025/08/01 15:47:43 OK   000004_create_tags_system.sql (16.77ms)
2025/08/01 15:47:43 OK   000005_create_auth_tables.sql (21.12ms)
2025/08/01 15:47:43 OK   000006_create_insights_table.sql (25.38ms)
2025/08/01 15:47:43 OK   000007_create_performance_indexes.sql (16.1ms)
2025/08/01 15:47:43 OK   000008_create_analytics_views.sql (29.14ms)
2025/08/01 15:47:43 OK   000009_development_data.sql (29.98ms)
2025/08/01 15:47:43 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36448/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36448/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:43 🐳 Stopping container: bbeb7c5bb331
2025/08/01 15:47:43 ✅ Container stopped: bbeb7c5bb331
2025/08/01 15:47:43 🐳 Terminating container: bbeb7c5bb331
2025/08/01 15:47:44 🚫 Container terminated: bbeb7c5bb331
--- PASS: TestUserHandler_ChangePassword (11.49s)
    --- PASS: TestUserHandler_ChangePassword/successful_password_change (3.18s)
    --- PASS: TestUserHandler_ChangePassword/invalid_current_password (2.92s)
    --- PASS: TestUserHandler_ChangePassword/invalid_request_body (2.84s)
    --- PASS: TestUserHandler_ChangePassword/unauthorized_password_change (2.54s)
=== RUN   TestUserHandler_DeleteAccount
=== RUN   TestUserHandler_DeleteAccount/successful_account_deletion
2025/08/01 15:47:44 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:44 ✅ Container created: 78fed1a9203a
2025/08/01 15:47:44 🐳 Starting container: 78fed1a9203a
2025/08/01 15:47:44 ✅ Container started: 78fed1a9203a
2025/08/01 15:47:44 ⏳ Waiting for container id 78fed1a9203a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc001e74a00 Strategies:[0xc0069bb020 0xc006c31e30]}
2025/08/01 15:47:46 🔔 Container is ready: 78fed1a9203a
DSN: postgres://testuser:testpass@localhost:36449/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36449/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:46 OK   000001_create_users_table.sql (10.42ms)
2025/08/01 15:47:46 OK   000002_create_projects_table.sql (12.78ms)
2025/08/01 15:47:46 OK   000003_create_log_entries_table.sql (19.56ms)
2025/08/01 15:47:46 OK   000004_create_tags_system.sql (16.53ms)
2025/08/01 15:47:46 OK   000005_create_auth_tables.sql (20.27ms)
2025/08/01 15:47:46 OK   000006_create_insights_table.sql (24.87ms)
2025/08/01 15:47:46 OK   000007_create_performance_indexes.sql (15.89ms)
2025/08/01 15:47:46 OK   000008_create_analytics_views.sql (27.96ms)
2025/08/01 15:47:46 OK   000009_development_data.sql (28.7ms)
2025/08/01 15:47:46 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36449/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36449/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:46 🐳 Stopping container: 78fed1a9203a
2025/08/01 15:47:46 ✅ Container stopped: 78fed1a9203a
2025/08/01 15:47:46 🐳 Terminating container: 78fed1a9203a
2025/08/01 15:47:46 🚫 Container terminated: 78fed1a9203a
=== RUN   TestUserHandler_DeleteAccount/unauthorized_account_deletion
2025/08/01 15:47:46 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:46 ✅ Container created: 857574cabed3
2025/08/01 15:47:46 🐳 Starting container: 857574cabed3
2025/08/01 15:47:47 ✅ Container started: 857574cabed3
2025/08/01 15:47:47 ⏳ Waiting for container id 857574cabed3 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc006a055b0 Strategies:[0xc00ca81bc0 0xc006a71320]}
2025/08/01 15:47:48 🔔 Container is ready: 857574cabed3
DSN: postgres://testuser:testpass@localhost:36450/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36450/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:48 OK   000001_create_users_table.sql (10.94ms)
2025/08/01 15:47:48 OK   000002_create_projects_table.sql (12.6ms)
2025/08/01 15:47:48 OK   000003_create_log_entries_table.sql (19.52ms)
2025/08/01 15:47:48 OK   000004_create_tags_system.sql (16.85ms)
2025/08/01 15:47:48 OK   000005_create_auth_tables.sql (20.46ms)
2025/08/01 15:47:48 OK   000006_create_insights_table.sql (24.69ms)
2025/08/01 15:47:48 OK   000007_create_performance_indexes.sql (15.45ms)
2025/08/01 15:47:49 OK   000008_create_analytics_views.sql (28.12ms)
2025/08/01 15:47:49 OK   000009_development_data.sql (29.06ms)
2025/08/01 15:47:49 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36450/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36450/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:49 🐳 Stopping container: 857574cabed3
2025/08/01 15:47:49 ✅ Container stopped: 857574cabed3
2025/08/01 15:47:49 🐳 Terminating container: 857574cabed3
2025/08/01 15:47:49 🚫 Container terminated: 857574cabed3
--- PASS: TestUserHandler_DeleteAccount (5.37s)
    --- PASS: TestUserHandler_DeleteAccount/successful_account_deletion (2.78s)
    --- PASS: TestUserHandler_DeleteAccount/unauthorized_account_deletion (2.60s)
=== RUN   TestUserHandler_ErrorHandling
=== RUN   TestUserHandler_ErrorHandling/profile_retrieval_with_invalid_token
2025/08/01 15:47:49 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:49 ✅ Container created: e79073b3ce5a
2025/08/01 15:47:49 🐳 Starting container: e79073b3ce5a
2025/08/01 15:47:49 ✅ Container started: e79073b3ce5a
2025/08/01 15:47:49 ⏳ Waiting for container id e79073b3ce5a image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc0071ece40 Strategies:[0xc00ca44480 0xc006a5d1a0]}
2025/08/01 15:47:51 🔔 Container is ready: e79073b3ce5a
DSN: postgres://testuser:testpass@localhost:36451/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36451/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:51 OK   000001_create_users_table.sql (10.82ms)
2025/08/01 15:47:51 OK   000002_create_projects_table.sql (12.68ms)
2025/08/01 15:47:51 OK   000003_create_log_entries_table.sql (19.45ms)
2025/08/01 15:47:51 OK   000004_create_tags_system.sql (16.82ms)
2025/08/01 15:47:51 OK   000005_create_auth_tables.sql (21.35ms)
2025/08/01 15:47:51 OK   000006_create_insights_table.sql (25.08ms)
2025/08/01 15:47:51 OK   000007_create_performance_indexes.sql (16.15ms)
2025/08/01 15:47:51 OK   000008_create_analytics_views.sql (29.01ms)
2025/08/01 15:47:51 OK   000009_development_data.sql (29.45ms)
2025/08/01 15:47:51 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36451/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36451/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:51 🐳 Stopping container: e79073b3ce5a
2025/08/01 15:47:51 ✅ Container stopped: e79073b3ce5a
2025/08/01 15:47:51 🐳 Terminating container: e79073b3ce5a
2025/08/01 15:47:51 🚫 Container terminated: e79073b3ce5a
=== RUN   TestUserHandler_ErrorHandling/profile_update_with_empty_body
2025/08/01 15:47:51 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:52 ✅ Container created: 858f8d5b5ac4
2025/08/01 15:47:52 🐳 Starting container: 858f8d5b5ac4
2025/08/01 15:47:52 ✅ Container started: 858f8d5b5ac4
2025/08/01 15:47:52 ⏳ Waiting for container id 858f8d5b5ac4 image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009a5ee10 Strategies:[0xc0039481e0 0xc006c50ed0]}
2025/08/01 15:47:53 🔔 Container is ready: 858f8d5b5ac4
DSN: postgres://testuser:testpass@localhost:36453/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36453/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:54 OK   000001_create_users_table.sql (9.95ms)
2025/08/01 15:47:54 OK   000002_create_projects_table.sql (11.89ms)
2025/08/01 15:47:54 OK   000003_create_log_entries_table.sql (19.17ms)
2025/08/01 15:47:54 OK   000004_create_tags_system.sql (16.07ms)
2025/08/01 15:47:54 OK   000005_create_auth_tables.sql (21.03ms)
2025/08/01 15:47:54 OK   000006_create_insights_table.sql (25.57ms)
2025/08/01 15:47:54 OK   000007_create_performance_indexes.sql (15.52ms)
2025/08/01 15:47:54 OK   000008_create_analytics_views.sql (28.75ms)
2025/08/01 15:47:54 OK   000009_development_data.sql (28.8ms)
2025/08/01 15:47:54 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36453/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36453/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:54 🐳 Stopping container: 858f8d5b5ac4
2025/08/01 15:47:54 ✅ Container stopped: 858f8d5b5ac4
2025/08/01 15:47:54 🐳 Terminating container: 858f8d5b5ac4
2025/08/01 15:47:54 🚫 Container terminated: 858f8d5b5ac4
=== RUN   TestUserHandler_ErrorHandling/password_change_with_missing_fields
2025/08/01 15:47:54 🐳 Creating container for image postgres:17-alpine
2025/08/01 15:47:54 ✅ Container created: 351acecbf79e
2025/08/01 15:47:54 🐳 Starting container: 351acecbf79e
2025/08/01 15:47:55 ✅ Container started: 351acecbf79e
2025/08/01 15:47:55 ⏳ Waiting for container id 351acecbf79e image: postgres:17-alpine. Waiting for: &{timeout:<nil> deadline:0xc009c99ed8 Strategies:[0xc007cbb020 0xc007d0a0f0]}
2025/08/01 15:47:56 🔔 Container is ready: 351acecbf79e
DSN: postgres://testuser:testpass@localhost:36454/englog_test?application_name=englog&connect_timeout=10&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36454/englog_test?application_name=englog&search_path=englog_test,public&connect_timeout=10&sslmode=disable
2025/08/01 15:47:56 OK   000001_create_users_table.sql (11.01ms)
2025/08/01 15:47:56 OK   000002_create_projects_table.sql (13.05ms)
2025/08/01 15:47:57 OK   000003_create_log_entries_table.sql (19.77ms)
2025/08/01 15:47:57 OK   000004_create_tags_system.sql (17.17ms)
2025/08/01 15:47:57 OK   000005_create_auth_tables.sql (21.27ms)
2025/08/01 15:47:57 OK   000006_create_insights_table.sql (25.94ms)
2025/08/01 15:47:57 OK   000007_create_performance_indexes.sql (17.41ms)
2025/08/01 15:47:57 OK   000008_create_analytics_views.sql (29.73ms)
2025/08/01 15:47:57 OK   000009_development_data.sql (29.34ms)
2025/08/01 15:47:57 goose: successfully migrated database to version: 9
DSN: postgres://testuser:testpass@localhost:36454/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
DSN: postgres://testuser:testpass@localhost:36454/englog_test?application_name=englog&search_path=englog_test,public&sslmode=disable
2025/08/01 15:47:57 🐳 Stopping container: 351acecbf79e
2025/08/01 15:47:57 ✅ Container stopped: 351acecbf79e
2025/08/01 15:47:57 🐳 Terminating container: 351acecbf79e
2025/08/01 15:47:57 🚫 Container terminated: 351acecbf79e
--- PASS: TestUserHandler_ErrorHandling (8.31s)
    --- PASS: TestUserHandler_ErrorHandling/profile_retrieval_with_invalid_token (2.57s)
    --- PASS: TestUserHandler_ErrorHandling/profile_update_with_empty_body (2.82s)
    --- PASS: TestUserHandler_ErrorHandling/password_change_with_missing_fields (2.91s)
=== RUN   TestSetupRoutesWithRateLimit
--- PASS: TestSetupRoutesWithRateLimit (0.00s)
PASS
ok      github.com/garnizeh/englog/internal/handlers    459.879s
?       github.com/garnizeh/englog/internal/logging     [no test files]
=== RUN   TestErrorLogger
--- PASS: TestErrorLogger (0.00s)
=== RUN   TestRequestLogger
--- PASS: TestRequestLogger (0.00s)
=== RUN   TestMiddlewareIntegration
=== RUN   TestMiddlewareIntegration/multiple_middlewares_working_together
=== RUN   TestMiddlewareIntegration/validation_middleware_failure_in_middleware_stack
--- PASS: TestMiddlewareIntegration (0.00s)
    --- PASS: TestMiddlewareIntegration/multiple_middlewares_working_together (0.00s)
    --- PASS: TestMiddlewareIntegration/validation_middleware_failure_in_middleware_stack (0.00s)
=== RUN   TestMiddlewareEdgeCases
=== RUN   TestMiddlewareEdgeCases/ValidationMiddleware_with_special_characters_in_param_name
=== RUN   TestMiddlewareEdgeCases/CORS_middleware_with_complex_requests
--- PASS: TestMiddlewareEdgeCases (0.00s)
    --- PASS: TestMiddlewareEdgeCases/ValidationMiddleware_with_special_characters_in_param_name (0.00s)
    --- PASS: TestMiddlewareEdgeCases/CORS_middleware_with_complex_requests (0.00s)
=== RUN   TestRateLimitMiddleware
=== RUN   TestRateLimitMiddleware/rate_limiting_disabled
=== RUN   TestRateLimitMiddleware/rate_limiting_enabled_but_redis_disabled
--- PASS: TestRateLimitMiddleware (0.00s)
    --- PASS: TestRateLimitMiddleware/rate_limiting_disabled (0.00s)
    --- PASS: TestRateLimitMiddleware/rate_limiting_enabled_but_redis_disabled (0.00s)
=== RUN   TestRecoveryLogger


2025/08/01 15:40:17 [Recovery] 2025/08/01 - 15:40:17 panic recovered:
test panic
/media/code/code/Go/garnizeh/englog/internal/middleware/recover_test.go:32 (0x6c7084)
        TestRecoveryLogger.func1: panic("test panic")
/home/user/go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/context.go:185 (0x6a64ae)
        (*Context).Next: c.handlers[c.index](c)
/home/user/go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/recovery.go:102 (0x6a649b)
        CustomRecoveryWithWriter.func1: c.Next()
/home/user/go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/context.go:185 (0x6a5271)
        (*Context).Next: c.handlers[c.index](c)
/home/user/go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/gin.go:644 (0x6a4d00)
        (*Engine).handleHTTPRequest: c.Next()
/home/user/go/pkg/mod/github.com/gin-gonic/gin@v1.10.1/gin.go:600 (0x6a4989)
        (*Engine).ServeHTTP: engine.handleHTTPRequest(c)
/media/code/code/Go/garnizeh/englog/internal/middleware/recover_test.go:39 (0x6c14d3)
        TestRecoveryLogger: router.ServeHTTP(w, req)
/usr/local/go/src/testing/testing.go:1792 (0x510b13)
        tRunner: fn(t)
/usr/local/go/src/runtime/asm_amd64.s:1700 (0x47db20)
        goexit: BYTE    $0x90   // NOP

--- PASS: TestRecoveryLogger (0.00s)
=== RUN   TestSecurityHeaders
--- PASS: TestSecurityHeaders (0.00s)
=== RUN   TestCORS
=== RUN   TestCORS/allowed_origin
=== RUN   TestCORS/options_request
=== RUN   TestCORS/disallowed_origin
--- PASS: TestCORS (0.00s)
    --- PASS: TestCORS/allowed_origin (0.00s)
    --- PASS: TestCORS/options_request (0.00s)
    --- PASS: TestCORS/disallowed_origin (0.00s)
=== RUN   TestRequestTimeout
    timeout_test.go:20: Skipping test due to data race conditions in timeout middleware
--- SKIP: TestRequestTimeout (0.00s)
=== RUN   TestRequestTimeoutContextCancellation
    timeout_test.go:115: Skipping test due to data race conditions in timeout middleware
--- SKIP: TestRequestTimeoutContextCancellation (0.00s)
=== RUN   TestRequestTimeoutMiddlewareOrder
    timeout_test.go:150: Skipping test due to data race conditions in concurrent middleware execution
--- SKIP: TestRequestTimeoutMiddlewareOrder (0.00s)
=== RUN   TestValidationMiddleware_ValidateUUIDParam
=== RUN   TestValidationMiddleware_ValidateUUIDParam/valid_UUID_parameter
=== RUN   TestValidationMiddleware_ValidateUUIDParam/missing_parameter
=== RUN   TestValidationMiddleware_ValidateUUIDParam/valid_UUID_parameter_(different_param_name)
=== RUN   TestValidationMiddleware_ValidateUUIDParam/invalid_UUID_format
--- PASS: TestValidationMiddleware_ValidateUUIDParam (0.00s)
    --- PASS: TestValidationMiddleware_ValidateUUIDParam/valid_UUID_parameter (0.00s)
    --- PASS: TestValidationMiddleware_ValidateUUIDParam/missing_parameter (0.00s)
    --- PASS: TestValidationMiddleware_ValidateUUIDParam/valid_UUID_parameter_(different_param_name) (0.00s)
    --- PASS: TestValidationMiddleware_ValidateUUIDParam/invalid_UUID_format (0.00s)
=== RUN   TestNewValidationMiddleware
--- PASS: TestNewValidationMiddleware (0.00s)
PASS
ok      github.com/garnizeh/englog/internal/middleware  0.015s
=== RUN   TestActivityType_IsValid
=== RUN   TestActivityType_IsValid/valid_development
=== RUN   TestActivityType_IsValid/valid_meeting
=== RUN   TestActivityType_IsValid/valid_code_review
=== RUN   TestActivityType_IsValid/valid_debugging
=== RUN   TestActivityType_IsValid/valid_documentation
=== RUN   TestActivityType_IsValid/valid_testing
=== RUN   TestActivityType_IsValid/valid_deployment
=== RUN   TestActivityType_IsValid/valid_research
=== RUN   TestActivityType_IsValid/valid_planning
=== RUN   TestActivityType_IsValid/valid_learning
=== RUN   TestActivityType_IsValid/valid_maintenance
=== RUN   TestActivityType_IsValid/valid_support
=== RUN   TestActivityType_IsValid/valid_other
=== RUN   TestActivityType_IsValid/invalid_type
--- PASS: TestActivityType_IsValid (0.00s)
    --- PASS: TestActivityType_IsValid/valid_development (0.00s)
    --- PASS: TestActivityType_IsValid/valid_meeting (0.00s)
    --- PASS: TestActivityType_IsValid/valid_code_review (0.00s)
    --- PASS: TestActivityType_IsValid/valid_debugging (0.00s)
    --- PASS: TestActivityType_IsValid/valid_documentation (0.00s)
    --- PASS: TestActivityType_IsValid/valid_testing (0.00s)
    --- PASS: TestActivityType_IsValid/valid_deployment (0.00s)
    --- PASS: TestActivityType_IsValid/valid_research (0.00s)
    --- PASS: TestActivityType_IsValid/valid_planning (0.00s)
    --- PASS: TestActivityType_IsValid/valid_learning (0.00s)
    --- PASS: TestActivityType_IsValid/valid_maintenance (0.00s)
    --- PASS: TestActivityType_IsValid/valid_support (0.00s)
    --- PASS: TestActivityType_IsValid/valid_other (0.00s)
    --- PASS: TestActivityType_IsValid/invalid_type (0.00s)
=== RUN   TestValueRating_IsValid
=== RUN   TestValueRating_IsValid/valid_low
=== RUN   TestValueRating_IsValid/valid_medium
=== RUN   TestValueRating_IsValid/valid_high
=== RUN   TestValueRating_IsValid/valid_critical
=== RUN   TestValueRating_IsValid/invalid_rating
--- PASS: TestValueRating_IsValid (0.00s)
    --- PASS: TestValueRating_IsValid/valid_low (0.00s)
    --- PASS: TestValueRating_IsValid/valid_medium (0.00s)
    --- PASS: TestValueRating_IsValid/valid_high (0.00s)
    --- PASS: TestValueRating_IsValid/valid_critical (0.00s)
    --- PASS: TestValueRating_IsValid/invalid_rating (0.00s)
=== RUN   TestImpactLevel_IsValid
=== RUN   TestImpactLevel_IsValid/valid_personal
=== RUN   TestImpactLevel_IsValid/valid_team
=== RUN   TestImpactLevel_IsValid/valid_department
=== RUN   TestImpactLevel_IsValid/valid_company
=== RUN   TestImpactLevel_IsValid/invalid_impact
--- PASS: TestImpactLevel_IsValid (0.00s)
    --- PASS: TestImpactLevel_IsValid/valid_personal (0.00s)
    --- PASS: TestImpactLevel_IsValid/valid_team (0.00s)
    --- PASS: TestImpactLevel_IsValid/valid_department (0.00s)
    --- PASS: TestImpactLevel_IsValid/valid_company (0.00s)
    --- PASS: TestImpactLevel_IsValid/invalid_impact (0.00s)
=== RUN   TestLogEntry_CalculateDuration
--- PASS: TestLogEntry_CalculateDuration (0.00s)
=== RUN   TestLogEntry_Validate
=== RUN   TestLogEntry_Validate/valid_log_entry
=== RUN   TestLogEntry_Validate/invalid_time_range
=== RUN   TestLogEntry_Validate/invalid_activity_type
=== RUN   TestLogEntry_Validate/invalid_value_rating
=== RUN   TestLogEntry_Validate/invalid_impact_level
--- PASS: TestLogEntry_Validate (0.00s)
    --- PASS: TestLogEntry_Validate/valid_log_entry (0.00s)
    --- PASS: TestLogEntry_Validate/invalid_time_range (0.00s)
    --- PASS: TestLogEntry_Validate/invalid_activity_type (0.00s)
    --- PASS: TestLogEntry_Validate/invalid_value_rating (0.00s)
    --- PASS: TestLogEntry_Validate/invalid_impact_level (0.00s)
=== RUN   TestReportType_IsValid
=== RUN   TestReportType_IsValid/valid_daily_summary
=== RUN   TestReportType_IsValid/valid_weekly_summary
=== RUN   TestReportType_IsValid/valid_monthly_summary
=== RUN   TestReportType_IsValid/valid_quarterly_summary
=== RUN   TestReportType_IsValid/valid_project_analysis
=== RUN   TestReportType_IsValid/valid_productivity_trends
=== RUN   TestReportType_IsValid/valid_time_distribution
=== RUN   TestReportType_IsValid/valid_performance_review
=== RUN   TestReportType_IsValid/valid_goal_progress
=== RUN   TestReportType_IsValid/valid_custom
=== RUN   TestReportType_IsValid/invalid_type
--- PASS: TestReportType_IsValid (0.00s)
    --- PASS: TestReportType_IsValid/valid_daily_summary (0.00s)
    --- PASS: TestReportType_IsValid/valid_weekly_summary (0.00s)
    --- PASS: TestReportType_IsValid/valid_monthly_summary (0.00s)
    --- PASS: TestReportType_IsValid/valid_quarterly_summary (0.00s)
    --- PASS: TestReportType_IsValid/valid_project_analysis (0.00s)
    --- PASS: TestReportType_IsValid/valid_productivity_trends (0.00s)
    --- PASS: TestReportType_IsValid/valid_time_distribution (0.00s)
    --- PASS: TestReportType_IsValid/valid_performance_review (0.00s)
    --- PASS: TestReportType_IsValid/valid_goal_progress (0.00s)
    --- PASS: TestReportType_IsValid/valid_custom (0.00s)
    --- PASS: TestReportType_IsValid/invalid_type (0.00s)
=== RUN   TestInsightStatus_IsValid
=== RUN   TestInsightStatus_IsValid/valid_active
=== RUN   TestInsightStatus_IsValid/valid_archived
=== RUN   TestInsightStatus_IsValid/valid_superseded
=== RUN   TestInsightStatus_IsValid/invalid_status
--- PASS: TestInsightStatus_IsValid (0.00s)
    --- PASS: TestInsightStatus_IsValid/valid_active (0.00s)
    --- PASS: TestInsightStatus_IsValid/valid_archived (0.00s)
    --- PASS: TestInsightStatus_IsValid/valid_superseded (0.00s)
    --- PASS: TestInsightStatus_IsValid/invalid_status (0.00s)
=== RUN   TestTaskType_IsValid
=== RUN   TestTaskType_IsValid/valid_generate_insight
=== RUN   TestTaskType_IsValid/valid_send_email
=== RUN   TestTaskType_IsValid/valid_export_data
=== RUN   TestTaskType_IsValid/valid_cleanup_data
=== RUN   TestTaskType_IsValid/valid_process_analytics
=== RUN   TestTaskType_IsValid/valid_generate_report
=== RUN   TestTaskType_IsValid/valid_backup_data
=== RUN   TestTaskType_IsValid/valid_custom
=== RUN   TestTaskType_IsValid/invalid_type
--- PASS: TestTaskType_IsValid (0.00s)
    --- PASS: TestTaskType_IsValid/valid_generate_insight (0.00s)
    --- PASS: TestTaskType_IsValid/valid_send_email (0.00s)
    --- PASS: TestTaskType_IsValid/valid_export_data (0.00s)
    --- PASS: TestTaskType_IsValid/valid_cleanup_data (0.00s)
    --- PASS: TestTaskType_IsValid/valid_process_analytics (0.00s)
    --- PASS: TestTaskType_IsValid/valid_generate_report (0.00s)
    --- PASS: TestTaskType_IsValid/valid_backup_data (0.00s)
    --- PASS: TestTaskType_IsValid/valid_custom (0.00s)
    --- PASS: TestTaskType_IsValid/invalid_type (0.00s)
=== RUN   TestTaskStatus_IsValid
=== RUN   TestTaskStatus_IsValid/valid_pending
=== RUN   TestTaskStatus_IsValid/valid_processing
=== RUN   TestTaskStatus_IsValid/valid_completed
=== RUN   TestTaskStatus_IsValid/valid_failed
=== RUN   TestTaskStatus_IsValid/valid_cancelled
=== RUN   TestTaskStatus_IsValid/valid_retrying
=== RUN   TestTaskStatus_IsValid/invalid_status
--- PASS: TestTaskStatus_IsValid (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_pending (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_processing (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_completed (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_failed (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_cancelled (0.00s)
    --- PASS: TestTaskStatus_IsValid/valid_retrying (0.00s)
    --- PASS: TestTaskStatus_IsValid/invalid_status (0.00s)
=== RUN   TestGeneratedInsight_Validate
=== RUN   TestGeneratedInsight_Validate/valid_insight
=== RUN   TestGeneratedInsight_Validate/invalid_time_range
=== RUN   TestGeneratedInsight_Validate/invalid_report_type
=== RUN   TestGeneratedInsight_Validate/invalid_status
=== RUN   TestGeneratedInsight_Validate/invalid_quality_score_-_too_low
=== RUN   TestGeneratedInsight_Validate/invalid_quality_score_-_too_high
--- PASS: TestGeneratedInsight_Validate (0.00s)
    --- PASS: TestGeneratedInsight_Validate/valid_insight (0.00s)
    --- PASS: TestGeneratedInsight_Validate/invalid_time_range (0.00s)
    --- PASS: TestGeneratedInsight_Validate/invalid_report_type (0.00s)
    --- PASS: TestGeneratedInsight_Validate/invalid_status (0.00s)
    --- PASS: TestGeneratedInsight_Validate/invalid_quality_score_-_too_low (0.00s)
    --- PASS: TestGeneratedInsight_Validate/invalid_quality_score_-_too_high (0.00s)
=== RUN   TestTask_Validate
=== RUN   TestTask_Validate/valid_task
=== RUN   TestTask_Validate/invalid_task_type
=== RUN   TestTask_Validate/invalid_status
=== RUN   TestTask_Validate/invalid_priority_-_too_low
=== RUN   TestTask_Validate/invalid_priority_-_too_high
--- PASS: TestTask_Validate (0.00s)
    --- PASS: TestTask_Validate/valid_task (0.00s)
    --- PASS: TestTask_Validate/invalid_task_type (0.00s)
    --- PASS: TestTask_Validate/invalid_status (0.00s)
    --- PASS: TestTask_Validate/invalid_priority_-_too_low (0.00s)
    --- PASS: TestTask_Validate/invalid_priority_-_too_high (0.00s)
=== RUN   TestUserJSON
--- PASS: TestUserJSON (0.00s)
=== RUN   TestLogEntryJSON
--- PASS: TestLogEntryJSON (0.00s)
=== RUN   TestProjectJSON
--- PASS: TestProjectJSON (0.00s)
=== RUN   TestGeneratedInsightJSON
--- PASS: TestGeneratedInsightJSON (0.00s)
=== RUN   TestProjectStatus_IsValid
=== RUN   TestProjectStatus_IsValid/valid_active
=== RUN   TestProjectStatus_IsValid/valid_completed
=== RUN   TestProjectStatus_IsValid/valid_on_hold
=== RUN   TestProjectStatus_IsValid/valid_cancelled
=== RUN   TestProjectStatus_IsValid/invalid_status
--- PASS: TestProjectStatus_IsValid (0.00s)
    --- PASS: TestProjectStatus_IsValid/valid_active (0.00s)
    --- PASS: TestProjectStatus_IsValid/valid_completed (0.00s)
    --- PASS: TestProjectStatus_IsValid/valid_on_hold (0.00s)
    --- PASS: TestProjectStatus_IsValid/valid_cancelled (0.00s)
    --- PASS: TestProjectStatus_IsValid/invalid_status (0.00s)
=== RUN   TestProject_Validate
=== RUN   TestProject_Validate/valid_project
=== RUN   TestProject_Validate/invalid_status
--- PASS: TestProject_Validate (0.00s)
    --- PASS: TestProject_Validate/valid_project (0.00s)
    --- PASS: TestProject_Validate/invalid_status (0.00s)
=== RUN   TestValidateTimezone
=== RUN   TestValidateTimezone/valid_UTC
=== RUN   TestValidateTimezone/valid_America/New_York
=== RUN   TestValidateTimezone/valid_Europe/London
=== RUN   TestValidateTimezone/valid_Asia/Tokyo
=== RUN   TestValidateTimezone/invalid_timezone
=== RUN   TestValidateTimezone/empty_timezone_defaults_to_UTC
--- PASS: TestValidateTimezone (0.00s)
    --- PASS: TestValidateTimezone/valid_UTC (0.00s)
    --- PASS: TestValidateTimezone/valid_America/New_York (0.00s)
    --- PASS: TestValidateTimezone/valid_Europe/London (0.00s)
    --- PASS: TestValidateTimezone/valid_Asia/Tokyo (0.00s)
    --- PASS: TestValidateTimezone/invalid_timezone (0.00s)
    --- PASS: TestValidateTimezone/empty_timezone_defaults_to_UTC (0.00s)
=== RUN   TestValidateHexColor
=== RUN   TestValidateHexColor/valid_color_lowercase
=== RUN   TestValidateHexColor/valid_color_uppercase
=== RUN   TestValidateHexColor/valid_color_mixed_case
=== RUN   TestValidateHexColor/invalid_without_hash
=== RUN   TestValidateHexColor/invalid_short_format
=== RUN   TestValidateHexColor/invalid_long_format
=== RUN   TestValidateHexColor/invalid_characters
=== RUN   TestValidateHexColor/empty_color
--- PASS: TestValidateHexColor (0.00s)
    --- PASS: TestValidateHexColor/valid_color_lowercase (0.00s)
    --- PASS: TestValidateHexColor/valid_color_uppercase (0.00s)
    --- PASS: TestValidateHexColor/valid_color_mixed_case (0.00s)
    --- PASS: TestValidateHexColor/invalid_without_hash (0.00s)
    --- PASS: TestValidateHexColor/invalid_short_format (0.00s)
    --- PASS: TestValidateHexColor/invalid_long_format (0.00s)
    --- PASS: TestValidateHexColor/invalid_characters (0.00s)
    --- PASS: TestValidateHexColor/empty_color (0.00s)
=== RUN   TestValidateTimeRange
=== RUN   TestValidateTimeRange/valid_range
=== RUN   TestValidateTimeRange/same_time
=== RUN   TestValidateTimeRange/invalid_range
--- PASS: TestValidateTimeRange (0.00s)
    --- PASS: TestValidateTimeRange/valid_range (0.00s)
    --- PASS: TestValidateTimeRange/same_time (0.00s)
    --- PASS: TestValidateTimeRange/invalid_range (0.00s)
=== RUN   TestValidateActivityType
=== RUN   TestValidateActivityType/empty_activity_type_(optional)
=== RUN   TestValidateActivityType/valid_activity_type_-_development
=== RUN   TestValidateActivityType/valid_activity_type_-_meeting
=== RUN   TestValidateActivityType/valid_activity_type_-_code_review
=== RUN   TestValidateActivityType/valid_activity_type_-_debugging
=== RUN   TestValidateActivityType/valid_activity_type_-_research
=== RUN   TestValidateActivityType/valid_activity_type_-_testing
=== RUN   TestValidateActivityType/valid_activity_type_-_documentation
=== RUN   TestValidateActivityType/valid_activity_type_-_deployment
=== RUN   TestValidateActivityType/valid_activity_type_-_learning
=== RUN   TestValidateActivityType/valid_activity_type_-_planning
=== RUN   TestValidateActivityType/valid_activity_type_-_maintenance
=== RUN   TestValidateActivityType/valid_activity_type_-_support
=== RUN   TestValidateActivityType/valid_activity_type_-_other
=== RUN   TestValidateActivityType/invalid_activity_type
=== RUN   TestValidateActivityType/case_sensitive_validation
--- PASS: TestValidateActivityType (0.00s)
    --- PASS: TestValidateActivityType/empty_activity_type_(optional) (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_development (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_meeting (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_code_review (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_debugging (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_research (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_testing (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_documentation (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_deployment (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_learning (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_planning (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_maintenance (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_support (0.00s)
    --- PASS: TestValidateActivityType/valid_activity_type_-_other (0.00s)
    --- PASS: TestValidateActivityType/invalid_activity_type (0.00s)
    --- PASS: TestValidateActivityType/case_sensitive_validation (0.00s)
=== RUN   TestValidateValueRating
=== RUN   TestValidateValueRating/empty_value_rating_(optional)
=== RUN   TestValidateValueRating/valid_value_rating_-_low
=== RUN   TestValidateValueRating/valid_value_rating_-_medium
=== RUN   TestValidateValueRating/valid_value_rating_-_high
=== RUN   TestValidateValueRating/valid_value_rating_-_critical
=== RUN   TestValidateValueRating/invalid_value_rating
=== RUN   TestValidateValueRating/case_sensitive_validation
--- PASS: TestValidateValueRating (0.00s)
    --- PASS: TestValidateValueRating/empty_value_rating_(optional) (0.00s)
    --- PASS: TestValidateValueRating/valid_value_rating_-_low (0.00s)
    --- PASS: TestValidateValueRating/valid_value_rating_-_medium (0.00s)
    --- PASS: TestValidateValueRating/valid_value_rating_-_high (0.00s)
    --- PASS: TestValidateValueRating/valid_value_rating_-_critical (0.00s)
    --- PASS: TestValidateValueRating/invalid_value_rating (0.00s)
    --- PASS: TestValidateValueRating/case_sensitive_validation (0.00s)
=== RUN   TestValidateImpactLevel
=== RUN   TestValidateImpactLevel/empty_impact_level_(optional)
=== RUN   TestValidateImpactLevel/valid_impact_level_-_personal
=== RUN   TestValidateImpactLevel/valid_impact_level_-_team
=== RUN   TestValidateImpactLevel/valid_impact_level_-_department
=== RUN   TestValidateImpactLevel/valid_impact_level_-_company
=== RUN   TestValidateImpactLevel/invalid_impact_level
=== RUN   TestValidateImpactLevel/case_sensitive_validation
--- PASS: TestValidateImpactLevel (0.00s)
    --- PASS: TestValidateImpactLevel/empty_impact_level_(optional) (0.00s)
    --- PASS: TestValidateImpactLevel/valid_impact_level_-_personal (0.00s)
    --- PASS: TestValidateImpactLevel/valid_impact_level_-_team (0.00s)
    --- PASS: TestValidateImpactLevel/valid_impact_level_-_department (0.00s)
    --- PASS: TestValidateImpactLevel/valid_impact_level_-_company (0.00s)
    --- PASS: TestValidateImpactLevel/invalid_impact_level (0.00s)
    --- PASS: TestValidateImpactLevel/case_sensitive_validation (0.00s)
=== RUN   TestValidateDateFormat
=== RUN   TestValidateDateFormat/empty_date_(optional)
=== RUN   TestValidateDateFormat/valid_date_format
=== RUN   TestValidateDateFormat/valid_date_format_-_different_date
=== RUN   TestValidateDateFormat/invalid_date_format_-_wrong_separator
=== RUN   TestValidateDateFormat/invalid_date_format_-_missing_day
=== RUN   TestValidateDateFormat/invalid_date_format_-_wrong_order
=== RUN   TestValidateDateFormat/invalid_date_format_-_with_time
=== RUN   TestValidateDateFormat/invalid_date_-_impossible_date
=== RUN   TestValidateDateFormat/invalid_date_-_wrong_month
=== RUN   TestValidateDateFormat/invalid_date_format_-_text
--- PASS: TestValidateDateFormat (0.00s)
    --- PASS: TestValidateDateFormat/empty_date_(optional) (0.00s)
    --- PASS: TestValidateDateFormat/valid_date_format (0.00s)
    --- PASS: TestValidateDateFormat/valid_date_format_-_different_date (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_format_-_wrong_separator (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_format_-_missing_day (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_format_-_wrong_order (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_format_-_with_time (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_-_impossible_date (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_-_wrong_month (0.00s)
    --- PASS: TestValidateDateFormat/invalid_date_format_-_text (0.00s)
=== RUN   TestMiddlewareEdgeCases
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/non-leap_year_February_29th
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/leap_year_February_29th
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/April_31st_(April_has_30_days)
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/valid_December_31st
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/year_zero
=== RUN   TestMiddlewareEdgeCases/date_validation_edge_cases/year_9999
--- PASS: TestMiddlewareEdgeCases (0.00s)
    --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/non-leap_year_February_29th (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/leap_year_February_29th (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/April_31st_(April_has_30_days) (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/valid_December_31st (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/year_zero (0.00s)
        --- PASS: TestMiddlewareEdgeCases/date_validation_edge_cases/year_9999 (0.00s)
PASS
ok      github.com/garnizeh/englog/internal/models      0.011s
=== RUN   TestAnalyticsService_ComprehensiveValidation
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_same_day
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_week
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_month
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_year
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/exact_same_timestamp
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/end_1_second_before_start
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/end_1_day_before_start
=== RUN   TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/very_large_range
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/valid_UUID
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/valid_UUID_with_uppercase
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/nil_UUID
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/empty_string
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_too_short
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_missing_hyphens
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_extra_characters
=== RUN   TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_characters
--- PASS: TestAnalyticsService_ComprehensiveValidation (0.00s)
    --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_same_day (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_week (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_month (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/valid_range_-_one_year (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/exact_same_timestamp (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/end_1_second_before_start (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/end_1_day_before_start (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/DateRangeValidationEdgeCases/very_large_range (0.00s)
    --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/valid_UUID (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/valid_UUID_with_uppercase (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/nil_UUID (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/empty_string (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_too_short (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_missing_hyphens (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_format_-_extra_characters (0.00s)
        --- PASS: TestAnalyticsService_ComprehensiveValidation/UserIDValidationEdgeCases/invalid_characters (0.00s)
=== RUN   TestAnalyticsService_MetricCalculations
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/perfect_score
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/half_high_value
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/no_high_value
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/no_activity
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/quarter_high_value
=== RUN   TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/three_quarters_high_value
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/single_value
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/multiple_values
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/empty_slice
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/values_with_decimals
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/negative_values
=== RUN   TestAnalyticsService_MetricCalculations/AverageCalculation/zero_values
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/even_distribution
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/uneven_distribution
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/three_categories
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/empty_map
=== RUN   TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/single_category
--- PASS: TestAnalyticsService_MetricCalculations (0.00s)
    --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/perfect_score (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/half_high_value (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/no_high_value (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/no_activity (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/quarter_high_value (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/ProductivityScoreCalculation/three_quarters_high_value (0.00s)
    --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/single_value (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/multiple_values (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/empty_slice (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/values_with_decimals (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/negative_values (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/AverageCalculation/zero_values (0.00s)
    --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/even_distribution (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/uneven_distribution (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/three_categories (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/empty_map (0.00s)
        --- PASS: TestAnalyticsService_MetricCalculations/PercentageDistributionCalculation/single_category (0.00s)
=== RUN   TestAnalyticsService_TimeZoneHandling
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/UTC
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/New_York
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/Tokyo
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/London
=== RUN   TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/Los_Angeles
--- PASS: TestAnalyticsService_TimeZoneHandling (0.00s)
    --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones (0.00s)
        --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/UTC (0.00s)
        --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/New_York (0.00s)
        --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/Tokyo (0.00s)
        --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/London (0.00s)
        --- PASS: TestAnalyticsService_TimeZoneHandling/DifferentTimeZones/Los_Angeles (0.00s)
=== RUN   TestAnalyticsService_DataAggregation
=== RUN   TestAnalyticsService_DataAggregation/WeeklyAggregation
=== RUN   TestAnalyticsService_DataAggregation/WeeklyAggregation/Monday
=== RUN   TestAnalyticsService_DataAggregation/WeeklyAggregation/Tuesday
=== RUN   TestAnalyticsService_DataAggregation/WeeklyAggregation/Sunday
=== RUN   TestAnalyticsService_DataAggregation/MonthlyAggregation
=== RUN   TestAnalyticsService_DataAggregation/MonthlyAggregation/first_day_of_month
=== RUN   TestAnalyticsService_DataAggregation/MonthlyAggregation/middle_of_month
=== RUN   TestAnalyticsService_DataAggregation/MonthlyAggregation/last_day_of_month
--- PASS: TestAnalyticsService_DataAggregation (0.00s)
    --- PASS: TestAnalyticsService_DataAggregation/WeeklyAggregation (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/WeeklyAggregation/Monday (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/WeeklyAggregation/Tuesday (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/WeeklyAggregation/Sunday (0.00s)
    --- PASS: TestAnalyticsService_DataAggregation/MonthlyAggregation (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/MonthlyAggregation/first_day_of_month (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/MonthlyAggregation/middle_of_month (0.00s)
        --- PASS: TestAnalyticsService_DataAggregation/MonthlyAggregation/last_day_of_month (0.00s)
=== RUN   TestAnalyticsService_BoundaryConditions
=== RUN   TestAnalyticsService_BoundaryConditions/ZeroValues
=== RUN   TestAnalyticsService_BoundaryConditions/LargeValues
=== RUN   TestAnalyticsService_BoundaryConditions/EdgeCaseDates
--- PASS: TestAnalyticsService_BoundaryConditions (0.00s)
    --- PASS: TestAnalyticsService_BoundaryConditions/ZeroValues (0.00s)
    --- PASS: TestAnalyticsService_BoundaryConditions/LargeValues (0.00s)
    --- PASS: TestAnalyticsService_BoundaryConditions/EdgeCaseDates (0.00s)
=== RUN   TestAnalyticsService_ValidationMethods
=== RUN   TestAnalyticsService_ValidationMethods/ValidateDateRange
=== RUN   TestAnalyticsService_ValidationMethods/ValidateDateRange/valid_date_range
=== RUN   TestAnalyticsService_ValidationMethods/ValidateDateRange/same_start_and_end_date
=== RUN   TestAnalyticsService_ValidationMethods/ValidateDateRange/end_before_start
=== RUN   TestAnalyticsService_ValidationMethods/ValidateDateRange/very_long_range
--- PASS: TestAnalyticsService_ValidationMethods (0.00s)
    --- PASS: TestAnalyticsService_ValidationMethods/ValidateDateRange (0.00s)
        --- PASS: TestAnalyticsService_ValidationMethods/ValidateDateRange/valid_date_range (0.00s)
        --- PASS: TestAnalyticsService_ValidationMethods/ValidateDateRange/same_start_and_end_date (0.00s)
        --- PASS: TestAnalyticsService_ValidationMethods/ValidateDateRange/end_before_start (0.00s)
        --- PASS: TestAnalyticsService_ValidationMethods/ValidateDateRange/very_long_range (0.00s)
=== RUN   TestAnalyticsService_HelperMethods
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/50%_increase
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/25%_decrease
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/no_change
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/from_zero
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/to_zero
=== RUN   TestAnalyticsService_HelperMethods/CalculatePercentageChange/both_zero
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore/normal_score
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore/negative_score
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore/over_100_score
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore/zero_score
=== RUN   TestAnalyticsService_HelperMethods/NormalizeProductivityScore/perfect_score
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Sunday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Monday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Tuesday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Wednesday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Thursday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Friday
=== RUN   TestAnalyticsService_HelperMethods/GetDayOfWeekName/Saturday
--- PASS: TestAnalyticsService_HelperMethods (0.00s)
    --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/50%_increase (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/25%_decrease (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/no_change (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/from_zero (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/to_zero (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/CalculatePercentageChange/both_zero (0.00s)
    --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore/normal_score (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore/negative_score (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore/over_100_score (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore/zero_score (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/NormalizeProductivityScore/perfect_score (0.00s)
    --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Sunday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Monday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Tuesday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Wednesday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Thursday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Friday (0.00s)
        --- PASS: TestAnalyticsService_HelperMethods/GetDayOfWeekName/Saturday (0.00s)
=== RUN   TestLogEntryService_ComprehensiveValidation
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/empty_title
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/whitespace_only_title
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/single_character_title
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/maximum_length_title
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_too_long
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_special_characters
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_unicode
=== RUN   TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_newlines
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/nil_description
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/empty_description
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/whitespace_only_description
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/maximum_length_description
=== RUN   TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/description_too_long
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/valid_time_range
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/same_start_and_end_time
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/end_before_start
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/1_minute_duration
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/1_second_duration
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/24_hour_duration
=== RUN   TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/very_long_duration
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/development_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/meeting_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/code_review_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/debugging_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/documentation_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/testing_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/deployment_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/research_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/planning_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/learning_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/maintenance_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/support_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/other_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/invalid_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/empty_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/case_sensitive_type
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/low_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/medium_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/high_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/critical_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/invalid_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/empty_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/case_sensitive_value
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/personal_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/team_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/department_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/company_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/invalid_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/empty_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/case_sensitive_impact
=== RUN   TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/nil_project_ID
=== RUN   TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/valid_project_ID
=== RUN   TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/zero_UUID
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/nil_tags
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/empty_tags
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/single_tag
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/multiple_tags
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/tags_with_special_characters
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/tags_with_unicode
=== RUN   TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/empty_tag_in_list
--- PASS: TestLogEntryService_ComprehensiveValidation (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/empty_title (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/whitespace_only_title (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/single_character_title (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/maximum_length_title (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_too_long (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_special_characters (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_unicode (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TitleValidationEdgeCases/title_with_newlines (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/nil_description (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/empty_description (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/whitespace_only_description (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/maximum_length_description (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/DescriptionValidationEdgeCases/description_too_long (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/valid_time_range (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/same_start_and_end_time (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/end_before_start (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/1_minute_duration (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/1_second_duration (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/24_hour_duration (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TimeValidationEdgeCases/very_long_duration (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/development_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/meeting_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/code_review_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/debugging_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/documentation_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/testing_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/deployment_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/research_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/planning_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/learning_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/maintenance_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/support_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/other_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/invalid_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/empty_type (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ActivityTypeValidationEdgeCases/case_sensitive_type (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/low_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/medium_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/high_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/critical_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/invalid_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/empty_value (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ValueRatingValidationEdgeCases/case_sensitive_value (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/personal_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/team_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/department_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/company_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/invalid_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/empty_impact (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ImpactLevelValidationEdgeCases/case_sensitive_impact (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/nil_project_ID (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/valid_project_ID (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/ProjectIDValidationEdgeCases/zero_UUID (0.00s)
    --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/nil_tags (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/empty_tags (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/single_tag (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/multiple_tags (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/tags_with_special_characters (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/tags_with_unicode (0.00s)
        --- PASS: TestLogEntryService_ComprehensiveValidation/TagsValidationEdgeCases/empty_tag_in_list (0.00s)
=== RUN   TestLogEntryService_DurationCalculation
=== RUN   TestLogEntryService_DurationCalculation/1_minute
=== RUN   TestLogEntryService_DurationCalculation/30_minutes
=== RUN   TestLogEntryService_DurationCalculation/1_hour
=== RUN   TestLogEntryService_DurationCalculation/1.5_hours
=== RUN   TestLogEntryService_DurationCalculation/8_hours
=== RUN   TestLogEntryService_DurationCalculation/partial_minute_(30_seconds)
=== RUN   TestLogEntryService_DurationCalculation/partial_minute_(90_seconds)
--- PASS: TestLogEntryService_DurationCalculation (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/1_minute (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/30_minutes (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/1_hour (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/1.5_hours (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/8_hours (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/partial_minute_(30_seconds) (0.00s)
    --- PASS: TestLogEntryService_DurationCalculation/partial_minute_(90_seconds) (0.00s)
=== RUN   TestLogEntryService_BoundaryConditions
=== RUN   TestLogEntryService_BoundaryConditions/MaxFieldLengthCombinations
=== RUN   TestLogEntryService_BoundaryConditions/EmptyOptionalFieldsCombinations
--- PASS: TestLogEntryService_BoundaryConditions (0.00s)
    --- PASS: TestLogEntryService_BoundaryConditions/MaxFieldLengthCombinations (0.00s)
    --- PASS: TestLogEntryService_BoundaryConditions/EmptyOptionalFieldsCombinations (0.00s)
=== RUN   TestLogEntryService_ValidationMethods
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/valid_request_minimal
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/valid_request_complete
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/missing_title
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/title_too_long
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_activity_type
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_value_rating
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_impact_level
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/end_time_before_start_time
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/same_start_and_end_time
=== RUN   TestLogEntryService_ValidationMethods/validateLogEntryRequest/description_too_long
--- PASS: TestLogEntryService_ValidationMethods (0.00s)
    --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/valid_request_minimal (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/valid_request_complete (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/missing_title (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/title_too_long (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_activity_type (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_value_rating (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/invalid_impact_level (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/end_time_before_start_time (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/same_start_and_end_time (0.00s)
        --- PASS: TestLogEntryService_ValidationMethods/validateLogEntryRequest/description_too_long (0.00s)
=== RUN   TestLogEntryService_ActivityTypeValidation
=== RUN   TestLogEntryService_ActivityTypeValidation/development_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/meeting_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/learning_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/debugging_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/other_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/invalid_activity
=== RUN   TestLogEntryService_ActivityTypeValidation/empty_activity
--- PASS: TestLogEntryService_ActivityTypeValidation (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/development_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/meeting_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/learning_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/debugging_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/other_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/invalid_activity (0.00s)
    --- PASS: TestLogEntryService_ActivityTypeValidation/empty_activity (0.00s)
=== RUN   TestLogEntryService_ValueRatingValidation
=== RUN   TestLogEntryService_ValueRatingValidation/low_value
=== RUN   TestLogEntryService_ValueRatingValidation/medium_value
=== RUN   TestLogEntryService_ValueRatingValidation/high_value
=== RUN   TestLogEntryService_ValueRatingValidation/critical_value
=== RUN   TestLogEntryService_ValueRatingValidation/invalid_value
=== RUN   TestLogEntryService_ValueRatingValidation/empty_value
--- PASS: TestLogEntryService_ValueRatingValidation (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/low_value (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/medium_value (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/high_value (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/critical_value (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/invalid_value (0.00s)
    --- PASS: TestLogEntryService_ValueRatingValidation/empty_value (0.00s)
=== RUN   TestLogEntryService_ImpactLevelValidation
=== RUN   TestLogEntryService_ImpactLevelValidation/personal_impact
=== RUN   TestLogEntryService_ImpactLevelValidation/team_impact
=== RUN   TestLogEntryService_ImpactLevelValidation/department_impact
=== RUN   TestLogEntryService_ImpactLevelValidation/company_impact
=== RUN   TestLogEntryService_ImpactLevelValidation/invalid_impact
=== RUN   TestLogEntryService_ImpactLevelValidation/empty_impact
--- PASS: TestLogEntryService_ImpactLevelValidation (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/personal_impact (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/team_impact (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/department_impact (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/company_impact (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/invalid_impact (0.00s)
    --- PASS: TestLogEntryService_ImpactLevelValidation/empty_impact (0.00s)
=== RUN   TestProjectService_BusinessLogicScenarios
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/single_character_name
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/name_with_spaces
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/name_with_special_chars
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/unicode_name
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/exactly_200_chars
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/201_chars_should_fail
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/empty_name
=== RUN   TestProjectService_BusinessLogicScenarios/project_name_edge_cases/whitespace_only_name
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_short
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_long
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_lowercase
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_mixed_case
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_with_numbers
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_no_hash
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_wrong_length
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_invalid_chars
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/empty_color
=== RUN   TestProjectService_BusinessLogicScenarios/project_color_validation/whitespace_color
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/active_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/completed_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/on_hold_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/cancelled_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/invalid_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_status_validation/empty_status
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/nil_description
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/empty_description
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/normal_description
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/unicode_description
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/description_with_newlines
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/description_with_special_chars
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/exactly_1000_chars
=== RUN   TestProjectService_BusinessLogicScenarios/project_description_validation/1001_chars_should_fail
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/both_dates_nil
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/only_start_date
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/only_end_date
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/valid_date_range
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/same_start_and_end_date
=== RUN   TestProjectService_BusinessLogicScenarios/project_dates_validation/end_date_before_start_date
=== RUN   TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios
=== RUN   TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios/default_project
=== RUN   TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios/non_default_project
--- PASS: TestProjectService_BusinessLogicScenarios (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/single_character_name (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/name_with_spaces (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/name_with_special_chars (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/unicode_name (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/exactly_200_chars (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/201_chars_should_fail (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/empty_name (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_name_edge_cases/whitespace_only_name (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_short (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_long (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_lowercase (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_mixed_case (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/valid_hex_color_with_numbers (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_no_hash (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_wrong_length (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/invalid_color_invalid_chars (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/empty_color (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_color_validation/whitespace_color (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/active_status (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/completed_status (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/on_hold_status (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/cancelled_status (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/invalid_status (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_status_validation/empty_status (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/nil_description (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/empty_description (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/normal_description (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/unicode_description (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/description_with_newlines (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/description_with_special_chars (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/exactly_1000_chars (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_description_validation/1001_chars_should_fail (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/both_dates_nil (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/only_start_date (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/only_end_date (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/valid_date_range (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/same_start_and_end_date (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_dates_validation/end_date_before_start_date (0.00s)
    --- PASS: TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios/default_project (0.00s)
        --- PASS: TestProjectService_BusinessLogicScenarios/project_default_flag_scenarios/non_default_project (0.00s)
=== RUN   TestProjectService_StatusTransitions
=== RUN   TestProjectService_StatusTransitions/transition_completed_to_active
=== RUN   TestProjectService_StatusTransitions/transition_cancelled_to_active
=== RUN   TestProjectService_StatusTransitions/transition_active_to_completed
=== RUN   TestProjectService_StatusTransitions/transition_active_to_on_hold
=== RUN   TestProjectService_StatusTransitions/transition_active_to_cancelled
=== RUN   TestProjectService_StatusTransitions/transition_on_hold_to_active
=== RUN   TestProjectService_StatusTransitions/transition_on_hold_to_completed
=== RUN   TestProjectService_StatusTransitions/transition_on_hold_to_cancelled
--- PASS: TestProjectService_StatusTransitions (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_completed_to_active (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_cancelled_to_active (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_active_to_completed (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_active_to_on_hold (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_active_to_cancelled (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_on_hold_to_active (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_on_hold_to_completed (0.00s)
    --- PASS: TestProjectService_StatusTransitions/transition_on_hold_to_cancelled (0.00s)
=== RUN   TestProjectService_ComplexValidationScenarios
=== RUN   TestProjectService_ComplexValidationScenarios/complete_project_validation
=== RUN   TestProjectService_ComplexValidationScenarios/complete_project_validation/fully_valid_project
=== RUN   TestProjectService_ComplexValidationScenarios/complete_project_validation/minimal_valid_project
=== RUN   TestProjectService_ComplexValidationScenarios/complete_project_validation/multiple_validation_errors
--- PASS: TestProjectService_ComplexValidationScenarios (0.00s)
    --- PASS: TestProjectService_ComplexValidationScenarios/complete_project_validation (0.00s)
        --- PASS: TestProjectService_ComplexValidationScenarios/complete_project_validation/fully_valid_project (0.00s)
        --- PASS: TestProjectService_ComplexValidationScenarios/complete_project_validation/minimal_valid_project (0.00s)
        --- PASS: TestProjectService_ComplexValidationScenarios/complete_project_validation/multiple_validation_errors (0.00s)
=== RUN   TestProjectService_ValidationMethods
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/valid_request
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/valid_request_with_description
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/valid_request_with_all_fields
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/missing_name
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/name_too_long
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/missing_color
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/invalid_color_format
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/invalid_status
=== RUN   TestProjectService_ValidationMethods/validateProjectRequest/description_too_long
--- PASS: TestProjectService_ValidationMethods (0.00s)
    --- PASS: TestProjectService_ValidationMethods/validateProjectRequest (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/valid_request (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/valid_request_with_description (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/valid_request_with_all_fields (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/missing_name (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/name_too_long (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/missing_color (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/invalid_color_format (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/invalid_status (0.00s)
        --- PASS: TestProjectService_ValidationMethods/validateProjectRequest/description_too_long (0.00s)
=== RUN   TestProjectService_ProjectStatusValidation
=== RUN   TestProjectService_ProjectStatusValidation/active_status
=== RUN   TestProjectService_ProjectStatusValidation/completed_status
=== RUN   TestProjectService_ProjectStatusValidation/on_hold_status
=== RUN   TestProjectService_ProjectStatusValidation/cancelled_status
=== RUN   TestProjectService_ProjectStatusValidation/invalid_status
=== RUN   TestProjectService_ProjectStatusValidation/empty_status
--- PASS: TestProjectService_ProjectStatusValidation (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/active_status (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/completed_status (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/on_hold_status (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/cancelled_status (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/invalid_status (0.00s)
    --- PASS: TestProjectService_ProjectStatusValidation/empty_status (0.00s)
=== RUN   TestTagService_BusinessLogicScenarios
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/single_character_name
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/name_with_spaces
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/name_with_special_chars
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/unicode_name
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/exactly_100_chars
=== RUN   TestTagService_BusinessLogicScenarios/tag_name_edge_cases/101_chars_should_fail
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/nil_description
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/empty_description
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/single_char_description
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/exactly_500_chars
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/501_chars_should_fail
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/description_with_newlines
=== RUN   TestTagService_BusinessLogicScenarios/description_edge_cases/description_with_unicode
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_red
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_green
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_blue
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_black
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_white
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_mixed_case
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_all_lowercase
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_all_uppercase
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_no_hash
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_too_short
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_too_long
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_non_hex
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_with_spaces
=== RUN   TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_empty
--- PASS: TestTagService_BusinessLogicScenarios (0.00s)
    --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/single_character_name (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/name_with_spaces (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/name_with_special_chars (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/unicode_name (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/exactly_100_chars (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/tag_name_edge_cases/101_chars_should_fail (0.00s)
    --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/nil_description (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/empty_description (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/single_char_description (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/exactly_500_chars (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/501_chars_should_fail (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/description_with_newlines (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/description_edge_cases/description_with_unicode (0.00s)
    --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_red (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_green (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_blue (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_black (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_white (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_mixed_case (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_all_lowercase (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/valid_all_uppercase (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_no_hash (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_too_short (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_too_long (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_non_hex (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_with_spaces (0.00s)
        --- PASS: TestTagService_BusinessLogicScenarios/color_validation_comprehensive/invalid_empty (0.00s)
=== RUN   TestTagService_LimitBoundaryConditions
=== RUN   TestTagService_LimitBoundaryConditions/zero_limit_popular_tags
=== RUN   TestTagService_LimitBoundaryConditions/negative_limit_popular_tags
=== RUN   TestTagService_LimitBoundaryConditions/max_int32_limit
=== RUN   TestTagService_LimitBoundaryConditions/one_limit
=== RUN   TestTagService_LimitBoundaryConditions/zero_limit_search_tags
=== RUN   TestTagService_LimitBoundaryConditions/negative_limit_search_tags
=== RUN   TestTagService_LimitBoundaryConditions/large_limit_search_tags
--- PASS: TestTagService_LimitBoundaryConditions (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/zero_limit_popular_tags (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/negative_limit_popular_tags (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/max_int32_limit (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/one_limit (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/zero_limit_search_tags (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/negative_limit_search_tags (0.00s)
    --- PASS: TestTagService_LimitBoundaryConditions/large_limit_search_tags (0.00s)
=== RUN   TestTagService_ErrorMessagePatterns
=== RUN   TestTagService_ErrorMessagePatterns/empty_name_error
=== RUN   TestTagService_ErrorMessagePatterns/empty_color_error
=== RUN   TestTagService_ErrorMessagePatterns/invalid_color_error
--- PASS: TestTagService_ErrorMessagePatterns (0.00s)
    --- PASS: TestTagService_ErrorMessagePatterns/empty_name_error (0.00s)
    --- PASS: TestTagService_ErrorMessagePatterns/empty_color_error (0.00s)
    --- PASS: TestTagService_ErrorMessagePatterns/invalid_color_error (0.00s)
=== RUN   TestTagService_RequestStructureValidation
=== RUN   TestTagService_RequestStructureValidation/valid_minimal_request
=== RUN   TestTagService_RequestStructureValidation/valid_complete_request
--- PASS: TestTagService_RequestStructureValidation (0.00s)
    --- PASS: TestTagService_RequestStructureValidation/valid_minimal_request (0.00s)
    --- PASS: TestTagService_RequestStructureValidation/valid_complete_request (0.00s)
=== RUN   TestTagService_ValidateTagRequest
=== RUN   TestTagService_ValidateTagRequest/valid_request
=== RUN   TestTagService_ValidateTagRequest/valid_request_with_description
=== RUN   TestTagService_ValidateTagRequest/missing_name
=== RUN   TestTagService_ValidateTagRequest/name_too_long
=== RUN   TestTagService_ValidateTagRequest/missing_color
=== RUN   TestTagService_ValidateTagRequest/invalid_color_format
=== RUN   TestTagService_ValidateTagRequest/description_too_long
=== RUN   TestTagService_ValidateTagRequest/valid_hex_color_six_characters
=== RUN   TestTagService_ValidateTagRequest/valid_hex_color_uppercase
=== RUN   TestTagService_ValidateTagRequest/valid_hex_color_lowercase
--- PASS: TestTagService_ValidateTagRequest (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/valid_request (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/valid_request_with_description (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/missing_name (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/name_too_long (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/missing_color (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/invalid_color_format (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/description_too_long (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/valid_hex_color_six_characters (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/valid_hex_color_uppercase (0.00s)
    --- PASS: TestTagService_ValidateTagRequest/valid_hex_color_lowercase (0.00s)
=== RUN   TestTagService_UUIDParsing
=== RUN   TestTagService_UUIDParsing/valid_UUID
=== RUN   TestTagService_UUIDParsing/invalid_UUID_format
=== RUN   TestTagService_UUIDParsing/empty_UUID
=== RUN   TestTagService_UUIDParsing/malformed_UUID
--- PASS: TestTagService_UUIDParsing (0.00s)
    --- PASS: TestTagService_UUIDParsing/valid_UUID (0.00s)
    --- PASS: TestTagService_UUIDParsing/invalid_UUID_format (0.00s)
    --- PASS: TestTagService_UUIDParsing/empty_UUID (0.00s)
    --- PASS: TestTagService_UUIDParsing/malformed_UUID (0.00s)
=== RUN   TestTagService_LimitValidation
=== RUN   TestTagService_LimitValidation/positive_limit
=== RUN   TestTagService_LimitValidation/zero_limit_should_use_default
=== RUN   TestTagService_LimitValidation/negative_limit_should_use_default
=== RUN   TestTagService_LimitValidation/large_limit
--- PASS: TestTagService_LimitValidation (0.00s)
    --- PASS: TestTagService_LimitValidation/positive_limit (0.00s)
    --- PASS: TestTagService_LimitValidation/zero_limit_should_use_default (0.00s)
    --- PASS: TestTagService_LimitValidation/negative_limit_should_use_default (0.00s)
    --- PASS: TestTagService_LimitValidation/large_limit (0.00s)
=== RUN   TestTagService_SearchLimitValidation
=== RUN   TestTagService_SearchLimitValidation/positive_limit
=== RUN   TestTagService_SearchLimitValidation/zero_limit_should_use_default
=== RUN   TestTagService_SearchLimitValidation/negative_limit_should_use_default
--- PASS: TestTagService_SearchLimitValidation (0.00s)
    --- PASS: TestTagService_SearchLimitValidation/positive_limit (0.00s)
    --- PASS: TestTagService_SearchLimitValidation/zero_limit_should_use_default (0.00s)
    --- PASS: TestTagService_SearchLimitValidation/negative_limit_should_use_default (0.00s)
=== RUN   TestTagService_ErrorHandling
=== RUN   TestTagService_ErrorHandling/create_tag_error
=== RUN   TestTagService_ErrorHandling/get_tag_error
=== RUN   TestTagService_ErrorHandling/delete_tag_error
--- PASS: TestTagService_ErrorHandling (0.00s)
    --- PASS: TestTagService_ErrorHandling/create_tag_error (0.00s)
    --- PASS: TestTagService_ErrorHandling/get_tag_error (0.00s)
    --- PASS: TestTagService_ErrorHandling/delete_tag_error (0.00s)
=== RUN   TestTagService_EnsureTagExistsLogic
--- PASS: TestTagService_EnsureTagExistsLogic (0.00s)
=== RUN   TestUserService_BusinessLogicScenarios
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/single_character_names
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/names_with_spaces
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/names_with_special_chars
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/unicode_names
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/long_timezone
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/whitespace_only_first_name
=== RUN   TestUserService_BusinessLogicScenarios/profile_request_edge_cases/whitespace_only_last_name
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/exactly_8_chars_password
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/exactly_100_chars_password
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/unicode_password
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/special_chars_password
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/same_passwords
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/7_chars_password
=== RUN   TestUserService_BusinessLogicScenarios/password_change_edge_cases/101_chars_password
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/minimal_valid_data
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/long_email
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/unicode_in_email_domain
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/special_chars_in_email
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/whitespace_in_email
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/empty_string_email
=== RUN   TestUserService_BusinessLogicScenarios/registration_edge_cases/whitespace_only_email
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/nil_preferences
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/empty_preferences
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/simple_preferences
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/nested_preferences
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/array_preferences
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/complex_nested_structure
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/unicode_values
=== RUN   TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/numeric_values
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/valid_empty_object
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/valid_json
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_missing_quote
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_trailing_comma
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_unmatched_brace
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_array_as_root
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_string_as_root
=== RUN   TestUserService_BusinessLogicScenarios/invalid_json_handling/completely_invalid_data
--- PASS: TestUserService_BusinessLogicScenarios (0.00s)
    --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/single_character_names (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/names_with_spaces (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/names_with_special_chars (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/unicode_names (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/long_timezone (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/whitespace_only_first_name (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/profile_request_edge_cases/whitespace_only_last_name (0.00s)
    --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/exactly_8_chars_password (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/exactly_100_chars_password (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/unicode_password (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/special_chars_password (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/same_passwords (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/7_chars_password (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/password_change_edge_cases/101_chars_password (0.00s)
    --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/minimal_valid_data (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/long_email (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/unicode_in_email_domain (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/special_chars_in_email (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/whitespace_in_email (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/empty_string_email (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/registration_edge_cases/whitespace_only_email (0.00s)
    --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/nil_preferences (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/empty_preferences (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/simple_preferences (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/nested_preferences (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/array_preferences (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/complex_nested_structure (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/unicode_values (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/preferences_serialization_scenarios/numeric_values (0.00s)
    --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/valid_empty_object (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/valid_json (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_missing_quote (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_trailing_comma (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_unmatched_brace (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_array_as_root (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/invalid_json_string_as_root (0.00s)
        --- PASS: TestUserService_BusinessLogicScenarios/invalid_json_handling/completely_invalid_data (0.00s)
=== RUN   TestUserService_ValidationBoundaryConditions
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_0
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_1
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_7
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_8
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_50
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_100
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_101
=== RUN   TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_200
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases/shortest_valid_email
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_numbers
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_dots_in_local_part
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_plus_in_local_part
=== RUN   TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_hyphen_in_domain
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases/utc
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases/short_timezone
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases/long_timezone
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases/timezone_with_slash
=== RUN   TestUserService_ValidationBoundaryConditions/timezone_edge_cases/timezone_with_underscore
--- PASS: TestUserService_ValidationBoundaryConditions (0.00s)
    --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_0 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_1 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_7 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_8 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_50 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_100 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_101 (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/password_length_boundaries/password_length_200 (0.00s)
    --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases/shortest_valid_email (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_numbers (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_dots_in_local_part (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_plus_in_local_part (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/email_edge_cases/email_with_hyphen_in_domain (0.00s)
    --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases/utc (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases/short_timezone (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases/long_timezone (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases/timezone_with_slash (0.00s)
        --- PASS: TestUserService_ValidationBoundaryConditions/timezone_edge_cases/timezone_with_underscore (0.00s)
=== RUN   TestUserService_ValidationMethods
=== RUN   TestUserService_ValidationMethods/validateProfileRequest
=== RUN   TestUserService_ValidationMethods/validateProfileRequest/valid_request
=== RUN   TestUserService_ValidationMethods/validateProfileRequest/missing_first_name
=== RUN   TestUserService_ValidationMethods/validateProfileRequest/missing_last_name
=== RUN   TestUserService_ValidationMethods/validateProfileRequest/missing_timezone
=== RUN   TestUserService_ValidationMethods/validatePasswordChangeRequest
=== RUN   TestUserService_ValidationMethods/validatePasswordChangeRequest/valid_request
=== RUN   TestUserService_ValidationMethods/validatePasswordChangeRequest/missing_current_password
=== RUN   TestUserService_ValidationMethods/validatePasswordChangeRequest/missing_new_password
=== RUN   TestUserService_ValidationMethods/validatePasswordChangeRequest/new_password_too_short
=== RUN   TestUserService_ValidationMethods/validateRegistrationRequest
=== RUN   TestUserService_ValidationMethods/validateRegistrationRequest/valid_request
=== RUN   TestUserService_ValidationMethods/validateRegistrationRequest/missing_email
=== RUN   TestUserService_ValidationMethods/validateRegistrationRequest/password_too_short
--- PASS: TestUserService_ValidationMethods (0.00s)
    --- PASS: TestUserService_ValidationMethods/validateProfileRequest (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateProfileRequest/valid_request (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateProfileRequest/missing_first_name (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateProfileRequest/missing_last_name (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateProfileRequest/missing_timezone (0.00s)
    --- PASS: TestUserService_ValidationMethods/validatePasswordChangeRequest (0.00s)
        --- PASS: TestUserService_ValidationMethods/validatePasswordChangeRequest/valid_request (0.00s)
        --- PASS: TestUserService_ValidationMethods/validatePasswordChangeRequest/missing_current_password (0.00s)
        --- PASS: TestUserService_ValidationMethods/validatePasswordChangeRequest/missing_new_password (0.00s)
        --- PASS: TestUserService_ValidationMethods/validatePasswordChangeRequest/new_password_too_short (0.00s)
    --- PASS: TestUserService_ValidationMethods/validateRegistrationRequest (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateRegistrationRequest/valid_request (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateRegistrationRequest/missing_email (0.00s)
        --- PASS: TestUserService_ValidationMethods/validateRegistrationRequest/password_too_short (0.00s)
=== RUN   TestUserService_PreferencesConversion
=== RUN   TestUserService_PreferencesConversion/ConvertToJSON
=== RUN   TestUserService_PreferencesConversion/ConvertFromJSON
=== RUN   TestUserService_PreferencesConversion/ConvertToBytes
=== RUN   TestUserService_PreferencesConversion/ConvertFromBytes
--- PASS: TestUserService_PreferencesConversion (0.00s)
    --- PASS: TestUserService_PreferencesConversion/ConvertToJSON (0.00s)
    --- PASS: TestUserService_PreferencesConversion/ConvertFromJSON (0.00s)
    --- PASS: TestUserService_PreferencesConversion/ConvertToBytes (0.00s)
    --- PASS: TestUserService_PreferencesConversion/ConvertFromBytes (0.00s)
=== RUN   TestUserService_PreferencesEdgeCases
=== RUN   TestUserService_PreferencesEdgeCases/EmptyPreferences
=== RUN   TestUserService_PreferencesEdgeCases/NilPreferences
=== RUN   TestUserService_PreferencesEdgeCases/EmptyJSONBytes
=== RUN   TestUserService_PreferencesEdgeCases/InvalidJSON
--- PASS: TestUserService_PreferencesEdgeCases (0.00s)
    --- PASS: TestUserService_PreferencesEdgeCases/EmptyPreferences (0.00s)
    --- PASS: TestUserService_PreferencesEdgeCases/NilPreferences (0.00s)
    --- PASS: TestUserService_PreferencesEdgeCases/EmptyJSONBytes (0.00s)
    --- PASS: TestUserService_PreferencesEdgeCases/InvalidJSON (0.00s)
PASS
ok      github.com/garnizeh/englog/internal/services    0.023s
?       github.com/garnizeh/englog/internal/sqlc        [no test files]
?       github.com/garnizeh/englog/internal/sqlc/schema/seed    [no test files]
?       github.com/garnizeh/englog/internal/store       [no test files]
?       github.com/garnizeh/englog/internal/store/testutils     [no test files]
?       github.com/garnizeh/englog/internal/worker      [no test files]
FAIL
make: *** [Makefile:72: test-unit] Error 1
user@userland:/media/code/code/Go/garnizeh/englog$ ^C
user@userland:/media/code/code/Go/garnizeh/englog$