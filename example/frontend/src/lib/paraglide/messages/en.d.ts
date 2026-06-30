/** @typedef {import('../runtime.js').LocalizedString} LocalizedString */
/** @typedef {{}} Common_LanguageInputs */
/** @typedef {{}} Common_EnglishInputs */
/** @typedef {{}} Common_RomanianInputs */
/** @typedef {{}} Common_CancelInputs */
/** @typedef {{}} Common_DeleteInputs */
/** @typedef {{}} Common_RetryInputs */
/** @typedef {{}} Common_PreviousInputs */
/** @typedef {{}} Common_NextInputs */
/** @typedef {{}} Common_Rows_Per_PageInputs */
/** @typedef {{}} Common_StrictInputs */
/** @typedef {{}} Common_FlexibleInputs */
/** @typedef {{}} Common_RequiredInputs */
/** @typedef {{}} Common_UnknownInputs */
/** @typedef {{}} Common_ActionsInputs */
/** @typedef {{}} Common_Toggle_ThemeInputs */
/** @typedef {{}} Header_Credits_UnavailableInputs */
/** @typedef {{ count: NonNullable<unknown> }} Header_CreditsInputs */
/** @typedef {{ message: NonNullable<unknown> }} Header_Credit_Balance_UnavailableInputs */
/** @typedef {{}} Nav_AccountInputs */
/** @typedef {{}} Nav_No_Email_AddressInputs */
/** @typedef {{}} Nav_NotificationsInputs */
/** @typedef {{}} Nav_Log_OutInputs */
/** @typedef {{}} Nav_Logout_TitleInputs */
/** @typedef {{}} Nav_Logout_DescriptionInputs */
/** @typedef {{}} Nav_Logout_FailedInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Nav_Account_LinkedInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Nav_Account_Link_ConflictInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Nav_Account_Link_DeniedInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Nav_Account_Link_Not_ConfiguredInputs */
/** @typedef {{}} Nav_Account_Link_Sign_In_AgainInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Nav_Account_Link_FailedInputs */
/** @typedef {{}} Nav_DashboardInputs */
/** @typedef {{}} Nav_SchemasInputs */
/** @typedef {{}} Nav_New_SchemaInputs */
/** @typedef {{}} Nav_Edit_SchemaInputs */
/** @typedef {{}} Nav_JobsInputs */
/** @typedef {{}} Nav_New_JobInputs */
/** @typedef {{}} Nav_BillingInputs */
/** @typedef {{}} Nav_Billing_OrdersInputs */
/** @typedef {{}} Nav_Credit_Usage_HistoryInputs */
/** @typedef {{}} Nav_Developer_SettingsInputs */
/** @typedef {{}} Nav_Get_HelpInputs */
/** @typedef {{}} Nav_Quick_OcrInputs */
/** @typedef {{}} Nav_Create_Quick_Ocr_JobInputs */
/** @typedef {{}} Nav_Create_SchemaInputs */
/** @typedef {{}} Nav_Create_JobInputs */
/** @typedef {{}} Dashboard_Metric_Documents_ProcessedInputs */
/** @typedef {{}} Dashboard_Page_DescriptionInputs */
/** @typedef {{}} Dashboard_RefreshingInputs */
/** @typedef {{}} Dashboard_Loading_TitleInputs */
/** @typedef {{}} Dashboard_Loading_DescriptionInputs */
/** @typedef {{}} Dashboard_Warning_TitleInputs */
/** @typedef {{}} Dashboard_Unavailable_TitleInputs */
/** @typedef {{}} Dashboard_Unavailable_DefaultInputs */
/** @typedef {{}} Dashboard_Metric_Pages_ProcessedInputs */
/** @typedef {{}} Dashboard_Metric_Completion_RateInputs */
/** @typedef {{}} Dashboard_Metric_Credits_SpentInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Jobs_In_Progress_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Jobs_In_Progress_OtherInputs */
/** @typedef {{}} Dashboard_Pages_CompletedInputs */
/** @typedef {{ completed: NonNullable<unknown>, failed: NonNullable<unknown> }} Dashboard_Completion_SummaryInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Credits_Available_ShortInputs */
/** @typedef {{}} Dashboard_Metrics_AriaInputs */
/** @typedef {{}} Dashboard_Documents_Processed_TitleInputs */
/** @typedef {{}} Dashboard_Chart_Documents_LabelInputs */
/** @typedef {{}} Dashboard_Select_RangeInputs */
/** @typedef {{}} Dashboard_Range_7dInputs */
/** @typedef {{}} Dashboard_Range_30dInputs */
/** @typedef {{}} Dashboard_Range_90dInputs */
/** @typedef {{}} Dashboard_Recent_Documents_TitleInputs */
/** @typedef {{}} Dashboard_Recent_Documents_DescriptionInputs */
/** @typedef {{}} Dashboard_ViewInputs */
/** @typedef {{}} Dashboard_No_Saved_SchemaInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Pages_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Pages_OtherInputs */
/** @typedef {{}} Dashboard_No_Completed_DocumentsInputs */
/** @typedef {{}} Dashboard_Schema_Throughput_TitleInputs */
/** @typedef {{}} Dashboard_Schema_Throughput_DescriptionInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Documents_Processed_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Documents_Processed_OtherInputs */
/** @typedef {{}} Dashboard_No_Schema_ThroughputInputs */
/** @typedef {{}} Dashboard_Datasets_TitleInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Total_Datasets_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Total_Datasets_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Fields_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Fields_OtherInputs */
/** @typedef {{}} Dashboard_No_DatasetsInputs */
/** @typedef {{}} Dashboard_Credits_TitleInputs */
/** @typedef {{}} Dashboard_Credits_DescriptionInputs */
/** @typedef {{}} Dashboard_Low_CreditInputs */
/** @typedef {{}} Dashboard_Available_CreditsInputs */
/** @typedef {{}} Dashboard_Credits_Spent_In_RangeInputs */
/** @typedef {{}} Dashboard_BillingInputs */
/** @typedef {{}} Dashboard_Onboarding_TitleInputs */
/** @typedef {{}} Dashboard_Onboarding_DescriptionInputs */
/** @typedef {{}} Dashboard_New_Ocr_JobInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Credits_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Credits_OtherInputs */
/** @typedef {{}} Dashboard_Step_SchemaInputs */
/** @typedef {{}} Dashboard_Step_Ocr_JobInputs */
/** @typedef {{}} Dashboard_Step_DatasetInputs */
/** @typedef {{}} Dashboard_Step_Api_KeyInputs */
/** @typedef {{}} Dashboard_Step_WebhookInputs */
/** @typedef {{}} Dashboard_Step_ReadyInputs */
/** @typedef {{}} Dashboard_Step_OpenInputs */
/** @typedef {{}} Admin_Nav_UsersInputs */
/** @typedef {{}} Admin_Nav_UserInputs */
/** @typedef {{}} Admin_Nav_InvoicesInputs */
/** @typedef {{}} Admin_Nav_OrdersInputs */
/** @typedef {{}} Admin_Nav_Json_RecipesInputs */
/** @typedef {{}} Admin_Nav_AdminInputs */
/** @typedef {{}} Admin_User_FallbackInputs */
/** @typedef {{}} Sidebar_SyncraInputs */
/** @typedef {{}} Sidebar_Syncra_AdminInputs */
/** @typedef {{}} Sidebar_User_SpaceInputs */
/** @typedef {{}} Sidebar_Admin_PortalInputs */
/** @typedef {{}} Sidebar_Switch_SpaceInputs */
/** @typedef {{}} Schemas_New_TitleInputs */
/** @typedef {{}} Schemas_LibraryInputs */
/** @typedef {{}} Schemas_New_DescriptionInputs */
/** @typedef {{}} Schemas_Edit_TitleInputs */
/** @typedef {{}} Schemas_Edit_DescriptionInputs */
/** @typedef {{}} Schemas_Save_SchemaInputs */
/** @typedef {{}} Schemas_Save_ChangesInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Saved_SuccessInputs */
/** @typedef {{ name: NonNullable<unknown>, id: NonNullable<unknown> }} Schemas_Saved_Success_With_IdInputs */
/** @typedef {{ name: NonNullable<unknown>, id: NonNullable<unknown> }} Schemas_Saved_FeedbackInputs */
/** @typedef {{}} Schemas_Empty_Schema_ErrorInputs */
/** @typedef {{}} Schemas_Delete_Single_TitleInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Delete_Single_DescriptionInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Delete_Bulk_Title_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Delete_Bulk_Title_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Delete_Bulk_Description_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Delete_Bulk_Description_OtherInputs */
/** @typedef {{}} Schemas_Select_All_On_PageInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Select_SchemaInputs */
/** @typedef {{}} Schemas_Name_ColumnInputs */
/** @typedef {{}} Schemas_Id_ColumnInputs */
/** @typedef {{}} Schemas_Id_LabelInputs */
/** @typedef {{}} Schemas_Copy_IdInputs */
/** @typedef {{ id: NonNullable<unknown> }} Schemas_Copy_Id_AriaInputs */
/** @typedef {{}} Schemas_Copy_Id_SuccessInputs */
/** @typedef {{}} Schemas_Copy_Id_ErrorInputs */
/** @typedef {{}} Schemas_Strict_Mode_ColumnInputs */
/** @typedef {{}} Schemas_Created_ColumnInputs */
/** @typedef {{}} Schemas_Updated_ColumnInputs */
/** @typedef {{}} Schemas_New_SchemaInputs */
/** @typedef {{}} Schemas_No_Schemas_FoundInputs */
/** @typedef {{}} Schemas_Empty_BodyInputs */
/** @typedef {{}} Schemas_Create_SchemaInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Showing_Schemas_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Showing_Schemas_OtherInputs */
/** @typedef {{}} Schemas_No_Schemas_To_ShowInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Selected_Count_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Schemas_Selected_Count_OtherInputs */
/** @typedef {{}} Schemas_DeletingInputs */
/** @typedef {{}} Schemas_No_DescriptionInputs */
/** @typedef {{}} Schemas_Sort_Created_AscendingInputs */
/** @typedef {{}} Schemas_Sort_Created_DescendingInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Edit_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Create_Job_WithInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Clone_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Schemas_Delete_AriaInputs */
/** @typedef {{}} Schemas_Loading_SchemaInputs */
/** @typedef {{}} Schemas_Not_Found_TitleInputs */
/** @typedef {{}} Schemas_Not_Found_BodyInputs */
/** @typedef {{}} Schemas_View_SchemasInputs */
/** @typedef {{}} Schemas_Could_Not_LoadInputs */
/** @typedef {{}} Schemas_Editor_BadgeInputs */
/** @typedef {{}} Schemas_General_SettingsInputs */
/** @typedef {{}} Schemas_Schema_Name_LabelInputs */
/** @typedef {{}} Schemas_Schema_Name_PlaceholderInputs */
/** @typedef {{}} Schemas_Description_LabelInputs */
/** @typedef {{}} Schemas_Description_PlaceholderInputs */
/** @typedef {{}} Schemas_Strict_ModeInputs */
/** @typedef {{}} Schemas_Flexible_ModeInputs */
/** @typedef {{}} Schemas_Strict_Mode_DescriptionInputs */
/** @typedef {{}} Schemas_Structure_DesignerInputs */
/** @typedef {{}} Schemas_Visual_Node_DesignerInputs */
/** @typedef {{}} Schemas_Validation_Name_RequiredInputs */
/** @typedef {{}} Schemas_Validation_Name_Too_LongInputs */
/** @typedef {{}} Schemas_Validation_Schema_ObjectInputs */
/** @typedef {{}} Schemas_CloneInputs */
/** @typedef {{}} Schemas_CloningInputs */
/** @typedef {{}} Schemas_SavingInputs */
/** @typedef {{}} Json_Recipes_TitleInputs */
/** @typedef {{}} Json_Recipes_DescriptionInputs */
/** @typedef {{}} Json_Recipes_New_RecipeInputs */
/** @typedef {{}} Json_Recipes_No_Recipes_FoundInputs */
/** @typedef {{}} Json_Recipes_Empty_BodyInputs */
/** @typedef {{}} Json_Recipes_LoadingInputs */
/** @typedef {{}} Json_Recipes_Loading_RecipeInputs */
/** @typedef {{}} Json_Recipes_Counter_ColumnInputs */
/** @typedef {{}} Json_Recipes_Created_ColumnInputs */
/** @typedef {{}} Json_Recipes_Updated_ColumnInputs */
/** @typedef {{}} Json_Recipes_Json_Fields_ColumnInputs */
/** @typedef {{}} Json_Recipes_Sort_Created_AscendingInputs */
/** @typedef {{}} Json_Recipes_Sort_Created_DescendingInputs */
/** @typedef {{ count: NonNullable<unknown> }} Json_Recipes_Showing_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Json_Recipes_Showing_OtherInputs */
/** @typedef {{}} Json_Recipes_No_Recipes_To_ShowInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipes_Edit_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipes_Delete_AriaInputs */
/** @typedef {{}} Json_Recipes_New_TitleInputs */
/** @typedef {{}} Json_Recipes_New_DescriptionInputs */
/** @typedef {{}} Json_Recipes_Edit_TitleInputs */
/** @typedef {{}} Json_Recipes_Edit_DescriptionInputs */
/** @typedef {{}} Json_Recipes_Save_RecipeInputs */
/** @typedef {{}} Json_Recipes_Save_ChangesInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipes_Created_SuccessInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipes_Saved_SuccessInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipes_Deleted_SuccessInputs */
/** @typedef {{}} Json_Recipes_Delete_ConfirmInputs */
/** @typedef {{}} Json_Recipes_Not_Found_TitleInputs */
/** @typedef {{}} Json_Recipes_Not_Found_BodyInputs */
/** @typedef {{}} Json_Recipes_View_RecipesInputs */
/** @typedef {{}} Json_Recipes_Could_Not_LoadInputs */
/** @typedef {{}} Json_Recipes_Editor_BadgeInputs */
/** @typedef {{}} Json_Recipes_General_SettingsInputs */
/** @typedef {{}} Json_Recipes_Title_LabelInputs */
/** @typedef {{}} Json_Recipes_Title_PlaceholderInputs */
/** @typedef {{}} Json_Recipes_Description_LabelInputs */
/** @typedef {{}} Json_Recipes_Description_PlaceholderInputs */
/** @typedef {{}} Json_Recipes_Structure_DesignerInputs */
/** @typedef {{}} Json_Recipes_Visual_Node_DesignerInputs */
/** @typedef {{}} Json_Recipes_Category_LabelInputs */
/** @typedef {{}} Json_Recipes_OthersInputs */
/** @typedef {{}} Json_Recipes_Manage_CategoriesInputs */
/** @typedef {{}} Json_Recipes_Validation_Title_RequiredInputs */
/** @typedef {{}} Json_Recipes_Validation_Title_Too_LongInputs */
/** @typedef {{}} Json_Recipes_Validation_Json_ObjectInputs */
/** @typedef {{}} Json_Recipes_SavingInputs */
/** @typedef {{}} Json_Recipes_DeletingInputs */
/** @typedef {{}} Json_Recipe_Categories_TitleInputs */
/** @typedef {{}} Json_Recipe_Categories_DescriptionInputs */
/** @typedef {{}} Json_Recipe_Categories_Title_En_LabelInputs */
/** @typedef {{}} Json_Recipe_Categories_Title_Ro_LabelInputs */
/** @typedef {{}} Json_Recipe_Categories_Create_CategoryInputs */
/** @typedef {{}} Json_Recipe_Categories_Save_CategoryInputs */
/** @typedef {{}} Json_Recipe_Categories_Edit_TitleInputs */
/** @typedef {{}} Json_Recipe_Categories_Delete_ConfirmInputs */
/** @typedef {{}} Json_Recipe_Categories_LoadingInputs */
/** @typedef {{}} Json_Recipe_Categories_Could_Not_LoadInputs */
/** @typedef {{}} Json_Recipe_Categories_Empty_TitleInputs */
/** @typedef {{}} Json_Recipe_Categories_Empty_BodyInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipe_Categories_Created_SuccessInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipe_Categories_Saved_SuccessInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipe_Categories_Deleted_SuccessInputs */
/** @typedef {{}} Json_Recipe_Categories_Validation_Titles_RequiredInputs */
/** @typedef {{}} Json_Recipe_Categories_Validation_Titles_Too_LongInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipe_Categories_Edit_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Json_Recipe_Categories_Delete_AriaInputs */
/** @typedef {{}} Ocr_Recipes_NavInputs */
/** @typedef {{}} Ocr_Recipes_TitleInputs */
/** @typedef {{}} Ocr_Recipes_Meta_DescriptionInputs */
/** @typedef {{}} Ocr_Recipes_EyebrowInputs */
/** @typedef {{}} Ocr_Recipes_Hero_TitleInputs */
/** @typedef {{}} Ocr_Recipes_Hero_DescriptionInputs */
/** @typedef {{}} Ocr_Recipes_Search_LabelInputs */
/** @typedef {{}} Ocr_Recipes_Search_PlaceholderInputs */
/** @typedef {{}} Ocr_Recipes_Category_FilterInputs */
/** @typedef {{}} Ocr_Recipes_All_CategoriesInputs */
/** @typedef {{}} Ocr_Recipes_Sort_LabelInputs */
/** @typedef {{}} Ocr_Recipes_Sort_PopularInputs */
/** @typedef {{}} Ocr_Recipes_Sort_NewestInputs */
/** @typedef {{}} Ocr_Recipes_Sort_AzInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Showing_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Showing_OtherInputs */
/** @typedef {{}} Ocr_Recipes_No_Matches_TitleInputs */
/** @typedef {{}} Ocr_Recipes_No_Matches_BodyInputs */
/** @typedef {{}} Ocr_Recipes_OthersInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Fields_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Fields_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Required_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Required_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Deploys_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Ocr_Recipes_Deploys_OtherInputs */
/** @typedef {{}} Ocr_Recipes_Json_FieldsInputs */
/** @typedef {{}} Ocr_Recipes_System_RecipeInputs */
/** @typedef {{}} Ocr_Recipes_Strict_SchemaInputs */
/** @typedef {{}} Ocr_Recipes_RequiredInputs */
/** @typedef {{}} Ocr_Recipes_Preview_JsonInputs */
/** @typedef {{}} Ocr_Recipes_No_FieldsInputs */
/** @typedef {{}} Ocr_Recipes_Clone_RecipeInputs */
/** @typedef {{ name: NonNullable<unknown> }} Ocr_Recipes_Clone_AriaInputs */
/** @typedef {{}} Ocr_Recipes_Log_In_To_CloneInputs */
/** @typedef {{}} Ocr_Recipes_Clone_FailedInputs */
/** @typedef {{}} Ocr_Recipes_Load_FailedInputs */
/** @typedef {{}} Jobs_Page_TitleInputs */
/** @typedef {{}} Jobs_Missing_Schema_IdInputs */
/** @typedef {{}} Jobs_Missing_Job_IdInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Delete_Bulk_Title_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Delete_Bulk_Title_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Delete_Bulk_Description_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Delete_Bulk_Description_OtherInputs */
/** @typedef {{}} Jobs_Delete_Single_TitleInputs */
/** @typedef {{ name: NonNullable<unknown> }} Jobs_Delete_Single_DescriptionInputs */
/** @typedef {{}} Jobs_Status_QueuedInputs */
/** @typedef {{}} Jobs_Status_PendingInputs */
/** @typedef {{}} Jobs_Status_ProcessingInputs */
/** @typedef {{}} Jobs_Status_CompletedInputs */
/** @typedef {{}} Jobs_Status_FailedInputs */
/** @typedef {{}} Jobs_Inline_SchemaInputs */
/** @typedef {{}} Jobs_No_SchemaInputs */
/** @typedef {{}} Jobs_SchemaInputs */
/** @typedef {{}} Jobs_Select_All_On_PageInputs */
/** @typedef {{ name: NonNullable<unknown> }} Jobs_Select_JobInputs */
/** @typedef {{}} Jobs_Filename_ColumnInputs */
/** @typedef {{}} Jobs_Status_ColumnInputs */
/** @typedef {{}} Jobs_Created_ColumnInputs */
/** @typedef {{}} Jobs_File_Size_ColumnInputs */
/** @typedef {{}} Jobs_Pages_ColumnInputs */
/** @typedef {{}} Jobs_New_JobInputs */
/** @typedef {{}} Jobs_No_Jobs_FoundInputs */
/** @typedef {{}} Jobs_Empty_BodyInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Showing_Jobs_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Showing_Jobs_OtherInputs */
/** @typedef {{}} Jobs_No_Jobs_To_ShowInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Selected_Count_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Jobs_Selected_Count_OtherInputs */
/** @typedef {{}} Jobs_DeletingInputs */
/** @typedef {{ name: NonNullable<unknown> }} Jobs_Delete_JobInputs */
/** @typedef {{}} Jobs_Saved_Extraction_SchemaInputs */
/** @typedef {{}} Jobs_Inline_Schema_DescriptionInputs */
/** @typedef {{}} Jobs_Extraction_Schema_DetailsInputs */
/** @typedef {{}} New_Job_Missing_Document_IdInputs */
/** @typedef {{}} New_Job_Failed_CreateInputs */
/** @typedef {{}} New_Job_Insufficient_Credits_BuyInputs */
/** @typedef {{}} New_Job_Failed_Load_DocumentInputs */
/** @typedef {{}} New_Job_Invalid_Document_ResponseInputs */
/** @typedef {{}} New_Job_Failed_Load_SchemasInputs */
/** @typedef {{}} New_Job_Invalid_Schema_ResponseInputs */
/** @typedef {{}} New_Job_Invalid_Job_ResponseInputs */
/** @typedef {{}} New_Job_Failed_Load_JobInputs */
/** @typedef {{}} New_Job_Failed_Poll_JobInputs */
/** @typedef {{}} New_Job_Select_SchemaInputs */
/** @typedef {{}} New_Job_Select_Schema_PlaceholderInputs */
/** @typedef {{}} New_Job_Configure_Payload_FormatInputs */
/** @typedef {{}} New_Job_Upload_DocumentsInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Files_Selected_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Files_Selected_OtherInputs */
/** @typedef {{}} New_Job_Drag_Or_Browse_FilesInputs */
/** @typedef {{}} New_Job_Run_MonitorInputs */
/** @typedef {{}} New_Job_Processing_BatchInputs */
/** @typedef {{}} New_Job_Start_Extraction_PipelineInputs */
/** @typedef {{}} New_Job_Select_Extraction_SchemaInputs */
/** @typedef {{}} New_Job_Select_Schema_DescriptionInputs */
/** @typedef {{}} New_Job_Select_Extraction_Schema_AriaInputs */
/** @typedef {{}} New_Job_Search_SchemasInputs */
/** @typedef {{}} New_Job_Loading_SchemasInputs */
/** @typedef {{}} New_Job_No_Schemas_FoundInputs */
/** @typedef {{}} New_Job_No_Schema_Ocr_OnlyInputs */
/** @typedef {{}} New_Job_No_Schema_DescriptionInputs */
/** @typedef {{}} New_Job_No_Personal_SchemasInputs */
/** @typedef {{}} New_Job_Create_OneInputs */
/** @typedef {{}} New_Job_Selected_Schema_HelpInputs */
/** @typedef {{}} New_Job_No_Schema_Selected_HelpInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Target_Mapped_FieldsInputs */
/** @typedef {{}} New_Job_No_Fields_DefinedInputs */
/** @typedef {{}} New_Job_Ocr_Only_Mode_ActiveInputs */
/** @typedef {{}} New_Job_Ocr_Only_Mode_BodyInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Upload_Documents_DescriptionInputs */
/** @typedef {{}} New_Job_Dropzone_TitleInputs */
/** @typedef {{ size: NonNullable<unknown> }} New_Job_Dropzone_DescriptionInputs */
/** @typedef {{}} New_Job_Browse_FilesInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Pending_Upload_QueueInputs */
/** @typedef {{}} New_Job_Clear_AllInputs */
/** @typedef {{}} New_Job_Remove_FileInputs */
/** @typedef {{}} New_Job_Extraction_Queue_ResultsInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_File_Count_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_File_Count_OtherInputs */
/** @typedef {{ label: NonNullable<unknown> }} New_Job_TotalInputs */
/** @typedef {{}} New_Job_Active_Batch_StatusInputs */
/** @typedef {{}} New_Job_Active_Batch_DescriptionInputs */
/** @typedef {{ progress: NonNullable<unknown> }} New_Job_ProgressInputs */
/** @typedef {{}} New_Job_Total_FilesInputs */
/** @typedef {{}} New_Job_CompletedInputs */
/** @typedef {{}} New_Job_ProcessingInputs */
/** @typedef {{}} New_Job_FailedInputs */
/** @typedef {{}} New_Job_No_Active_Extraction_JobsInputs */
/** @typedef {{}} New_Job_No_Active_Extraction_Jobs_BodyInputs */
/** @typedef {{}} New_Job_Preview_DocumentInputs */
/** @typedef {{}} New_Job_Preview_UnavailableInputs */
/** @typedef {{}} New_Job_Remove_Failed_JobInputs */
/** @typedef {{}} New_Job_Queueing_DocumentsInputs */
/** @typedef {{}} New_Job_Extracting_ContentInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Run_Extraction_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} New_Job_Run_Extraction_OtherInputs */
/** @typedef {{}} New_Job_Insufficient_Credits_DocumentInputs */
/** @typedef {{}} New_Job_Processing_FailedInputs */
/** @typedef {{}} New_Job_ProcessedInputs */
/** @typedef {{ id: NonNullable<unknown> }} New_Job_Document_IdInputs */
/** @typedef {{}} New_Job_Creating_JobInputs */
/** @typedef {{}} New_Job_Queued_ProcessingInputs */
/** @typedef {{}} New_Job_Extracting_EntitiesInputs */
/** @typedef {{}} Common_ApplyInputs */
/** @typedef {{}} Common_ClearInputs */
/** @typedef {{}} Common_SavingInputs */
/** @typedef {{}} Common_LoadingInputs */
/** @typedef {{}} Common_RefreshInputs */
/** @typedef {{}} Common_ConnectedInputs */
/** @typedef {{}} Common_ConnectInputs */
/** @typedef {{}} Common_DownloadInputs */
/** @typedef {{}} Common_TodayInputs */
/** @typedef {{}} Common_This_WeekInputs */
/** @typedef {{}} Common_This_MonthInputs */
/** @typedef {{}} Common_AnyInputs */
/** @typedef {{}} Billing_UnavailableInputs */
/** @typedef {{}} Billing_Credit_Blocks_ErrorInputs */
/** @typedef {{}} Billing_Checkout_UnavailableInputs */
/** @typedef {{}} Billing_Payment_Received_TitleInputs */
/** @typedef {{}} Billing_Payment_Received_BodyInputs */
/** @typedef {{}} Billing_Checkout_Canceled_TitleInputs */
/** @typedef {{}} Billing_Checkout_Canceled_BodyInputs */
/** @typedef {{}} Billing_Available_BalanceInputs */
/** @typedef {{}} Billing_ConversionInputs */
/** @typedef {{}} Billing_Conversion_RateInputs */
/** @typedef {{}} Billing_Balance_Checked_UploadInputs */
/** @typedef {{}} Billing_Debited_After_SuccessInputs */
/** @typedef {{}} Billing_Secure_Stripe_CheckoutInputs */
/** @typedef {{}} Billing_Purchase_CreditsInputs */
/** @typedef {{}} Billing_Credits_To_PurchaseInputs */
/** @typedef {{}} Billing_Volume_Discount_TiersInputs */
/** @typedef {{}} Billing_Total_To_PayInputs */
/** @typedef {{}} Billing_Base_PriceInputs */
/** @typedef {{}} Billing_Volume_DiscountInputs */
/** @typedef {{}} Billing_Starting_CheckoutInputs */
/** @typedef {{}} Billing_Secure_CheckoutInputs */
/** @typedef {{}} Billing_Buy_CreditsInputs */
/** @typedef {{}} Billing_Orders_Page_TitleInputs */
/** @typedef {{}} Billing_Orders_Order_Date_FilterInputs */
/** @typedef {{}} Billing_Orders_Amount_ColumnInputs */
/** @typedef {{}} Billing_Orders_Credits_ColumnInputs */
/** @typedef {{}} Billing_Orders_Status_ColumnInputs */
/** @typedef {{}} Billing_Orders_Payment_Datetime_ColumnInputs */
/** @typedef {{}} Billing_Orders_Invoice_ColumnInputs */
/** @typedef {{}} Billing_Orders_PresetsInputs */
/** @typedef {{}} Billing_Orders_Filter_StatusInputs */
/** @typedef {{}} Billing_Orders_All_OrdersInputs */
/** @typedef {{}} Billing_Orders_Clear_FiltersInputs */
/** @typedef {{}} Billing_Orders_Clear_Filters_ActionInputs */
/** @typedef {{}} Billing_Orders_No_Orders_FoundInputs */
/** @typedef {{}} Billing_Orders_No_Orders_YetInputs */
/** @typedef {{}} Billing_Orders_No_Orders_MatchInputs */
/** @typedef {{}} Billing_Orders_Empty_BodyInputs */
/** @typedef {{ count: NonNullable<unknown> }} Billing_Orders_Showing_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Billing_Orders_Showing_OtherInputs */
/** @typedef {{}} Billing_Orders_None_To_ShowInputs */
/** @typedef {{}} Billing_Orders_Sort_Order_Date_AscendingInputs */
/** @typedef {{}} Billing_Orders_Sort_Order_Date_DescendingInputs */
/** @typedef {{}} Billing_Order_Status_PendingInputs */
/** @typedef {{}} Billing_Order_Status_PaidInputs */
/** @typedef {{}} Billing_Order_Status_FailedInputs */
/** @typedef {{}} Billing_Order_Status_RefundedInputs */
/** @typedef {{}} Billing_Order_Status_CanceledInputs */
/** @typedef {{ invoice: NonNullable<unknown> }} Billing_Orders_Invoice_Pdf_TitleInputs */
/** @typedef {{ invoice: NonNullable<unknown> }} Billing_Orders_Invoice_Preview_TitleInputs */
/** @typedef {{}} Billing_Orders_Invoice_Preview_DescriptionInputs */
/** @typedef {{ invoice: NonNullable<unknown> }} Billing_Orders_Invoice_Iframe_TitleInputs */
/** @typedef {{}} Billing_Orders_Download_InvoiceInputs */
/** @typedef {{}} Credit_Usage_Page_TitleInputs */
/** @typedef {{}} Credit_Usage_Date_Range_FilterInputs */
/** @typedef {{}} Credit_Usage_Created_ColumnInputs */
/** @typedef {{}} Credit_Usage_Type_ColumnInputs */
/** @typedef {{}} Credit_Usage_Credits_ColumnInputs */
/** @typedef {{}} Credit_Usage_Related_Id_ColumnInputs */
/** @typedef {{}} Credit_Usage_Filter_TypeInputs */
/** @typedef {{}} Credit_Usage_All_ActivityInputs */
/** @typedef {{}} Credit_Usage_Type_PurchaseInputs */
/** @typedef {{}} Credit_Usage_Type_DebitInputs */
/** @typedef {{}} Credit_Usage_No_Usage_FoundInputs */
/** @typedef {{}} Credit_Usage_No_Usage_YetInputs */
/** @typedef {{}} Credit_Usage_No_Usage_MatchInputs */
/** @typedef {{}} Credit_Usage_Empty_BodyInputs */
/** @typedef {{ count: NonNullable<unknown> }} Credit_Usage_Showing_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Credit_Usage_Showing_OtherInputs */
/** @typedef {{}} Credit_Usage_None_To_ShowInputs */
/** @typedef {{}} Credit_Usage_Sort_Created_AscendingInputs */
/** @typedef {{}} Credit_Usage_Sort_Created_DescendingInputs */
/** @typedef {{}} Account_Settings_TitleInputs */
/** @typedef {{}} Account_Settings_DescriptionInputs */
/** @typedef {{}} Account_Settings_Nav_LabelInputs */
/** @typedef {{}} Account_Settings_Account_FallbackInputs */
/** @typedef {{}} Account_Settings_No_Email_AddressInputs */
/** @typedef {{}} Account_Settings_GeneralInputs */
/** @typedef {{}} Account_Settings_SecurityInputs */
/** @typedef {{}} Account_Settings_SessionsInputs */
/** @typedef {{}} Account_Settings_Linked_AccountsInputs */
/** @typedef {{}} Account_Settings_Update_ErrorInputs */
/** @typedef {{}} Account_Settings_Save_ErrorInputs */
/** @typedef {{}} Account_Settings_Revoke_Session_TitleInputs */
/** @typedef {{ session: NonNullable<unknown> }} Account_Settings_Revoke_Session_DescriptionInputs */
/** @typedef {{}} Account_Settings_RevokeInputs */
/** @typedef {{}} Account_Settings_Session_RevokedInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Account_Settings_Unlink_Provider_TitleInputs */
/** @typedef {{ provider: NonNullable<unknown> }} Account_Settings_Unlink_Provider_DescriptionInputs */
/** @typedef {{}} Account_Settings_UnlinkInputs */
/** @typedef {{}} Account_Settings_Linked_Account_RemovedInputs */
/** @typedef {{}} Account_Settings_Avatar_SavedInputs */
/** @typedef {{}} Account_Settings_Name_SavedInputs */
/** @typedef {{}} Account_Settings_Email_SavedInputs */
/** @typedef {{}} Account_Settings_Language_SavedInputs */
/** @typedef {{}} Account_Settings_Password_UpdatedInputs */
/** @typedef {{}} Account_Settings_Current_SessionInputs */
/** @typedef {{}} Account_Settings_Browser_SessionInputs */
/** @typedef {{ date: NonNullable<unknown> }} Account_Settings_Session_Created_AtInputs */
/** @typedef {{ ip: NonNullable<unknown>, date: NonNullable<unknown> }} Account_Settings_Session_Ip_Created_AtInputs */
/** @typedef {{}} Account_Settings_UnknownInputs */
/** @typedef {{}} Account_Settings_AvatarInputs */
/** @typedef {{}} Account_Settings_Avatar_DescriptionInputs */
/** @typedef {{}} Account_Settings_Avatar_UploadingInputs */
/** @typedef {{}} Account_Settings_Avatar_UploadInputs */
/** @typedef {{}} Account_Settings_Avatar_File_HintInputs */
/** @typedef {{}} Account_Settings_Crop_AvatarInputs */
/** @typedef {{}} Account_Settings_Crop_Avatar_DescriptionInputs */
/** @typedef {{}} Account_Settings_Display_NameInputs */
/** @typedef {{}} Account_Settings_Email_AddressInputs */
/** @typedef {{}} Account_Settings_LanguageInputs */
/** @typedef {{}} Account_Settings_Save_NameInputs */
/** @typedef {{}} Account_Settings_Save_EmailInputs */
/** @typedef {{}} Account_Settings_Save_LanguageInputs */
/** @typedef {{}} Account_Settings_Save_PasswordInputs */
/** @typedef {{}} Account_Settings_New_PasswordInputs */
/** @typedef {{}} Account_Settings_Confirm_PasswordInputs */
/** @typedef {{}} Account_Settings_Security_DescriptionInputs */
/** @typedef {{}} Account_Settings_Sessions_DescriptionInputs */
/** @typedef {{}} Account_Settings_Loading_SessionsInputs */
/** @typedef {{}} Account_Settings_No_SessionsInputs */
/** @typedef {{}} Account_Settings_CurrentInputs */
/** @typedef {{ date: NonNullable<unknown> }} Account_Settings_ExpiresInputs */
/** @typedef {{}} Account_Settings_Current_Session_Cannot_RevokeInputs */
/** @typedef {{}} Account_Settings_Revoke_SessionInputs */
/** @typedef {{}} Account_Settings_RevokingInputs */
/** @typedef {{}} Account_Settings_Linked_Accounts_DescriptionInputs */
/** @typedef {{}} Account_Settings_Loading_Linked_AccountsInputs */
/** @typedef {{}} Account_Settings_No_Sign_In_MethodsInputs */
/** @typedef {{}} Account_Settings_Email_PasswordInputs */
/** @typedef {{ email: NonNullable<unknown> }} Account_Settings_Password_EnabledInputs */
/** @typedef {{}} Account_Settings_Add_PasswordInputs */
/** @typedef {{}} Account_Settings_Set_PasswordInputs */
/** @typedef {{}} Account_Settings_Provider_Google_DescriptionInputs */
/** @typedef {{}} Account_Settings_Provider_Github_DescriptionInputs */
/** @typedef {{ date: NonNullable<unknown> }} Account_Settings_Linked_AtInputs */
/** @typedef {{}} Account_Settings_UnlinkingInputs */
/** @typedef {{}} Account_Settings_Unavailable_TitleInputs */
/** @typedef {{}} Account_Settings_Unavailable_BodyInputs */
/** @typedef {{}} Billing_Profile_TitleInputs */
/** @typedef {{}} Billing_Profile_DescriptionInputs */
/** @typedef {{}} Billing_Profile_Load_ErrorInputs */
/** @typedef {{}} Billing_Profile_Save_ErrorInputs */
/** @typedef {{}} Billing_Profile_SavedInputs */
/** @typedef {{}} Billing_Profile_Company_NameInputs */
/** @typedef {{}} Billing_Profile_Full_NameInputs */
/** @typedef {{}} Billing_Profile_Error_TitleInputs */
/** @typedef {{}} Billing_Profile_LoadingInputs */
/** @typedef {{}} Billing_Profile_Loading_BodyInputs */
/** @typedef {{}} Billing_Profile_Failed_LoadInputs */
/** @typedef {{}} Billing_Profile_Retry_LoadingInputs */
/** @typedef {{}} Billing_Profile_Billing_EntityInputs */
/** @typedef {{}} Billing_Profile_Entity_DescriptionInputs */
/** @typedef {{}} Billing_Profile_IndividualInputs */
/** @typedef {{}} Billing_Profile_CompanyInputs */
/** @typedef {{}} Billing_Profile_General_DetailsInputs */
/** @typedef {{}} Billing_Profile_Billing_EmailInputs */
/** @typedef {{}} Billing_Profile_Billing_AddressInputs */
/** @typedef {{}} Billing_Profile_Address_Line1Inputs */
/** @typedef {{}} Billing_Profile_Address_Line2Inputs */
/** @typedef {{}} Billing_Profile_CityInputs */
/** @typedef {{}} Billing_Profile_Region_StateInputs */
/** @typedef {{}} Billing_Profile_CountryInputs */
/** @typedef {{}} Billing_Profile_Postal_CodeInputs */
/** @typedef {{}} Billing_Profile_Company_DetailsInputs */
/** @typedef {{}} Billing_Profile_Fiscal_CodeInputs */
/** @typedef {{}} Billing_Profile_Registration_NumberInputs */
/** @typedef {{}} Billing_Profile_Save_ButtonInputs */
/** @typedef {{}} Datasets_Page_TitleInputs */
/** @typedef {{}} Datasets_Detail_Page_TitleInputs */
/** @typedef {{}} Datasets_Name_ColumnInputs */
/** @typedef {{}} Datasets_Schema_ColumnInputs */
/** @typedef {{}} Datasets_Fields_ColumnInputs */
/** @typedef {{}} Datasets_Created_ColumnInputs */
/** @typedef {{}} Datasets_Actions_ColumnInputs */
/** @typedef {{}} Datasets_Sort_Created_AscendingInputs */
/** @typedef {{}} Datasets_Sort_Created_DescendingInputs */
/** @typedef {{}} Datasets_RetryInputs */
/** @typedef {{}} Datasets_OpenInputs */
/** @typedef {{}} Datasets_No_Datasets_FoundInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Showing_Datasets_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Showing_Datasets_OtherInputs */
/** @typedef {{}} Datasets_No_Datasets_To_ShowInputs */
/** @typedef {{}} Datasets_Rows_Per_PageInputs */
/** @typedef {{}} Datasets_Previous_PageInputs */
/** @typedef {{}} Datasets_Next_PageInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Field_Count_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Field_Count_OtherInputs */
/** @typedef {{}} Datasets_Date_RangeInputs */
/** @typedef {{}} Datasets_Any_DateInputs */
/** @typedef {{ start: NonNullable<unknown>, end: NonNullable<unknown> }} Datasets_Date_Range_ValueInputs */
/** @typedef {{}} Datasets_PresetsInputs */
/** @typedef {{}} Datasets_TodayInputs */
/** @typedef {{}} Datasets_This_WeekInputs */
/** @typedef {{}} Datasets_This_MonthInputs */
/** @typedef {{}} Datasets_ClearInputs */
/** @typedef {{}} Datasets_ApplyInputs */
/** @typedef {{}} Datasets_Document_Id_ColumnInputs */
/** @typedef {{}} Datasets_Filename_ColumnInputs */
/** @typedef {{}} Datasets_Not_Found_TitleInputs */
/** @typedef {{}} Datasets_Not_Found_BodyInputs */
/** @typedef {{}} Datasets_View_DatasetsInputs */
/** @typedef {{ documentId: NonNullable<unknown> }} Datasets_Preview_DocumentInputs */
/** @typedef {{}} Datasets_No_Documents_ExtractedInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Showing_Rows_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Showing_Rows_OtherInputs */
/** @typedef {{}} Datasets_No_Rows_To_ShowInputs */
/** @typedef {{}} Datasets_Export_CsvInputs */
/** @typedef {{}} Datasets_Export_XlsxInputs */
/** @typedef {{}} Datasets_Failed_ExportInputs */
/** @typedef {{}} Datasets_Invalid_DateInputs */
/** @typedef {{}} Datasets_Missing_Document_IdInputs */
/** @typedef {{}} Datasets_Add_DatasetInputs */
/** @typedef {{}} Datasets_All_DatasetsInputs */
/** @typedef {{}} Datasets_Retry_DatasetsInputs */
/** @typedef {{}} Datasets_No_DatasetsInputs */
/** @typedef {{}} Datasets_Dataset_ActionsInputs */
/** @typedef {{}} Datasets_EditInputs */
/** @typedef {{}} Datasets_DeleteInputs */
/** @typedef {{}} Datasets_Delete_FailedInputs */
/** @typedef {{}} Datasets_Delete_Confirm_TitleInputs */
/** @typedef {{ name: NonNullable<unknown> }} Datasets_Delete_Confirm_DescriptionInputs */
/** @typedef {{}} Datasets_Dialog_Title_NewInputs */
/** @typedef {{}} Datasets_Dialog_Title_EditInputs */
/** @typedef {{}} Datasets_Save_ChangesInputs */
/** @typedef {{}} Datasets_Create_DatasetInputs */
/** @typedef {{}} Datasets_Selected_SchemaInputs */
/** @typedef {{}} Datasets_Loading_SchemasInputs */
/** @typedef {{}} Datasets_Select_SchemaInputs */
/** @typedef {{}} Datasets_No_Fields_SelectedInputs */
/** @typedef {{}} Datasets_One_Field_SelectedInputs */
/** @typedef {{ count: NonNullable<unknown> }} Datasets_Fields_SelectedInputs */
/** @typedef {{ label: NonNullable<unknown> }} Datasets_Collapse_FieldInputs */
/** @typedef {{ label: NonNullable<unknown> }} Datasets_Expand_FieldInputs */
/** @typedef {{ label: NonNullable<unknown> }} Datasets_Select_FieldInputs */
/** @typedef {{}} Datasets_Name_PlaceholderInputs */
/** @typedef {{}} Datasets_Search_SchemasInputs */
/** @typedef {{}} Datasets_No_Schemas_FoundInputs */
/** @typedef {{}} Datasets_No_FieldsInputs */
/** @typedef {{}} Datasets_CancelInputs */
/** @typedef {{}} Datasets_Json_BadgeInputs */
/** @typedef {{}} Documents_Page_TitleInputs */
/** @typedef {{}} Documents_New_Ocr_JobInputs */
/** @typedef {{}} Documents_Search_Filename_PlaceholderInputs */
/** @typedef {{}} Documents_Search_FilenameInputs */
/** @typedef {{}} Documents_Date_RangeInputs */
/** @typedef {{}} Documents_Any_DateInputs */
/** @typedef {{ start: NonNullable<unknown>, end: NonNullable<unknown> }} Documents_Date_Range_ValueInputs */
/** @typedef {{}} Documents_PresetsInputs */
/** @typedef {{}} Documents_TodayInputs */
/** @typedef {{}} Documents_This_WeekInputs */
/** @typedef {{}} Documents_This_MonthInputs */
/** @typedef {{}} Documents_ClearInputs */
/** @typedef {{}} Documents_ApplyInputs */
/** @typedef {{}} Documents_Filter_By_CollectionInputs */
/** @typedef {{}} Documents_Filter_By_SchemaInputs */
/** @typedef {{}} Documents_Unknown_CollectionInputs */
/** @typedef {{}} Documents_All_CollectionsInputs */
/** @typedef {{}} Documents_All_SchemasInputs */
/** @typedef {{}} Documents_Missing_Document_IdInputs */
/** @typedef {{}} Documents_Failed_Load_DocumentsInputs */
/** @typedef {{}} Documents_Failed_Load_DocumentInputs */
/** @typedef {{}} Documents_Failed_Delete_DocumentInputs */
/** @typedef {{}} Documents_Failed_Update_DocumentInputs */
/** @typedef {{}} Documents_Failed_Delete_DocumentsInputs */
/** @typedef {{}} Documents_Failed_Move_DocumentsInputs */
/** @typedef {{}} Documents_Failed_DownloadInputs */
/** @typedef {{}} Documents_Invalid_DateInputs */
/** @typedef {{}} Documents_Select_All_On_PageInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Select_DocumentInputs */
/** @typedef {{}} Documents_Filename_ColumnInputs */
/** @typedef {{}} Documents_Collections_ColumnInputs */
/** @typedef {{}} Documents_Pages_ColumnInputs */
/** @typedef {{}} Documents_Created_ColumnInputs */
/** @typedef {{}} Documents_File_Size_ColumnInputs */
/** @typedef {{}} Documents_Sort_Created_AscendingInputs */
/** @typedef {{}} Documents_Sort_Created_DescendingInputs */
/** @typedef {{}} Documents_Collection_Not_Found_TitleInputs */
/** @typedef {{}} Documents_Collection_Not_Found_BodyInputs */
/** @typedef {{}} Documents_View_All_DocumentsInputs */
/** @typedef {{}} Documents_RetryInputs */
/** @typedef {{}} Documents_No_Documents_FoundInputs */
/** @typedef {{}} Documents_Empty_Filtered_BodyInputs */
/** @typedef {{}} Documents_Empty_Unfiltered_BodyInputs */
/** @typedef {{}} Documents_Clear_FiltersInputs */
/** @typedef {{}} Documents_Process_First_DocumentInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Showing_Documents_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Showing_Documents_OtherInputs */
/** @typedef {{}} Documents_No_Documents_To_ShowInputs */
/** @typedef {{}} Documents_Rows_Per_PageInputs */
/** @typedef {{}} Documents_PreviousInputs */
/** @typedef {{}} Documents_NextInputs */
/** @typedef {{}} Documents_DeleteInputs */
/** @typedef {{}} Documents_Delete_Single_TitleInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Delete_Single_DescriptionInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Delete_Bulk_Title_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Delete_Bulk_Title_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Delete_Bulk_Description_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Delete_Bulk_Description_OtherInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Selected_Count_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Selected_Count_OtherInputs */
/** @typedef {{}} Documents_Download_SelectedInputs */
/** @typedef {{}} Documents_DownloadInputs */
/** @typedef {{}} Documents_DownloadingInputs */
/** @typedef {{}} Documents_MoveInputs */
/** @typedef {{}} Documents_MovingInputs */
/** @typedef {{}} Documents_DeletingInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Open_Actions_ForInputs */
/** @typedef {{}} Documents_PreviewInputs */
/** @typedef {{}} Documents_RenameInputs */
/** @typedef {{}} Documents_Failed_RenameInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Rename_FileInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Preview_FileInputs */
/** @typedef {{}} Documents_Download_Dialog_Title_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Download_Dialog_Title_OtherInputs */
/** @typedef {{}} Documents_Selected_DocumentsInputs */
/** @typedef {{}} Documents_Format_MarkdownInputs */
/** @typedef {{}} Documents_Format_HtmlInputs */
/** @typedef {{}} Documents_Format_JsonInputs */
/** @typedef {{}} Documents_Preparing_DownloadInputs */
/** @typedef {{}} Documents_No_Collections_SelectedInputs */
/** @typedef {{}} Documents_One_Collection_SelectedInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Collections_SelectedInputs */
/** @typedef {{}} Documents_Remove_From_AllInputs */
/** @typedef {{}} Documents_Move_DocumentsInputs */
/** @typedef {{}} Documents_Move_Description_OneInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Move_Description_OtherInputs */
/** @typedef {{}} Documents_Collections_LabelInputs */
/** @typedef {{}} Documents_Search_CollectionsInputs */
/** @typedef {{}} Documents_Loading_CollectionsInputs */
/** @typedef {{}} Documents_No_Collections_FoundInputs */
/** @typedef {{}} Documents_CancelInputs */
/** @typedef {{}} Documents_Collections_Nav_LabelInputs */
/** @typedef {{}} Documents_Add_CollectionInputs */
/** @typedef {{}} Documents_All_DocumentsInputs */
/** @typedef {{}} Documents_Retry_CollectionsInputs */
/** @typedef {{}} Documents_No_CollectionsInputs */
/** @typedef {{}} Documents_Collection_ActionsInputs */
/** @typedef {{}} Documents_EditInputs */
/** @typedef {{}} Documents_Delete_FailedInputs */
/** @typedef {{}} Documents_Delete_Collection_TitleInputs */
/** @typedef {{ name: NonNullable<unknown> }} Documents_Delete_Collection_DescriptionInputs */
/** @typedef {{}} Documents_Collection_Dialog_Title_NewInputs */
/** @typedef {{}} Documents_Collection_Dialog_Title_EditInputs */
/** @typedef {{}} Documents_Collection_Dialog_Description_NewInputs */
/** @typedef {{}} Documents_Collection_Dialog_Description_EditInputs */
/** @typedef {{}} Documents_Save_ChangesInputs */
/** @typedef {{}} Documents_Create_CollectionInputs */
/** @typedef {{}} Documents_Name_ColumnInputs */
/** @typedef {{}} Documents_Collection_Name_PlaceholderInputs */
/** @typedef {{}} Documents_Schemas_LabelInputs */
/** @typedef {{}} Documents_No_Schemas_SelectedInputs */
/** @typedef {{}} Documents_One_Schema_SelectedInputs */
/** @typedef {{ count: NonNullable<unknown> }} Documents_Schemas_SelectedInputs */
/** @typedef {{}} Documents_Search_SchemasInputs */
/** @typedef {{}} Documents_Loading_SchemasInputs */
/** @typedef {{}} Documents_No_Schemas_FoundInputs */
/** @typedef {{}} Documents_Collection_Schema_HintInputs */
/** @typedef {{}} Documents_Preview_Fallback_TitleInputs */
/** @typedef {{}} Documents_Preview_DescriptionInputs */
/** @typedef {{}} Documents_Rename_Document_TitleInputs */
/** @typedef {{}} Documents_Loading_DocumentInputs */
/** @typedef {{}} Documents_Copy_MarkdownInputs */
/** @typedef {{}} Documents_Copy_HtmlInputs */
/** @typedef {{}} Documents_Copy_JsonInputs */
/** @typedef {{}} Documents_CopiedInputs */
/** @typedef {{}} Documents_No_Json_AnnotationInputs */
/** @typedef {{}} Documents_No_Markdown_ContentInputs */
/** @typedef {{}} Documents_No_Preview_AvailableInputs */
/** @typedef {{}} Documents_CloseInputs */
/** @typedef {{}} Documents_MoreInputs */
/** @typedef {{}} Documents_OpenInputs */
/** @typedef {{}} Documents_ShareInputs */
export const common_language: (inputs: Common_LanguageInputs) => LocalizedString;
export const common_english: (inputs: Common_EnglishInputs) => LocalizedString;
export const common_romanian: (inputs: Common_RomanianInputs) => LocalizedString;
export const common_cancel: (inputs: Common_CancelInputs) => LocalizedString;
export const common_delete: (inputs: Common_DeleteInputs) => LocalizedString;
export const common_retry: (inputs: Common_RetryInputs) => LocalizedString;
export const common_previous: (inputs: Common_PreviousInputs) => LocalizedString;
export const common_next: (inputs: Common_NextInputs) => LocalizedString;
export const common_rows_per_page: (inputs: Common_Rows_Per_PageInputs) => LocalizedString;
export const common_strict: (inputs: Common_StrictInputs) => LocalizedString;
export const common_flexible: (inputs: Common_FlexibleInputs) => LocalizedString;
export const common_required: (inputs: Common_RequiredInputs) => LocalizedString;
export const common_unknown: (inputs: Common_UnknownInputs) => LocalizedString;
export const common_actions: (inputs: Common_ActionsInputs) => LocalizedString;
export const common_toggle_theme: (inputs: Common_Toggle_ThemeInputs) => LocalizedString;
export const header_credits_unavailable: (inputs: Header_Credits_UnavailableInputs) => LocalizedString;
export const header_credits: (inputs: Header_CreditsInputs) => LocalizedString;
export const header_credit_balance_unavailable: (inputs: Header_Credit_Balance_UnavailableInputs) => LocalizedString;
export const nav_account: (inputs: Nav_AccountInputs) => LocalizedString;
export const nav_no_email_address: (inputs: Nav_No_Email_AddressInputs) => LocalizedString;
export const nav_notifications: (inputs: Nav_NotificationsInputs) => LocalizedString;
export const nav_log_out: (inputs: Nav_Log_OutInputs) => LocalizedString;
export const nav_logout_title: (inputs: Nav_Logout_TitleInputs) => LocalizedString;
export const nav_logout_description: (inputs: Nav_Logout_DescriptionInputs) => LocalizedString;
export const nav_logout_failed: (inputs: Nav_Logout_FailedInputs) => LocalizedString;
export const nav_account_linked: (inputs: Nav_Account_LinkedInputs) => LocalizedString;
export const nav_account_link_conflict: (inputs: Nav_Account_Link_ConflictInputs) => LocalizedString;
export const nav_account_link_denied: (inputs: Nav_Account_Link_DeniedInputs) => LocalizedString;
export const nav_account_link_not_configured: (inputs: Nav_Account_Link_Not_ConfiguredInputs) => LocalizedString;
export const nav_account_link_sign_in_again: (inputs: Nav_Account_Link_Sign_In_AgainInputs) => LocalizedString;
export const nav_account_link_failed: (inputs: Nav_Account_Link_FailedInputs) => LocalizedString;
export const nav_dashboard: (inputs: Nav_DashboardInputs) => LocalizedString;
export const nav_schemas: (inputs: Nav_SchemasInputs) => LocalizedString;
export const nav_new_schema: (inputs: Nav_New_SchemaInputs) => LocalizedString;
export const nav_edit_schema: (inputs: Nav_Edit_SchemaInputs) => LocalizedString;
export const nav_jobs: (inputs: Nav_JobsInputs) => LocalizedString;
export const nav_new_job: (inputs: Nav_New_JobInputs) => LocalizedString;
export const nav_billing: (inputs: Nav_BillingInputs) => LocalizedString;
export const nav_billing_orders: (inputs: Nav_Billing_OrdersInputs) => LocalizedString;
export const nav_credit_usage_history: (inputs: Nav_Credit_Usage_HistoryInputs) => LocalizedString;
export const nav_developer_settings: (inputs: Nav_Developer_SettingsInputs) => LocalizedString;
export const nav_get_help: (inputs: Nav_Get_HelpInputs) => LocalizedString;
export const nav_quick_ocr: (inputs: Nav_Quick_OcrInputs) => LocalizedString;
export const nav_create_quick_ocr_job: (inputs: Nav_Create_Quick_Ocr_JobInputs) => LocalizedString;
export const nav_create_schema: (inputs: Nav_Create_SchemaInputs) => LocalizedString;
export const nav_create_job: (inputs: Nav_Create_JobInputs) => LocalizedString;
export const dashboard_metric_documents_processed: (inputs: Dashboard_Metric_Documents_ProcessedInputs) => LocalizedString;
export const dashboard_page_description: (inputs: Dashboard_Page_DescriptionInputs) => LocalizedString;
export const dashboard_refreshing: (inputs: Dashboard_RefreshingInputs) => LocalizedString;
export const dashboard_loading_title: (inputs: Dashboard_Loading_TitleInputs) => LocalizedString;
export const dashboard_loading_description: (inputs: Dashboard_Loading_DescriptionInputs) => LocalizedString;
export const dashboard_warning_title: (inputs: Dashboard_Warning_TitleInputs) => LocalizedString;
export const dashboard_unavailable_title: (inputs: Dashboard_Unavailable_TitleInputs) => LocalizedString;
export const dashboard_unavailable_default: (inputs: Dashboard_Unavailable_DefaultInputs) => LocalizedString;
export const dashboard_metric_pages_processed: (inputs: Dashboard_Metric_Pages_ProcessedInputs) => LocalizedString;
export const dashboard_metric_completion_rate: (inputs: Dashboard_Metric_Completion_RateInputs) => LocalizedString;
export const dashboard_metric_credits_spent: (inputs: Dashboard_Metric_Credits_SpentInputs) => LocalizedString;
export const dashboard_jobs_in_progress_one: (inputs: Dashboard_Jobs_In_Progress_OneInputs) => LocalizedString;
export const dashboard_jobs_in_progress_other: (inputs: Dashboard_Jobs_In_Progress_OtherInputs) => LocalizedString;
export const dashboard_pages_completed: (inputs: Dashboard_Pages_CompletedInputs) => LocalizedString;
export const dashboard_completion_summary: (inputs: Dashboard_Completion_SummaryInputs) => LocalizedString;
export const dashboard_credits_available_short: (inputs: Dashboard_Credits_Available_ShortInputs) => LocalizedString;
export const dashboard_metrics_aria: (inputs: Dashboard_Metrics_AriaInputs) => LocalizedString;
export const dashboard_documents_processed_title: (inputs: Dashboard_Documents_Processed_TitleInputs) => LocalizedString;
export const dashboard_chart_documents_label: (inputs: Dashboard_Chart_Documents_LabelInputs) => LocalizedString;
export const dashboard_select_range: (inputs: Dashboard_Select_RangeInputs) => LocalizedString;
export const dashboard_range_7d: (inputs: Dashboard_Range_7dInputs) => LocalizedString;
export const dashboard_range_30d: (inputs: Dashboard_Range_30dInputs) => LocalizedString;
export const dashboard_range_90d: (inputs: Dashboard_Range_90dInputs) => LocalizedString;
export const dashboard_recent_documents_title: (inputs: Dashboard_Recent_Documents_TitleInputs) => LocalizedString;
export const dashboard_recent_documents_description: (inputs: Dashboard_Recent_Documents_DescriptionInputs) => LocalizedString;
export const dashboard_view: (inputs: Dashboard_ViewInputs) => LocalizedString;
export const dashboard_no_saved_schema: (inputs: Dashboard_No_Saved_SchemaInputs) => LocalizedString;
export const dashboard_pages_one: (inputs: Dashboard_Pages_OneInputs) => LocalizedString;
export const dashboard_pages_other: (inputs: Dashboard_Pages_OtherInputs) => LocalizedString;
export const dashboard_no_completed_documents: (inputs: Dashboard_No_Completed_DocumentsInputs) => LocalizedString;
export const dashboard_schema_throughput_title: (inputs: Dashboard_Schema_Throughput_TitleInputs) => LocalizedString;
export const dashboard_schema_throughput_description: (inputs: Dashboard_Schema_Throughput_DescriptionInputs) => LocalizedString;
export const dashboard_documents_processed_one: (inputs: Dashboard_Documents_Processed_OneInputs) => LocalizedString;
export const dashboard_documents_processed_other: (inputs: Dashboard_Documents_Processed_OtherInputs) => LocalizedString;
export const dashboard_no_schema_throughput: (inputs: Dashboard_No_Schema_ThroughputInputs) => LocalizedString;
export const dashboard_datasets_title: (inputs: Dashboard_Datasets_TitleInputs) => LocalizedString;
export const dashboard_total_datasets_one: (inputs: Dashboard_Total_Datasets_OneInputs) => LocalizedString;
export const dashboard_total_datasets_other: (inputs: Dashboard_Total_Datasets_OtherInputs) => LocalizedString;
export const dashboard_fields_one: (inputs: Dashboard_Fields_OneInputs) => LocalizedString;
export const dashboard_fields_other: (inputs: Dashboard_Fields_OtherInputs) => LocalizedString;
export const dashboard_no_datasets: (inputs: Dashboard_No_DatasetsInputs) => LocalizedString;
export const dashboard_credits_title: (inputs: Dashboard_Credits_TitleInputs) => LocalizedString;
export const dashboard_credits_description: (inputs: Dashboard_Credits_DescriptionInputs) => LocalizedString;
export const dashboard_low_credit: (inputs: Dashboard_Low_CreditInputs) => LocalizedString;
export const dashboard_available_credits: (inputs: Dashboard_Available_CreditsInputs) => LocalizedString;
export const dashboard_credits_spent_in_range: (inputs: Dashboard_Credits_Spent_In_RangeInputs) => LocalizedString;
export const dashboard_billing: (inputs: Dashboard_BillingInputs) => LocalizedString;
export const dashboard_onboarding_title: (inputs: Dashboard_Onboarding_TitleInputs) => LocalizedString;
export const dashboard_onboarding_description: (inputs: Dashboard_Onboarding_DescriptionInputs) => LocalizedString;
export const dashboard_new_ocr_job: (inputs: Dashboard_New_Ocr_JobInputs) => LocalizedString;
export const dashboard_credits_one: (inputs: Dashboard_Credits_OneInputs) => LocalizedString;
export const dashboard_credits_other: (inputs: Dashboard_Credits_OtherInputs) => LocalizedString;
export const dashboard_step_schema: (inputs: Dashboard_Step_SchemaInputs) => LocalizedString;
export const dashboard_step_ocr_job: (inputs: Dashboard_Step_Ocr_JobInputs) => LocalizedString;
export const dashboard_step_dataset: (inputs: Dashboard_Step_DatasetInputs) => LocalizedString;
export const dashboard_step_api_key: (inputs: Dashboard_Step_Api_KeyInputs) => LocalizedString;
export const dashboard_step_webhook: (inputs: Dashboard_Step_WebhookInputs) => LocalizedString;
export const dashboard_step_ready: (inputs: Dashboard_Step_ReadyInputs) => LocalizedString;
export const dashboard_step_open: (inputs: Dashboard_Step_OpenInputs) => LocalizedString;
export const admin_nav_users: (inputs: Admin_Nav_UsersInputs) => LocalizedString;
export const admin_nav_user: (inputs: Admin_Nav_UserInputs) => LocalizedString;
export const admin_nav_invoices: (inputs: Admin_Nav_InvoicesInputs) => LocalizedString;
export const admin_nav_orders: (inputs: Admin_Nav_OrdersInputs) => LocalizedString;
export const admin_nav_json_recipes: (inputs: Admin_Nav_Json_RecipesInputs) => LocalizedString;
export const admin_nav_admin: (inputs: Admin_Nav_AdminInputs) => LocalizedString;
export const admin_user_fallback: (inputs: Admin_User_FallbackInputs) => LocalizedString;
export const sidebar_syncra: (inputs: Sidebar_SyncraInputs) => LocalizedString;
export const sidebar_syncra_admin: (inputs: Sidebar_Syncra_AdminInputs) => LocalizedString;
export const sidebar_user_space: (inputs: Sidebar_User_SpaceInputs) => LocalizedString;
export const sidebar_admin_portal: (inputs: Sidebar_Admin_PortalInputs) => LocalizedString;
export const sidebar_switch_space: (inputs: Sidebar_Switch_SpaceInputs) => LocalizedString;
export const schemas_new_title: (inputs: Schemas_New_TitleInputs) => LocalizedString;
export const schemas_library: (inputs: Schemas_LibraryInputs) => LocalizedString;
export const schemas_new_description: (inputs: Schemas_New_DescriptionInputs) => LocalizedString;
export const schemas_edit_title: (inputs: Schemas_Edit_TitleInputs) => LocalizedString;
export const schemas_edit_description: (inputs: Schemas_Edit_DescriptionInputs) => LocalizedString;
export const schemas_save_schema: (inputs: Schemas_Save_SchemaInputs) => LocalizedString;
export const schemas_save_changes: (inputs: Schemas_Save_ChangesInputs) => LocalizedString;
export const schemas_saved_success: (inputs: Schemas_Saved_SuccessInputs) => LocalizedString;
export const schemas_saved_success_with_id: (inputs: Schemas_Saved_Success_With_IdInputs) => LocalizedString;
export const schemas_saved_feedback: (inputs: Schemas_Saved_FeedbackInputs) => LocalizedString;
export const schemas_empty_schema_error: (inputs: Schemas_Empty_Schema_ErrorInputs) => LocalizedString;
export const schemas_delete_single_title: (inputs: Schemas_Delete_Single_TitleInputs) => LocalizedString;
export const schemas_delete_single_description: (inputs: Schemas_Delete_Single_DescriptionInputs) => LocalizedString;
export const schemas_delete_bulk_title_one: (inputs: Schemas_Delete_Bulk_Title_OneInputs) => LocalizedString;
export const schemas_delete_bulk_title_other: (inputs: Schemas_Delete_Bulk_Title_OtherInputs) => LocalizedString;
export const schemas_delete_bulk_description_one: (inputs: Schemas_Delete_Bulk_Description_OneInputs) => LocalizedString;
export const schemas_delete_bulk_description_other: (inputs: Schemas_Delete_Bulk_Description_OtherInputs) => LocalizedString;
export const schemas_select_all_on_page: (inputs: Schemas_Select_All_On_PageInputs) => LocalizedString;
export const schemas_select_schema: (inputs: Schemas_Select_SchemaInputs) => LocalizedString;
export const schemas_name_column: (inputs: Schemas_Name_ColumnInputs) => LocalizedString;
export const schemas_id_column: (inputs: Schemas_Id_ColumnInputs) => LocalizedString;
export const schemas_id_label: (inputs: Schemas_Id_LabelInputs) => LocalizedString;
export const schemas_copy_id: (inputs: Schemas_Copy_IdInputs) => LocalizedString;
export const schemas_copy_id_aria: (inputs: Schemas_Copy_Id_AriaInputs) => LocalizedString;
export const schemas_copy_id_success: (inputs: Schemas_Copy_Id_SuccessInputs) => LocalizedString;
export const schemas_copy_id_error: (inputs: Schemas_Copy_Id_ErrorInputs) => LocalizedString;
export const schemas_strict_mode_column: (inputs: Schemas_Strict_Mode_ColumnInputs) => LocalizedString;
export const schemas_created_column: (inputs: Schemas_Created_ColumnInputs) => LocalizedString;
export const schemas_updated_column: (inputs: Schemas_Updated_ColumnInputs) => LocalizedString;
export const schemas_new_schema: (inputs: Schemas_New_SchemaInputs) => LocalizedString;
export const schemas_no_schemas_found: (inputs: Schemas_No_Schemas_FoundInputs) => LocalizedString;
export const schemas_empty_body: (inputs: Schemas_Empty_BodyInputs) => LocalizedString;
export const schemas_create_schema: (inputs: Schemas_Create_SchemaInputs) => LocalizedString;
export const schemas_showing_schemas_one: (inputs: Schemas_Showing_Schemas_OneInputs) => LocalizedString;
export const schemas_showing_schemas_other: (inputs: Schemas_Showing_Schemas_OtherInputs) => LocalizedString;
export const schemas_no_schemas_to_show: (inputs: Schemas_No_Schemas_To_ShowInputs) => LocalizedString;
export const schemas_selected_count_one: (inputs: Schemas_Selected_Count_OneInputs) => LocalizedString;
export const schemas_selected_count_other: (inputs: Schemas_Selected_Count_OtherInputs) => LocalizedString;
export const schemas_deleting: (inputs: Schemas_DeletingInputs) => LocalizedString;
export const schemas_no_description: (inputs: Schemas_No_DescriptionInputs) => LocalizedString;
export const schemas_sort_created_ascending: (inputs: Schemas_Sort_Created_AscendingInputs) => LocalizedString;
export const schemas_sort_created_descending: (inputs: Schemas_Sort_Created_DescendingInputs) => LocalizedString;
export const schemas_edit_aria: (inputs: Schemas_Edit_AriaInputs) => LocalizedString;
export const schemas_create_job_with: (inputs: Schemas_Create_Job_WithInputs) => LocalizedString;
export const schemas_clone_aria: (inputs: Schemas_Clone_AriaInputs) => LocalizedString;
export const schemas_delete_aria: (inputs: Schemas_Delete_AriaInputs) => LocalizedString;
export const schemas_loading_schema: (inputs: Schemas_Loading_SchemaInputs) => LocalizedString;
export const schemas_not_found_title: (inputs: Schemas_Not_Found_TitleInputs) => LocalizedString;
export const schemas_not_found_body: (inputs: Schemas_Not_Found_BodyInputs) => LocalizedString;
export const schemas_view_schemas: (inputs: Schemas_View_SchemasInputs) => LocalizedString;
export const schemas_could_not_load: (inputs: Schemas_Could_Not_LoadInputs) => LocalizedString;
export const schemas_editor_badge: (inputs: Schemas_Editor_BadgeInputs) => LocalizedString;
export const schemas_general_settings: (inputs: Schemas_General_SettingsInputs) => LocalizedString;
export const schemas_schema_name_label: (inputs: Schemas_Schema_Name_LabelInputs) => LocalizedString;
export const schemas_schema_name_placeholder: (inputs: Schemas_Schema_Name_PlaceholderInputs) => LocalizedString;
export const schemas_description_label: (inputs: Schemas_Description_LabelInputs) => LocalizedString;
export const schemas_description_placeholder: (inputs: Schemas_Description_PlaceholderInputs) => LocalizedString;
export const schemas_strict_mode: (inputs: Schemas_Strict_ModeInputs) => LocalizedString;
export const schemas_flexible_mode: (inputs: Schemas_Flexible_ModeInputs) => LocalizedString;
export const schemas_strict_mode_description: (inputs: Schemas_Strict_Mode_DescriptionInputs) => LocalizedString;
export const schemas_structure_designer: (inputs: Schemas_Structure_DesignerInputs) => LocalizedString;
export const schemas_visual_node_designer: (inputs: Schemas_Visual_Node_DesignerInputs) => LocalizedString;
export const schemas_validation_name_required: (inputs: Schemas_Validation_Name_RequiredInputs) => LocalizedString;
export const schemas_validation_name_too_long: (inputs: Schemas_Validation_Name_Too_LongInputs) => LocalizedString;
export const schemas_validation_schema_object: (inputs: Schemas_Validation_Schema_ObjectInputs) => LocalizedString;
export const schemas_clone: (inputs: Schemas_CloneInputs) => LocalizedString;
export const schemas_cloning: (inputs: Schemas_CloningInputs) => LocalizedString;
export const schemas_saving: (inputs: Schemas_SavingInputs) => LocalizedString;
export const json_recipes_title: (inputs: Json_Recipes_TitleInputs) => LocalizedString;
export const json_recipes_description: (inputs: Json_Recipes_DescriptionInputs) => LocalizedString;
export const json_recipes_new_recipe: (inputs: Json_Recipes_New_RecipeInputs) => LocalizedString;
export const json_recipes_no_recipes_found: (inputs: Json_Recipes_No_Recipes_FoundInputs) => LocalizedString;
export const json_recipes_empty_body: (inputs: Json_Recipes_Empty_BodyInputs) => LocalizedString;
export const json_recipes_loading: (inputs: Json_Recipes_LoadingInputs) => LocalizedString;
export const json_recipes_loading_recipe: (inputs: Json_Recipes_Loading_RecipeInputs) => LocalizedString;
export const json_recipes_counter_column: (inputs: Json_Recipes_Counter_ColumnInputs) => LocalizedString;
export const json_recipes_created_column: (inputs: Json_Recipes_Created_ColumnInputs) => LocalizedString;
export const json_recipes_updated_column: (inputs: Json_Recipes_Updated_ColumnInputs) => LocalizedString;
export const json_recipes_json_fields_column: (inputs: Json_Recipes_Json_Fields_ColumnInputs) => LocalizedString;
export const json_recipes_sort_created_ascending: (inputs: Json_Recipes_Sort_Created_AscendingInputs) => LocalizedString;
export const json_recipes_sort_created_descending: (inputs: Json_Recipes_Sort_Created_DescendingInputs) => LocalizedString;
export const json_recipes_showing_one: (inputs: Json_Recipes_Showing_OneInputs) => LocalizedString;
export const json_recipes_showing_other: (inputs: Json_Recipes_Showing_OtherInputs) => LocalizedString;
export const json_recipes_no_recipes_to_show: (inputs: Json_Recipes_No_Recipes_To_ShowInputs) => LocalizedString;
export const json_recipes_edit_aria: (inputs: Json_Recipes_Edit_AriaInputs) => LocalizedString;
export const json_recipes_delete_aria: (inputs: Json_Recipes_Delete_AriaInputs) => LocalizedString;
export const json_recipes_new_title: (inputs: Json_Recipes_New_TitleInputs) => LocalizedString;
export const json_recipes_new_description: (inputs: Json_Recipes_New_DescriptionInputs) => LocalizedString;
export const json_recipes_edit_title: (inputs: Json_Recipes_Edit_TitleInputs) => LocalizedString;
export const json_recipes_edit_description: (inputs: Json_Recipes_Edit_DescriptionInputs) => LocalizedString;
export const json_recipes_save_recipe: (inputs: Json_Recipes_Save_RecipeInputs) => LocalizedString;
export const json_recipes_save_changes: (inputs: Json_Recipes_Save_ChangesInputs) => LocalizedString;
export const json_recipes_created_success: (inputs: Json_Recipes_Created_SuccessInputs) => LocalizedString;
export const json_recipes_saved_success: (inputs: Json_Recipes_Saved_SuccessInputs) => LocalizedString;
export const json_recipes_deleted_success: (inputs: Json_Recipes_Deleted_SuccessInputs) => LocalizedString;
export const json_recipes_delete_confirm: (inputs: Json_Recipes_Delete_ConfirmInputs) => LocalizedString;
export const json_recipes_not_found_title: (inputs: Json_Recipes_Not_Found_TitleInputs) => LocalizedString;
export const json_recipes_not_found_body: (inputs: Json_Recipes_Not_Found_BodyInputs) => LocalizedString;
export const json_recipes_view_recipes: (inputs: Json_Recipes_View_RecipesInputs) => LocalizedString;
export const json_recipes_could_not_load: (inputs: Json_Recipes_Could_Not_LoadInputs) => LocalizedString;
export const json_recipes_editor_badge: (inputs: Json_Recipes_Editor_BadgeInputs) => LocalizedString;
export const json_recipes_general_settings: (inputs: Json_Recipes_General_SettingsInputs) => LocalizedString;
export const json_recipes_title_label: (inputs: Json_Recipes_Title_LabelInputs) => LocalizedString;
export const json_recipes_title_placeholder: (inputs: Json_Recipes_Title_PlaceholderInputs) => LocalizedString;
export const json_recipes_description_label: (inputs: Json_Recipes_Description_LabelInputs) => LocalizedString;
export const json_recipes_description_placeholder: (inputs: Json_Recipes_Description_PlaceholderInputs) => LocalizedString;
export const json_recipes_structure_designer: (inputs: Json_Recipes_Structure_DesignerInputs) => LocalizedString;
export const json_recipes_visual_node_designer: (inputs: Json_Recipes_Visual_Node_DesignerInputs) => LocalizedString;
export const json_recipes_category_label: (inputs: Json_Recipes_Category_LabelInputs) => LocalizedString;
export const json_recipes_others: (inputs: Json_Recipes_OthersInputs) => LocalizedString;
export const json_recipes_manage_categories: (inputs: Json_Recipes_Manage_CategoriesInputs) => LocalizedString;
export const json_recipes_validation_title_required: (inputs: Json_Recipes_Validation_Title_RequiredInputs) => LocalizedString;
export const json_recipes_validation_title_too_long: (inputs: Json_Recipes_Validation_Title_Too_LongInputs) => LocalizedString;
export const json_recipes_validation_json_object: (inputs: Json_Recipes_Validation_Json_ObjectInputs) => LocalizedString;
export const json_recipes_saving: (inputs: Json_Recipes_SavingInputs) => LocalizedString;
export const json_recipes_deleting: (inputs: Json_Recipes_DeletingInputs) => LocalizedString;
export const json_recipe_categories_title: (inputs: Json_Recipe_Categories_TitleInputs) => LocalizedString;
export const json_recipe_categories_description: (inputs: Json_Recipe_Categories_DescriptionInputs) => LocalizedString;
export const json_recipe_categories_title_en_label: (inputs: Json_Recipe_Categories_Title_En_LabelInputs) => LocalizedString;
export const json_recipe_categories_title_ro_label: (inputs: Json_Recipe_Categories_Title_Ro_LabelInputs) => LocalizedString;
export const json_recipe_categories_create_category: (inputs: Json_Recipe_Categories_Create_CategoryInputs) => LocalizedString;
export const json_recipe_categories_save_category: (inputs: Json_Recipe_Categories_Save_CategoryInputs) => LocalizedString;
export const json_recipe_categories_edit_title: (inputs: Json_Recipe_Categories_Edit_TitleInputs) => LocalizedString;
export const json_recipe_categories_delete_confirm: (inputs: Json_Recipe_Categories_Delete_ConfirmInputs) => LocalizedString;
export const json_recipe_categories_loading: (inputs: Json_Recipe_Categories_LoadingInputs) => LocalizedString;
export const json_recipe_categories_could_not_load: (inputs: Json_Recipe_Categories_Could_Not_LoadInputs) => LocalizedString;
export const json_recipe_categories_empty_title: (inputs: Json_Recipe_Categories_Empty_TitleInputs) => LocalizedString;
export const json_recipe_categories_empty_body: (inputs: Json_Recipe_Categories_Empty_BodyInputs) => LocalizedString;
export const json_recipe_categories_created_success: (inputs: Json_Recipe_Categories_Created_SuccessInputs) => LocalizedString;
export const json_recipe_categories_saved_success: (inputs: Json_Recipe_Categories_Saved_SuccessInputs) => LocalizedString;
export const json_recipe_categories_deleted_success: (inputs: Json_Recipe_Categories_Deleted_SuccessInputs) => LocalizedString;
export const json_recipe_categories_validation_titles_required: (inputs: Json_Recipe_Categories_Validation_Titles_RequiredInputs) => LocalizedString;
export const json_recipe_categories_validation_titles_too_long: (inputs: Json_Recipe_Categories_Validation_Titles_Too_LongInputs) => LocalizedString;
export const json_recipe_categories_edit_aria: (inputs: Json_Recipe_Categories_Edit_AriaInputs) => LocalizedString;
export const json_recipe_categories_delete_aria: (inputs: Json_Recipe_Categories_Delete_AriaInputs) => LocalizedString;
export const ocr_recipes_nav: (inputs: Ocr_Recipes_NavInputs) => LocalizedString;
export const ocr_recipes_title: (inputs: Ocr_Recipes_TitleInputs) => LocalizedString;
export const ocr_recipes_meta_description: (inputs: Ocr_Recipes_Meta_DescriptionInputs) => LocalizedString;
export const ocr_recipes_eyebrow: (inputs: Ocr_Recipes_EyebrowInputs) => LocalizedString;
export const ocr_recipes_hero_title: (inputs: Ocr_Recipes_Hero_TitleInputs) => LocalizedString;
export const ocr_recipes_hero_description: (inputs: Ocr_Recipes_Hero_DescriptionInputs) => LocalizedString;
export const ocr_recipes_search_label: (inputs: Ocr_Recipes_Search_LabelInputs) => LocalizedString;
export const ocr_recipes_search_placeholder: (inputs: Ocr_Recipes_Search_PlaceholderInputs) => LocalizedString;
export const ocr_recipes_category_filter: (inputs: Ocr_Recipes_Category_FilterInputs) => LocalizedString;
export const ocr_recipes_all_categories: (inputs: Ocr_Recipes_All_CategoriesInputs) => LocalizedString;
export const ocr_recipes_sort_label: (inputs: Ocr_Recipes_Sort_LabelInputs) => LocalizedString;
export const ocr_recipes_sort_popular: (inputs: Ocr_Recipes_Sort_PopularInputs) => LocalizedString;
export const ocr_recipes_sort_newest: (inputs: Ocr_Recipes_Sort_NewestInputs) => LocalizedString;
export const ocr_recipes_sort_az: (inputs: Ocr_Recipes_Sort_AzInputs) => LocalizedString;
export const ocr_recipes_showing_one: (inputs: Ocr_Recipes_Showing_OneInputs) => LocalizedString;
export const ocr_recipes_showing_other: (inputs: Ocr_Recipes_Showing_OtherInputs) => LocalizedString;
export const ocr_recipes_no_matches_title: (inputs: Ocr_Recipes_No_Matches_TitleInputs) => LocalizedString;
export const ocr_recipes_no_matches_body: (inputs: Ocr_Recipes_No_Matches_BodyInputs) => LocalizedString;
export const ocr_recipes_others: (inputs: Ocr_Recipes_OthersInputs) => LocalizedString;
export const ocr_recipes_fields_one: (inputs: Ocr_Recipes_Fields_OneInputs) => LocalizedString;
export const ocr_recipes_fields_other: (inputs: Ocr_Recipes_Fields_OtherInputs) => LocalizedString;
export const ocr_recipes_required_one: (inputs: Ocr_Recipes_Required_OneInputs) => LocalizedString;
export const ocr_recipes_required_other: (inputs: Ocr_Recipes_Required_OtherInputs) => LocalizedString;
export const ocr_recipes_deploys_one: (inputs: Ocr_Recipes_Deploys_OneInputs) => LocalizedString;
export const ocr_recipes_deploys_other: (inputs: Ocr_Recipes_Deploys_OtherInputs) => LocalizedString;
export const ocr_recipes_json_fields: (inputs: Ocr_Recipes_Json_FieldsInputs) => LocalizedString;
export const ocr_recipes_system_recipe: (inputs: Ocr_Recipes_System_RecipeInputs) => LocalizedString;
export const ocr_recipes_strict_schema: (inputs: Ocr_Recipes_Strict_SchemaInputs) => LocalizedString;
export const ocr_recipes_required: (inputs: Ocr_Recipes_RequiredInputs) => LocalizedString;
export const ocr_recipes_preview_json: (inputs: Ocr_Recipes_Preview_JsonInputs) => LocalizedString;
export const ocr_recipes_no_fields: (inputs: Ocr_Recipes_No_FieldsInputs) => LocalizedString;
export const ocr_recipes_clone_recipe: (inputs: Ocr_Recipes_Clone_RecipeInputs) => LocalizedString;
export const ocr_recipes_clone_aria: (inputs: Ocr_Recipes_Clone_AriaInputs) => LocalizedString;
export const ocr_recipes_log_in_to_clone: (inputs: Ocr_Recipes_Log_In_To_CloneInputs) => LocalizedString;
export const ocr_recipes_clone_failed: (inputs: Ocr_Recipes_Clone_FailedInputs) => LocalizedString;
export const ocr_recipes_load_failed: (inputs: Ocr_Recipes_Load_FailedInputs) => LocalizedString;
export const jobs_page_title: (inputs: Jobs_Page_TitleInputs) => LocalizedString;
export const jobs_missing_schema_id: (inputs: Jobs_Missing_Schema_IdInputs) => LocalizedString;
export const jobs_missing_job_id: (inputs: Jobs_Missing_Job_IdInputs) => LocalizedString;
export const jobs_delete_bulk_title_one: (inputs: Jobs_Delete_Bulk_Title_OneInputs) => LocalizedString;
export const jobs_delete_bulk_title_other: (inputs: Jobs_Delete_Bulk_Title_OtherInputs) => LocalizedString;
export const jobs_delete_bulk_description_one: (inputs: Jobs_Delete_Bulk_Description_OneInputs) => LocalizedString;
export const jobs_delete_bulk_description_other: (inputs: Jobs_Delete_Bulk_Description_OtherInputs) => LocalizedString;
export const jobs_delete_single_title: (inputs: Jobs_Delete_Single_TitleInputs) => LocalizedString;
export const jobs_delete_single_description: (inputs: Jobs_Delete_Single_DescriptionInputs) => LocalizedString;
export const jobs_status_queued: (inputs: Jobs_Status_QueuedInputs) => LocalizedString;
export const jobs_status_pending: (inputs: Jobs_Status_PendingInputs) => LocalizedString;
export const jobs_status_processing: (inputs: Jobs_Status_ProcessingInputs) => LocalizedString;
export const jobs_status_completed: (inputs: Jobs_Status_CompletedInputs) => LocalizedString;
export const jobs_status_failed: (inputs: Jobs_Status_FailedInputs) => LocalizedString;
export const jobs_inline_schema: (inputs: Jobs_Inline_SchemaInputs) => LocalizedString;
export const jobs_no_schema: (inputs: Jobs_No_SchemaInputs) => LocalizedString;
export const jobs_schema: (inputs: Jobs_SchemaInputs) => LocalizedString;
export const jobs_select_all_on_page: (inputs: Jobs_Select_All_On_PageInputs) => LocalizedString;
export const jobs_select_job: (inputs: Jobs_Select_JobInputs) => LocalizedString;
export const jobs_filename_column: (inputs: Jobs_Filename_ColumnInputs) => LocalizedString;
export const jobs_status_column: (inputs: Jobs_Status_ColumnInputs) => LocalizedString;
export const jobs_created_column: (inputs: Jobs_Created_ColumnInputs) => LocalizedString;
export const jobs_file_size_column: (inputs: Jobs_File_Size_ColumnInputs) => LocalizedString;
export const jobs_pages_column: (inputs: Jobs_Pages_ColumnInputs) => LocalizedString;
export const jobs_new_job: (inputs: Jobs_New_JobInputs) => LocalizedString;
export const jobs_no_jobs_found: (inputs: Jobs_No_Jobs_FoundInputs) => LocalizedString;
export const jobs_empty_body: (inputs: Jobs_Empty_BodyInputs) => LocalizedString;
export const jobs_showing_jobs_one: (inputs: Jobs_Showing_Jobs_OneInputs) => LocalizedString;
export const jobs_showing_jobs_other: (inputs: Jobs_Showing_Jobs_OtherInputs) => LocalizedString;
export const jobs_no_jobs_to_show: (inputs: Jobs_No_Jobs_To_ShowInputs) => LocalizedString;
export const jobs_selected_count_one: (inputs: Jobs_Selected_Count_OneInputs) => LocalizedString;
export const jobs_selected_count_other: (inputs: Jobs_Selected_Count_OtherInputs) => LocalizedString;
export const jobs_deleting: (inputs: Jobs_DeletingInputs) => LocalizedString;
export const jobs_delete_job: (inputs: Jobs_Delete_JobInputs) => LocalizedString;
export const jobs_saved_extraction_schema: (inputs: Jobs_Saved_Extraction_SchemaInputs) => LocalizedString;
export const jobs_inline_schema_description: (inputs: Jobs_Inline_Schema_DescriptionInputs) => LocalizedString;
export const jobs_extraction_schema_details: (inputs: Jobs_Extraction_Schema_DetailsInputs) => LocalizedString;
export const new_job_missing_document_id: (inputs: New_Job_Missing_Document_IdInputs) => LocalizedString;
export const new_job_failed_create: (inputs: New_Job_Failed_CreateInputs) => LocalizedString;
export const new_job_insufficient_credits_buy: (inputs: New_Job_Insufficient_Credits_BuyInputs) => LocalizedString;
export const new_job_failed_load_document: (inputs: New_Job_Failed_Load_DocumentInputs) => LocalizedString;
export const new_job_invalid_document_response: (inputs: New_Job_Invalid_Document_ResponseInputs) => LocalizedString;
export const new_job_failed_load_schemas: (inputs: New_Job_Failed_Load_SchemasInputs) => LocalizedString;
export const new_job_invalid_schema_response: (inputs: New_Job_Invalid_Schema_ResponseInputs) => LocalizedString;
export const new_job_invalid_job_response: (inputs: New_Job_Invalid_Job_ResponseInputs) => LocalizedString;
export const new_job_failed_load_job: (inputs: New_Job_Failed_Load_JobInputs) => LocalizedString;
export const new_job_failed_poll_job: (inputs: New_Job_Failed_Poll_JobInputs) => LocalizedString;
export const new_job_select_schema: (inputs: New_Job_Select_SchemaInputs) => LocalizedString;
export const new_job_select_schema_placeholder: (inputs: New_Job_Select_Schema_PlaceholderInputs) => LocalizedString;
export const new_job_configure_payload_format: (inputs: New_Job_Configure_Payload_FormatInputs) => LocalizedString;
export const new_job_upload_documents: (inputs: New_Job_Upload_DocumentsInputs) => LocalizedString;
export const new_job_files_selected_one: (inputs: New_Job_Files_Selected_OneInputs) => LocalizedString;
export const new_job_files_selected_other: (inputs: New_Job_Files_Selected_OtherInputs) => LocalizedString;
export const new_job_drag_or_browse_files: (inputs: New_Job_Drag_Or_Browse_FilesInputs) => LocalizedString;
export const new_job_run_monitor: (inputs: New_Job_Run_MonitorInputs) => LocalizedString;
export const new_job_processing_batch: (inputs: New_Job_Processing_BatchInputs) => LocalizedString;
export const new_job_start_extraction_pipeline: (inputs: New_Job_Start_Extraction_PipelineInputs) => LocalizedString;
export const new_job_select_extraction_schema: (inputs: New_Job_Select_Extraction_SchemaInputs) => LocalizedString;
export const new_job_select_schema_description: (inputs: New_Job_Select_Schema_DescriptionInputs) => LocalizedString;
export const new_job_select_extraction_schema_aria: (inputs: New_Job_Select_Extraction_Schema_AriaInputs) => LocalizedString;
export const new_job_search_schemas: (inputs: New_Job_Search_SchemasInputs) => LocalizedString;
export const new_job_loading_schemas: (inputs: New_Job_Loading_SchemasInputs) => LocalizedString;
export const new_job_no_schemas_found: (inputs: New_Job_No_Schemas_FoundInputs) => LocalizedString;
export const new_job_no_schema_ocr_only: (inputs: New_Job_No_Schema_Ocr_OnlyInputs) => LocalizedString;
export const new_job_no_schema_description: (inputs: New_Job_No_Schema_DescriptionInputs) => LocalizedString;
export const new_job_no_personal_schemas: (inputs: New_Job_No_Personal_SchemasInputs) => LocalizedString;
export const new_job_create_one: (inputs: New_Job_Create_OneInputs) => LocalizedString;
export const new_job_selected_schema_help: (inputs: New_Job_Selected_Schema_HelpInputs) => LocalizedString;
export const new_job_no_schema_selected_help: (inputs: New_Job_No_Schema_Selected_HelpInputs) => LocalizedString;
export const new_job_target_mapped_fields: (inputs: New_Job_Target_Mapped_FieldsInputs) => LocalizedString;
export const new_job_no_fields_defined: (inputs: New_Job_No_Fields_DefinedInputs) => LocalizedString;
export const new_job_ocr_only_mode_active: (inputs: New_Job_Ocr_Only_Mode_ActiveInputs) => LocalizedString;
export const new_job_ocr_only_mode_body: (inputs: New_Job_Ocr_Only_Mode_BodyInputs) => LocalizedString;
export const new_job_upload_documents_description: (inputs: New_Job_Upload_Documents_DescriptionInputs) => LocalizedString;
export const new_job_dropzone_title: (inputs: New_Job_Dropzone_TitleInputs) => LocalizedString;
export const new_job_dropzone_description: (inputs: New_Job_Dropzone_DescriptionInputs) => LocalizedString;
export const new_job_browse_files: (inputs: New_Job_Browse_FilesInputs) => LocalizedString;
export const new_job_pending_upload_queue: (inputs: New_Job_Pending_Upload_QueueInputs) => LocalizedString;
export const new_job_clear_all: (inputs: New_Job_Clear_AllInputs) => LocalizedString;
export const new_job_remove_file: (inputs: New_Job_Remove_FileInputs) => LocalizedString;
export const new_job_extraction_queue_results: (inputs: New_Job_Extraction_Queue_ResultsInputs) => LocalizedString;
export const new_job_file_count_one: (inputs: New_Job_File_Count_OneInputs) => LocalizedString;
export const new_job_file_count_other: (inputs: New_Job_File_Count_OtherInputs) => LocalizedString;
export const new_job_total: (inputs: New_Job_TotalInputs) => LocalizedString;
export const new_job_active_batch_status: (inputs: New_Job_Active_Batch_StatusInputs) => LocalizedString;
export const new_job_active_batch_description: (inputs: New_Job_Active_Batch_DescriptionInputs) => LocalizedString;
export const new_job_progress: (inputs: New_Job_ProgressInputs) => LocalizedString;
export const new_job_total_files: (inputs: New_Job_Total_FilesInputs) => LocalizedString;
export const new_job_completed: (inputs: New_Job_CompletedInputs) => LocalizedString;
export const new_job_processing: (inputs: New_Job_ProcessingInputs) => LocalizedString;
export const new_job_failed: (inputs: New_Job_FailedInputs) => LocalizedString;
export const new_job_no_active_extraction_jobs: (inputs: New_Job_No_Active_Extraction_JobsInputs) => LocalizedString;
export const new_job_no_active_extraction_jobs_body: (inputs: New_Job_No_Active_Extraction_Jobs_BodyInputs) => LocalizedString;
export const new_job_preview_document: (inputs: New_Job_Preview_DocumentInputs) => LocalizedString;
export const new_job_preview_unavailable: (inputs: New_Job_Preview_UnavailableInputs) => LocalizedString;
export const new_job_remove_failed_job: (inputs: New_Job_Remove_Failed_JobInputs) => LocalizedString;
export const new_job_queueing_documents: (inputs: New_Job_Queueing_DocumentsInputs) => LocalizedString;
export const new_job_extracting_content: (inputs: New_Job_Extracting_ContentInputs) => LocalizedString;
export const new_job_run_extraction_one: (inputs: New_Job_Run_Extraction_OneInputs) => LocalizedString;
export const new_job_run_extraction_other: (inputs: New_Job_Run_Extraction_OtherInputs) => LocalizedString;
export const new_job_insufficient_credits_document: (inputs: New_Job_Insufficient_Credits_DocumentInputs) => LocalizedString;
export const new_job_processing_failed: (inputs: New_Job_Processing_FailedInputs) => LocalizedString;
export const new_job_processed: (inputs: New_Job_ProcessedInputs) => LocalizedString;
export const new_job_document_id: (inputs: New_Job_Document_IdInputs) => LocalizedString;
export const new_job_creating_job: (inputs: New_Job_Creating_JobInputs) => LocalizedString;
export const new_job_queued_processing: (inputs: New_Job_Queued_ProcessingInputs) => LocalizedString;
export const new_job_extracting_entities: (inputs: New_Job_Extracting_EntitiesInputs) => LocalizedString;
export const common_apply: (inputs: Common_ApplyInputs) => LocalizedString;
export const common_clear: (inputs: Common_ClearInputs) => LocalizedString;
export const common_saving: (inputs: Common_SavingInputs) => LocalizedString;
export const common_loading: (inputs: Common_LoadingInputs) => LocalizedString;
export const common_refresh: (inputs: Common_RefreshInputs) => LocalizedString;
export const common_connected: (inputs: Common_ConnectedInputs) => LocalizedString;
export const common_connect: (inputs: Common_ConnectInputs) => LocalizedString;
export const common_download: (inputs: Common_DownloadInputs) => LocalizedString;
export const common_today: (inputs: Common_TodayInputs) => LocalizedString;
export const common_this_week: (inputs: Common_This_WeekInputs) => LocalizedString;
export const common_this_month: (inputs: Common_This_MonthInputs) => LocalizedString;
export const common_any: (inputs: Common_AnyInputs) => LocalizedString;
export const billing_unavailable: (inputs: Billing_UnavailableInputs) => LocalizedString;
export const billing_credit_blocks_error: (inputs: Billing_Credit_Blocks_ErrorInputs) => LocalizedString;
export const billing_checkout_unavailable: (inputs: Billing_Checkout_UnavailableInputs) => LocalizedString;
export const billing_payment_received_title: (inputs: Billing_Payment_Received_TitleInputs) => LocalizedString;
export const billing_payment_received_body: (inputs: Billing_Payment_Received_BodyInputs) => LocalizedString;
export const billing_checkout_canceled_title: (inputs: Billing_Checkout_Canceled_TitleInputs) => LocalizedString;
export const billing_checkout_canceled_body: (inputs: Billing_Checkout_Canceled_BodyInputs) => LocalizedString;
export const billing_available_balance: (inputs: Billing_Available_BalanceInputs) => LocalizedString;
export const billing_conversion: (inputs: Billing_ConversionInputs) => LocalizedString;
export const billing_conversion_rate: (inputs: Billing_Conversion_RateInputs) => LocalizedString;
export const billing_balance_checked_upload: (inputs: Billing_Balance_Checked_UploadInputs) => LocalizedString;
export const billing_debited_after_success: (inputs: Billing_Debited_After_SuccessInputs) => LocalizedString;
export const billing_secure_stripe_checkout: (inputs: Billing_Secure_Stripe_CheckoutInputs) => LocalizedString;
export const billing_purchase_credits: (inputs: Billing_Purchase_CreditsInputs) => LocalizedString;
export const billing_credits_to_purchase: (inputs: Billing_Credits_To_PurchaseInputs) => LocalizedString;
export const billing_volume_discount_tiers: (inputs: Billing_Volume_Discount_TiersInputs) => LocalizedString;
export const billing_total_to_pay: (inputs: Billing_Total_To_PayInputs) => LocalizedString;
export const billing_base_price: (inputs: Billing_Base_PriceInputs) => LocalizedString;
export const billing_volume_discount: (inputs: Billing_Volume_DiscountInputs) => LocalizedString;
export const billing_starting_checkout: (inputs: Billing_Starting_CheckoutInputs) => LocalizedString;
export const billing_secure_checkout: (inputs: Billing_Secure_CheckoutInputs) => LocalizedString;
export const billing_buy_credits: (inputs: Billing_Buy_CreditsInputs) => LocalizedString;
export const billing_orders_page_title: (inputs: Billing_Orders_Page_TitleInputs) => LocalizedString;
export const billing_orders_order_date_filter: (inputs: Billing_Orders_Order_Date_FilterInputs) => LocalizedString;
export const billing_orders_amount_column: (inputs: Billing_Orders_Amount_ColumnInputs) => LocalizedString;
export const billing_orders_credits_column: (inputs: Billing_Orders_Credits_ColumnInputs) => LocalizedString;
export const billing_orders_status_column: (inputs: Billing_Orders_Status_ColumnInputs) => LocalizedString;
export const billing_orders_payment_datetime_column: (inputs: Billing_Orders_Payment_Datetime_ColumnInputs) => LocalizedString;
export const billing_orders_invoice_column: (inputs: Billing_Orders_Invoice_ColumnInputs) => LocalizedString;
export const billing_orders_presets: (inputs: Billing_Orders_PresetsInputs) => LocalizedString;
export const billing_orders_filter_status: (inputs: Billing_Orders_Filter_StatusInputs) => LocalizedString;
export const billing_orders_all_orders: (inputs: Billing_Orders_All_OrdersInputs) => LocalizedString;
export const billing_orders_clear_filters: (inputs: Billing_Orders_Clear_FiltersInputs) => LocalizedString;
export const billing_orders_clear_filters_action: (inputs: Billing_Orders_Clear_Filters_ActionInputs) => LocalizedString;
export const billing_orders_no_orders_found: (inputs: Billing_Orders_No_Orders_FoundInputs) => LocalizedString;
export const billing_orders_no_orders_yet: (inputs: Billing_Orders_No_Orders_YetInputs) => LocalizedString;
export const billing_orders_no_orders_match: (inputs: Billing_Orders_No_Orders_MatchInputs) => LocalizedString;
export const billing_orders_empty_body: (inputs: Billing_Orders_Empty_BodyInputs) => LocalizedString;
export const billing_orders_showing_one: (inputs: Billing_Orders_Showing_OneInputs) => LocalizedString;
export const billing_orders_showing_other: (inputs: Billing_Orders_Showing_OtherInputs) => LocalizedString;
export const billing_orders_none_to_show: (inputs: Billing_Orders_None_To_ShowInputs) => LocalizedString;
export const billing_orders_sort_order_date_ascending: (inputs: Billing_Orders_Sort_Order_Date_AscendingInputs) => LocalizedString;
export const billing_orders_sort_order_date_descending: (inputs: Billing_Orders_Sort_Order_Date_DescendingInputs) => LocalizedString;
export const billing_order_status_pending: (inputs: Billing_Order_Status_PendingInputs) => LocalizedString;
export const billing_order_status_paid: (inputs: Billing_Order_Status_PaidInputs) => LocalizedString;
export const billing_order_status_failed: (inputs: Billing_Order_Status_FailedInputs) => LocalizedString;
export const billing_order_status_refunded: (inputs: Billing_Order_Status_RefundedInputs) => LocalizedString;
export const billing_order_status_canceled: (inputs: Billing_Order_Status_CanceledInputs) => LocalizedString;
export const billing_orders_invoice_pdf_title: (inputs: Billing_Orders_Invoice_Pdf_TitleInputs) => LocalizedString;
export const billing_orders_invoice_preview_title: (inputs: Billing_Orders_Invoice_Preview_TitleInputs) => LocalizedString;
export const billing_orders_invoice_preview_description: (inputs: Billing_Orders_Invoice_Preview_DescriptionInputs) => LocalizedString;
export const billing_orders_invoice_iframe_title: (inputs: Billing_Orders_Invoice_Iframe_TitleInputs) => LocalizedString;
export const billing_orders_download_invoice: (inputs: Billing_Orders_Download_InvoiceInputs) => LocalizedString;
export const credit_usage_page_title: (inputs: Credit_Usage_Page_TitleInputs) => LocalizedString;
export const credit_usage_date_range_filter: (inputs: Credit_Usage_Date_Range_FilterInputs) => LocalizedString;
export const credit_usage_created_column: (inputs: Credit_Usage_Created_ColumnInputs) => LocalizedString;
export const credit_usage_type_column: (inputs: Credit_Usage_Type_ColumnInputs) => LocalizedString;
export const credit_usage_credits_column: (inputs: Credit_Usage_Credits_ColumnInputs) => LocalizedString;
export const credit_usage_related_id_column: (inputs: Credit_Usage_Related_Id_ColumnInputs) => LocalizedString;
export const credit_usage_filter_type: (inputs: Credit_Usage_Filter_TypeInputs) => LocalizedString;
export const credit_usage_all_activity: (inputs: Credit_Usage_All_ActivityInputs) => LocalizedString;
export const credit_usage_type_purchase: (inputs: Credit_Usage_Type_PurchaseInputs) => LocalizedString;
export const credit_usage_type_debit: (inputs: Credit_Usage_Type_DebitInputs) => LocalizedString;
export const credit_usage_no_usage_found: (inputs: Credit_Usage_No_Usage_FoundInputs) => LocalizedString;
export const credit_usage_no_usage_yet: (inputs: Credit_Usage_No_Usage_YetInputs) => LocalizedString;
export const credit_usage_no_usage_match: (inputs: Credit_Usage_No_Usage_MatchInputs) => LocalizedString;
export const credit_usage_empty_body: (inputs: Credit_Usage_Empty_BodyInputs) => LocalizedString;
export const credit_usage_showing_one: (inputs: Credit_Usage_Showing_OneInputs) => LocalizedString;
export const credit_usage_showing_other: (inputs: Credit_Usage_Showing_OtherInputs) => LocalizedString;
export const credit_usage_none_to_show: (inputs: Credit_Usage_None_To_ShowInputs) => LocalizedString;
export const credit_usage_sort_created_ascending: (inputs: Credit_Usage_Sort_Created_AscendingInputs) => LocalizedString;
export const credit_usage_sort_created_descending: (inputs: Credit_Usage_Sort_Created_DescendingInputs) => LocalizedString;
export const account_settings_title: (inputs: Account_Settings_TitleInputs) => LocalizedString;
export const account_settings_description: (inputs: Account_Settings_DescriptionInputs) => LocalizedString;
export const account_settings_nav_label: (inputs: Account_Settings_Nav_LabelInputs) => LocalizedString;
export const account_settings_account_fallback: (inputs: Account_Settings_Account_FallbackInputs) => LocalizedString;
export const account_settings_no_email_address: (inputs: Account_Settings_No_Email_AddressInputs) => LocalizedString;
export const account_settings_general: (inputs: Account_Settings_GeneralInputs) => LocalizedString;
export const account_settings_security: (inputs: Account_Settings_SecurityInputs) => LocalizedString;
export const account_settings_sessions: (inputs: Account_Settings_SessionsInputs) => LocalizedString;
export const account_settings_linked_accounts: (inputs: Account_Settings_Linked_AccountsInputs) => LocalizedString;
export const account_settings_update_error: (inputs: Account_Settings_Update_ErrorInputs) => LocalizedString;
export const account_settings_save_error: (inputs: Account_Settings_Save_ErrorInputs) => LocalizedString;
export const account_settings_revoke_session_title: (inputs: Account_Settings_Revoke_Session_TitleInputs) => LocalizedString;
export const account_settings_revoke_session_description: (inputs: Account_Settings_Revoke_Session_DescriptionInputs) => LocalizedString;
export const account_settings_revoke: (inputs: Account_Settings_RevokeInputs) => LocalizedString;
export const account_settings_session_revoked: (inputs: Account_Settings_Session_RevokedInputs) => LocalizedString;
export const account_settings_unlink_provider_title: (inputs: Account_Settings_Unlink_Provider_TitleInputs) => LocalizedString;
export const account_settings_unlink_provider_description: (inputs: Account_Settings_Unlink_Provider_DescriptionInputs) => LocalizedString;
export const account_settings_unlink: (inputs: Account_Settings_UnlinkInputs) => LocalizedString;
export const account_settings_linked_account_removed: (inputs: Account_Settings_Linked_Account_RemovedInputs) => LocalizedString;
export const account_settings_avatar_saved: (inputs: Account_Settings_Avatar_SavedInputs) => LocalizedString;
export const account_settings_name_saved: (inputs: Account_Settings_Name_SavedInputs) => LocalizedString;
export const account_settings_email_saved: (inputs: Account_Settings_Email_SavedInputs) => LocalizedString;
export const account_settings_language_saved: (inputs: Account_Settings_Language_SavedInputs) => LocalizedString;
export const account_settings_password_updated: (inputs: Account_Settings_Password_UpdatedInputs) => LocalizedString;
export const account_settings_current_session: (inputs: Account_Settings_Current_SessionInputs) => LocalizedString;
export const account_settings_browser_session: (inputs: Account_Settings_Browser_SessionInputs) => LocalizedString;
export const account_settings_session_created_at: (inputs: Account_Settings_Session_Created_AtInputs) => LocalizedString;
export const account_settings_session_ip_created_at: (inputs: Account_Settings_Session_Ip_Created_AtInputs) => LocalizedString;
export const account_settings_unknown: (inputs: Account_Settings_UnknownInputs) => LocalizedString;
export const account_settings_avatar: (inputs: Account_Settings_AvatarInputs) => LocalizedString;
export const account_settings_avatar_description: (inputs: Account_Settings_Avatar_DescriptionInputs) => LocalizedString;
export const account_settings_avatar_uploading: (inputs: Account_Settings_Avatar_UploadingInputs) => LocalizedString;
export const account_settings_avatar_upload: (inputs: Account_Settings_Avatar_UploadInputs) => LocalizedString;
export const account_settings_avatar_file_hint: (inputs: Account_Settings_Avatar_File_HintInputs) => LocalizedString;
export const account_settings_crop_avatar: (inputs: Account_Settings_Crop_AvatarInputs) => LocalizedString;
export const account_settings_crop_avatar_description: (inputs: Account_Settings_Crop_Avatar_DescriptionInputs) => LocalizedString;
export const account_settings_display_name: (inputs: Account_Settings_Display_NameInputs) => LocalizedString;
export const account_settings_email_address: (inputs: Account_Settings_Email_AddressInputs) => LocalizedString;
export const account_settings_language: (inputs: Account_Settings_LanguageInputs) => LocalizedString;
export const account_settings_save_name: (inputs: Account_Settings_Save_NameInputs) => LocalizedString;
export const account_settings_save_email: (inputs: Account_Settings_Save_EmailInputs) => LocalizedString;
export const account_settings_save_language: (inputs: Account_Settings_Save_LanguageInputs) => LocalizedString;
export const account_settings_save_password: (inputs: Account_Settings_Save_PasswordInputs) => LocalizedString;
export const account_settings_new_password: (inputs: Account_Settings_New_PasswordInputs) => LocalizedString;
export const account_settings_confirm_password: (inputs: Account_Settings_Confirm_PasswordInputs) => LocalizedString;
export const account_settings_security_description: (inputs: Account_Settings_Security_DescriptionInputs) => LocalizedString;
export const account_settings_sessions_description: (inputs: Account_Settings_Sessions_DescriptionInputs) => LocalizedString;
export const account_settings_loading_sessions: (inputs: Account_Settings_Loading_SessionsInputs) => LocalizedString;
export const account_settings_no_sessions: (inputs: Account_Settings_No_SessionsInputs) => LocalizedString;
export const account_settings_current: (inputs: Account_Settings_CurrentInputs) => LocalizedString;
export const account_settings_expires: (inputs: Account_Settings_ExpiresInputs) => LocalizedString;
export const account_settings_current_session_cannot_revoke: (inputs: Account_Settings_Current_Session_Cannot_RevokeInputs) => LocalizedString;
export const account_settings_revoke_session: (inputs: Account_Settings_Revoke_SessionInputs) => LocalizedString;
export const account_settings_revoking: (inputs: Account_Settings_RevokingInputs) => LocalizedString;
export const account_settings_linked_accounts_description: (inputs: Account_Settings_Linked_Accounts_DescriptionInputs) => LocalizedString;
export const account_settings_loading_linked_accounts: (inputs: Account_Settings_Loading_Linked_AccountsInputs) => LocalizedString;
export const account_settings_no_sign_in_methods: (inputs: Account_Settings_No_Sign_In_MethodsInputs) => LocalizedString;
export const account_settings_email_password: (inputs: Account_Settings_Email_PasswordInputs) => LocalizedString;
export const account_settings_password_enabled: (inputs: Account_Settings_Password_EnabledInputs) => LocalizedString;
export const account_settings_add_password: (inputs: Account_Settings_Add_PasswordInputs) => LocalizedString;
export const account_settings_set_password: (inputs: Account_Settings_Set_PasswordInputs) => LocalizedString;
export const account_settings_provider_google_description: (inputs: Account_Settings_Provider_Google_DescriptionInputs) => LocalizedString;
export const account_settings_provider_github_description: (inputs: Account_Settings_Provider_Github_DescriptionInputs) => LocalizedString;
export const account_settings_linked_at: (inputs: Account_Settings_Linked_AtInputs) => LocalizedString;
export const account_settings_unlinking: (inputs: Account_Settings_UnlinkingInputs) => LocalizedString;
export const account_settings_unavailable_title: (inputs: Account_Settings_Unavailable_TitleInputs) => LocalizedString;
export const account_settings_unavailable_body: (inputs: Account_Settings_Unavailable_BodyInputs) => LocalizedString;
export const billing_profile_title: (inputs: Billing_Profile_TitleInputs) => LocalizedString;
export const billing_profile_description: (inputs: Billing_Profile_DescriptionInputs) => LocalizedString;
export const billing_profile_load_error: (inputs: Billing_Profile_Load_ErrorInputs) => LocalizedString;
export const billing_profile_save_error: (inputs: Billing_Profile_Save_ErrorInputs) => LocalizedString;
export const billing_profile_saved: (inputs: Billing_Profile_SavedInputs) => LocalizedString;
export const billing_profile_company_name: (inputs: Billing_Profile_Company_NameInputs) => LocalizedString;
export const billing_profile_full_name: (inputs: Billing_Profile_Full_NameInputs) => LocalizedString;
export const billing_profile_error_title: (inputs: Billing_Profile_Error_TitleInputs) => LocalizedString;
export const billing_profile_loading: (inputs: Billing_Profile_LoadingInputs) => LocalizedString;
export const billing_profile_loading_body: (inputs: Billing_Profile_Loading_BodyInputs) => LocalizedString;
export const billing_profile_failed_load: (inputs: Billing_Profile_Failed_LoadInputs) => LocalizedString;
export const billing_profile_retry_loading: (inputs: Billing_Profile_Retry_LoadingInputs) => LocalizedString;
export const billing_profile_billing_entity: (inputs: Billing_Profile_Billing_EntityInputs) => LocalizedString;
export const billing_profile_entity_description: (inputs: Billing_Profile_Entity_DescriptionInputs) => LocalizedString;
export const billing_profile_individual: (inputs: Billing_Profile_IndividualInputs) => LocalizedString;
export const billing_profile_company: (inputs: Billing_Profile_CompanyInputs) => LocalizedString;
export const billing_profile_general_details: (inputs: Billing_Profile_General_DetailsInputs) => LocalizedString;
export const billing_profile_billing_email: (inputs: Billing_Profile_Billing_EmailInputs) => LocalizedString;
export const billing_profile_billing_address: (inputs: Billing_Profile_Billing_AddressInputs) => LocalizedString;
export const billing_profile_address_line1: (inputs: Billing_Profile_Address_Line1Inputs) => LocalizedString;
export const billing_profile_address_line2: (inputs: Billing_Profile_Address_Line2Inputs) => LocalizedString;
export const billing_profile_city: (inputs: Billing_Profile_CityInputs) => LocalizedString;
export const billing_profile_region_state: (inputs: Billing_Profile_Region_StateInputs) => LocalizedString;
export const billing_profile_country: (inputs: Billing_Profile_CountryInputs) => LocalizedString;
export const billing_profile_postal_code: (inputs: Billing_Profile_Postal_CodeInputs) => LocalizedString;
export const billing_profile_company_details: (inputs: Billing_Profile_Company_DetailsInputs) => LocalizedString;
export const billing_profile_fiscal_code: (inputs: Billing_Profile_Fiscal_CodeInputs) => LocalizedString;
export const billing_profile_registration_number: (inputs: Billing_Profile_Registration_NumberInputs) => LocalizedString;
export const billing_profile_save_button: (inputs: Billing_Profile_Save_ButtonInputs) => LocalizedString;
export const datasets_page_title: (inputs: Datasets_Page_TitleInputs) => LocalizedString;
export const datasets_detail_page_title: (inputs: Datasets_Detail_Page_TitleInputs) => LocalizedString;
export const datasets_name_column: (inputs: Datasets_Name_ColumnInputs) => LocalizedString;
export const datasets_schema_column: (inputs: Datasets_Schema_ColumnInputs) => LocalizedString;
export const datasets_fields_column: (inputs: Datasets_Fields_ColumnInputs) => LocalizedString;
export const datasets_created_column: (inputs: Datasets_Created_ColumnInputs) => LocalizedString;
export const datasets_actions_column: (inputs: Datasets_Actions_ColumnInputs) => LocalizedString;
export const datasets_sort_created_ascending: (inputs: Datasets_Sort_Created_AscendingInputs) => LocalizedString;
export const datasets_sort_created_descending: (inputs: Datasets_Sort_Created_DescendingInputs) => LocalizedString;
export const datasets_retry: (inputs: Datasets_RetryInputs) => LocalizedString;
export const datasets_open: (inputs: Datasets_OpenInputs) => LocalizedString;
export const datasets_no_datasets_found: (inputs: Datasets_No_Datasets_FoundInputs) => LocalizedString;
export const datasets_showing_datasets_one: (inputs: Datasets_Showing_Datasets_OneInputs) => LocalizedString;
export const datasets_showing_datasets_other: (inputs: Datasets_Showing_Datasets_OtherInputs) => LocalizedString;
export const datasets_no_datasets_to_show: (inputs: Datasets_No_Datasets_To_ShowInputs) => LocalizedString;
export const datasets_rows_per_page: (inputs: Datasets_Rows_Per_PageInputs) => LocalizedString;
export const datasets_previous_page: (inputs: Datasets_Previous_PageInputs) => LocalizedString;
export const datasets_next_page: (inputs: Datasets_Next_PageInputs) => LocalizedString;
export const datasets_field_count_one: (inputs: Datasets_Field_Count_OneInputs) => LocalizedString;
export const datasets_field_count_other: (inputs: Datasets_Field_Count_OtherInputs) => LocalizedString;
export const datasets_date_range: (inputs: Datasets_Date_RangeInputs) => LocalizedString;
export const datasets_any_date: (inputs: Datasets_Any_DateInputs) => LocalizedString;
export const datasets_date_range_value: (inputs: Datasets_Date_Range_ValueInputs) => LocalizedString;
export const datasets_presets: (inputs: Datasets_PresetsInputs) => LocalizedString;
export const datasets_today: (inputs: Datasets_TodayInputs) => LocalizedString;
export const datasets_this_week: (inputs: Datasets_This_WeekInputs) => LocalizedString;
export const datasets_this_month: (inputs: Datasets_This_MonthInputs) => LocalizedString;
export const datasets_clear: (inputs: Datasets_ClearInputs) => LocalizedString;
export const datasets_apply: (inputs: Datasets_ApplyInputs) => LocalizedString;
export const datasets_document_id_column: (inputs: Datasets_Document_Id_ColumnInputs) => LocalizedString;
export const datasets_filename_column: (inputs: Datasets_Filename_ColumnInputs) => LocalizedString;
export const datasets_not_found_title: (inputs: Datasets_Not_Found_TitleInputs) => LocalizedString;
export const datasets_not_found_body: (inputs: Datasets_Not_Found_BodyInputs) => LocalizedString;
export const datasets_view_datasets: (inputs: Datasets_View_DatasetsInputs) => LocalizedString;
export const datasets_preview_document: (inputs: Datasets_Preview_DocumentInputs) => LocalizedString;
export const datasets_no_documents_extracted: (inputs: Datasets_No_Documents_ExtractedInputs) => LocalizedString;
export const datasets_showing_rows_one: (inputs: Datasets_Showing_Rows_OneInputs) => LocalizedString;
export const datasets_showing_rows_other: (inputs: Datasets_Showing_Rows_OtherInputs) => LocalizedString;
export const datasets_no_rows_to_show: (inputs: Datasets_No_Rows_To_ShowInputs) => LocalizedString;
export const datasets_export_csv: (inputs: Datasets_Export_CsvInputs) => LocalizedString;
export const datasets_export_xlsx: (inputs: Datasets_Export_XlsxInputs) => LocalizedString;
export const datasets_failed_export: (inputs: Datasets_Failed_ExportInputs) => LocalizedString;
export const datasets_invalid_date: (inputs: Datasets_Invalid_DateInputs) => LocalizedString;
export const datasets_missing_document_id: (inputs: Datasets_Missing_Document_IdInputs) => LocalizedString;
export const datasets_add_dataset: (inputs: Datasets_Add_DatasetInputs) => LocalizedString;
export const datasets_all_datasets: (inputs: Datasets_All_DatasetsInputs) => LocalizedString;
export const datasets_retry_datasets: (inputs: Datasets_Retry_DatasetsInputs) => LocalizedString;
export const datasets_no_datasets: (inputs: Datasets_No_DatasetsInputs) => LocalizedString;
export const datasets_dataset_actions: (inputs: Datasets_Dataset_ActionsInputs) => LocalizedString;
export const datasets_edit: (inputs: Datasets_EditInputs) => LocalizedString;
export const datasets_delete: (inputs: Datasets_DeleteInputs) => LocalizedString;
export const datasets_delete_failed: (inputs: Datasets_Delete_FailedInputs) => LocalizedString;
export const datasets_delete_confirm_title: (inputs: Datasets_Delete_Confirm_TitleInputs) => LocalizedString;
export const datasets_delete_confirm_description: (inputs: Datasets_Delete_Confirm_DescriptionInputs) => LocalizedString;
export const datasets_dialog_title_new: (inputs: Datasets_Dialog_Title_NewInputs) => LocalizedString;
export const datasets_dialog_title_edit: (inputs: Datasets_Dialog_Title_EditInputs) => LocalizedString;
export const datasets_save_changes: (inputs: Datasets_Save_ChangesInputs) => LocalizedString;
export const datasets_create_dataset: (inputs: Datasets_Create_DatasetInputs) => LocalizedString;
export const datasets_selected_schema: (inputs: Datasets_Selected_SchemaInputs) => LocalizedString;
export const datasets_loading_schemas: (inputs: Datasets_Loading_SchemasInputs) => LocalizedString;
export const datasets_select_schema: (inputs: Datasets_Select_SchemaInputs) => LocalizedString;
export const datasets_no_fields_selected: (inputs: Datasets_No_Fields_SelectedInputs) => LocalizedString;
export const datasets_one_field_selected: (inputs: Datasets_One_Field_SelectedInputs) => LocalizedString;
export const datasets_fields_selected: (inputs: Datasets_Fields_SelectedInputs) => LocalizedString;
export const datasets_collapse_field: (inputs: Datasets_Collapse_FieldInputs) => LocalizedString;
export const datasets_expand_field: (inputs: Datasets_Expand_FieldInputs) => LocalizedString;
export const datasets_select_field: (inputs: Datasets_Select_FieldInputs) => LocalizedString;
export const datasets_name_placeholder: (inputs: Datasets_Name_PlaceholderInputs) => LocalizedString;
export const datasets_search_schemas: (inputs: Datasets_Search_SchemasInputs) => LocalizedString;
export const datasets_no_schemas_found: (inputs: Datasets_No_Schemas_FoundInputs) => LocalizedString;
export const datasets_no_fields: (inputs: Datasets_No_FieldsInputs) => LocalizedString;
export const datasets_cancel: (inputs: Datasets_CancelInputs) => LocalizedString;
export const datasets_json_badge: (inputs: Datasets_Json_BadgeInputs) => LocalizedString;
export const documents_page_title: (inputs: Documents_Page_TitleInputs) => LocalizedString;
export const documents_new_ocr_job: (inputs: Documents_New_Ocr_JobInputs) => LocalizedString;
export const documents_search_filename_placeholder: (inputs: Documents_Search_Filename_PlaceholderInputs) => LocalizedString;
export const documents_search_filename: (inputs: Documents_Search_FilenameInputs) => LocalizedString;
export const documents_date_range: (inputs: Documents_Date_RangeInputs) => LocalizedString;
export const documents_any_date: (inputs: Documents_Any_DateInputs) => LocalizedString;
export const documents_date_range_value: (inputs: Documents_Date_Range_ValueInputs) => LocalizedString;
export const documents_presets: (inputs: Documents_PresetsInputs) => LocalizedString;
export const documents_today: (inputs: Documents_TodayInputs) => LocalizedString;
export const documents_this_week: (inputs: Documents_This_WeekInputs) => LocalizedString;
export const documents_this_month: (inputs: Documents_This_MonthInputs) => LocalizedString;
export const documents_clear: (inputs: Documents_ClearInputs) => LocalizedString;
export const documents_apply: (inputs: Documents_ApplyInputs) => LocalizedString;
export const documents_filter_by_collection: (inputs: Documents_Filter_By_CollectionInputs) => LocalizedString;
export const documents_filter_by_schema: (inputs: Documents_Filter_By_SchemaInputs) => LocalizedString;
export const documents_unknown_collection: (inputs: Documents_Unknown_CollectionInputs) => LocalizedString;
export const documents_all_collections: (inputs: Documents_All_CollectionsInputs) => LocalizedString;
export const documents_all_schemas: (inputs: Documents_All_SchemasInputs) => LocalizedString;
export const documents_missing_document_id: (inputs: Documents_Missing_Document_IdInputs) => LocalizedString;
export const documents_failed_load_documents: (inputs: Documents_Failed_Load_DocumentsInputs) => LocalizedString;
export const documents_failed_load_document: (inputs: Documents_Failed_Load_DocumentInputs) => LocalizedString;
export const documents_failed_delete_document: (inputs: Documents_Failed_Delete_DocumentInputs) => LocalizedString;
export const documents_failed_update_document: (inputs: Documents_Failed_Update_DocumentInputs) => LocalizedString;
export const documents_failed_delete_documents: (inputs: Documents_Failed_Delete_DocumentsInputs) => LocalizedString;
export const documents_failed_move_documents: (inputs: Documents_Failed_Move_DocumentsInputs) => LocalizedString;
export const documents_failed_download: (inputs: Documents_Failed_DownloadInputs) => LocalizedString;
export const documents_invalid_date: (inputs: Documents_Invalid_DateInputs) => LocalizedString;
export const documents_select_all_on_page: (inputs: Documents_Select_All_On_PageInputs) => LocalizedString;
export const documents_select_document: (inputs: Documents_Select_DocumentInputs) => LocalizedString;
export const documents_filename_column: (inputs: Documents_Filename_ColumnInputs) => LocalizedString;
export const documents_collections_column: (inputs: Documents_Collections_ColumnInputs) => LocalizedString;
export const documents_pages_column: (inputs: Documents_Pages_ColumnInputs) => LocalizedString;
export const documents_created_column: (inputs: Documents_Created_ColumnInputs) => LocalizedString;
export const documents_file_size_column: (inputs: Documents_File_Size_ColumnInputs) => LocalizedString;
export const documents_sort_created_ascending: (inputs: Documents_Sort_Created_AscendingInputs) => LocalizedString;
export const documents_sort_created_descending: (inputs: Documents_Sort_Created_DescendingInputs) => LocalizedString;
export const documents_collection_not_found_title: (inputs: Documents_Collection_Not_Found_TitleInputs) => LocalizedString;
export const documents_collection_not_found_body: (inputs: Documents_Collection_Not_Found_BodyInputs) => LocalizedString;
export const documents_view_all_documents: (inputs: Documents_View_All_DocumentsInputs) => LocalizedString;
export const documents_retry: (inputs: Documents_RetryInputs) => LocalizedString;
export const documents_no_documents_found: (inputs: Documents_No_Documents_FoundInputs) => LocalizedString;
export const documents_empty_filtered_body: (inputs: Documents_Empty_Filtered_BodyInputs) => LocalizedString;
export const documents_empty_unfiltered_body: (inputs: Documents_Empty_Unfiltered_BodyInputs) => LocalizedString;
export const documents_clear_filters: (inputs: Documents_Clear_FiltersInputs) => LocalizedString;
export const documents_process_first_document: (inputs: Documents_Process_First_DocumentInputs) => LocalizedString;
export const documents_showing_documents_one: (inputs: Documents_Showing_Documents_OneInputs) => LocalizedString;
export const documents_showing_documents_other: (inputs: Documents_Showing_Documents_OtherInputs) => LocalizedString;
export const documents_no_documents_to_show: (inputs: Documents_No_Documents_To_ShowInputs) => LocalizedString;
export const documents_rows_per_page: (inputs: Documents_Rows_Per_PageInputs) => LocalizedString;
export const documents_previous: (inputs: Documents_PreviousInputs) => LocalizedString;
export const documents_next: (inputs: Documents_NextInputs) => LocalizedString;
export const documents_delete: (inputs: Documents_DeleteInputs) => LocalizedString;
export const documents_delete_single_title: (inputs: Documents_Delete_Single_TitleInputs) => LocalizedString;
export const documents_delete_single_description: (inputs: Documents_Delete_Single_DescriptionInputs) => LocalizedString;
export const documents_delete_bulk_title_one: (inputs: Documents_Delete_Bulk_Title_OneInputs) => LocalizedString;
export const documents_delete_bulk_title_other: (inputs: Documents_Delete_Bulk_Title_OtherInputs) => LocalizedString;
export const documents_delete_bulk_description_one: (inputs: Documents_Delete_Bulk_Description_OneInputs) => LocalizedString;
export const documents_delete_bulk_description_other: (inputs: Documents_Delete_Bulk_Description_OtherInputs) => LocalizedString;
export const documents_selected_count_one: (inputs: Documents_Selected_Count_OneInputs) => LocalizedString;
export const documents_selected_count_other: (inputs: Documents_Selected_Count_OtherInputs) => LocalizedString;
export const documents_download_selected: (inputs: Documents_Download_SelectedInputs) => LocalizedString;
export const documents_download: (inputs: Documents_DownloadInputs) => LocalizedString;
export const documents_downloading: (inputs: Documents_DownloadingInputs) => LocalizedString;
export const documents_move: (inputs: Documents_MoveInputs) => LocalizedString;
export const documents_moving: (inputs: Documents_MovingInputs) => LocalizedString;
export const documents_deleting: (inputs: Documents_DeletingInputs) => LocalizedString;
export const documents_open_actions_for: (inputs: Documents_Open_Actions_ForInputs) => LocalizedString;
export const documents_preview: (inputs: Documents_PreviewInputs) => LocalizedString;
export const documents_rename: (inputs: Documents_RenameInputs) => LocalizedString;
export const documents_failed_rename: (inputs: Documents_Failed_RenameInputs) => LocalizedString;
export const documents_rename_file: (inputs: Documents_Rename_FileInputs) => LocalizedString;
export const documents_preview_file: (inputs: Documents_Preview_FileInputs) => LocalizedString;
export const documents_download_dialog_title_one: (inputs: Documents_Download_Dialog_Title_OneInputs) => LocalizedString;
export const documents_download_dialog_title_other: (inputs: Documents_Download_Dialog_Title_OtherInputs) => LocalizedString;
export const documents_selected_documents: (inputs: Documents_Selected_DocumentsInputs) => LocalizedString;
export const documents_format_markdown: (inputs: Documents_Format_MarkdownInputs) => LocalizedString;
export const documents_format_html: (inputs: Documents_Format_HtmlInputs) => LocalizedString;
export const documents_format_json: (inputs: Documents_Format_JsonInputs) => LocalizedString;
export const documents_preparing_download: (inputs: Documents_Preparing_DownloadInputs) => LocalizedString;
export const documents_no_collections_selected: (inputs: Documents_No_Collections_SelectedInputs) => LocalizedString;
export const documents_one_collection_selected: (inputs: Documents_One_Collection_SelectedInputs) => LocalizedString;
export const documents_collections_selected: (inputs: Documents_Collections_SelectedInputs) => LocalizedString;
export const documents_remove_from_all: (inputs: Documents_Remove_From_AllInputs) => LocalizedString;
export const documents_move_documents: (inputs: Documents_Move_DocumentsInputs) => LocalizedString;
export const documents_move_description_one: (inputs: Documents_Move_Description_OneInputs) => LocalizedString;
export const documents_move_description_other: (inputs: Documents_Move_Description_OtherInputs) => LocalizedString;
export const documents_collections_label: (inputs: Documents_Collections_LabelInputs) => LocalizedString;
export const documents_search_collections: (inputs: Documents_Search_CollectionsInputs) => LocalizedString;
export const documents_loading_collections: (inputs: Documents_Loading_CollectionsInputs) => LocalizedString;
export const documents_no_collections_found: (inputs: Documents_No_Collections_FoundInputs) => LocalizedString;
export const documents_cancel: (inputs: Documents_CancelInputs) => LocalizedString;
export const documents_collections_nav_label: (inputs: Documents_Collections_Nav_LabelInputs) => LocalizedString;
export const documents_add_collection: (inputs: Documents_Add_CollectionInputs) => LocalizedString;
export const documents_all_documents: (inputs: Documents_All_DocumentsInputs) => LocalizedString;
export const documents_retry_collections: (inputs: Documents_Retry_CollectionsInputs) => LocalizedString;
export const documents_no_collections: (inputs: Documents_No_CollectionsInputs) => LocalizedString;
export const documents_collection_actions: (inputs: Documents_Collection_ActionsInputs) => LocalizedString;
export const documents_edit: (inputs: Documents_EditInputs) => LocalizedString;
export const documents_delete_failed: (inputs: Documents_Delete_FailedInputs) => LocalizedString;
export const documents_delete_collection_title: (inputs: Documents_Delete_Collection_TitleInputs) => LocalizedString;
export const documents_delete_collection_description: (inputs: Documents_Delete_Collection_DescriptionInputs) => LocalizedString;
export const documents_collection_dialog_title_new: (inputs: Documents_Collection_Dialog_Title_NewInputs) => LocalizedString;
export const documents_collection_dialog_title_edit: (inputs: Documents_Collection_Dialog_Title_EditInputs) => LocalizedString;
export const documents_collection_dialog_description_new: (inputs: Documents_Collection_Dialog_Description_NewInputs) => LocalizedString;
export const documents_collection_dialog_description_edit: (inputs: Documents_Collection_Dialog_Description_EditInputs) => LocalizedString;
export const documents_save_changes: (inputs: Documents_Save_ChangesInputs) => LocalizedString;
export const documents_create_collection: (inputs: Documents_Create_CollectionInputs) => LocalizedString;
export const documents_name_column: (inputs: Documents_Name_ColumnInputs) => LocalizedString;
export const documents_collection_name_placeholder: (inputs: Documents_Collection_Name_PlaceholderInputs) => LocalizedString;
export const documents_schemas_label: (inputs: Documents_Schemas_LabelInputs) => LocalizedString;
export const documents_no_schemas_selected: (inputs: Documents_No_Schemas_SelectedInputs) => LocalizedString;
export const documents_one_schema_selected: (inputs: Documents_One_Schema_SelectedInputs) => LocalizedString;
export const documents_schemas_selected: (inputs: Documents_Schemas_SelectedInputs) => LocalizedString;
export const documents_search_schemas: (inputs: Documents_Search_SchemasInputs) => LocalizedString;
export const documents_loading_schemas: (inputs: Documents_Loading_SchemasInputs) => LocalizedString;
export const documents_no_schemas_found: (inputs: Documents_No_Schemas_FoundInputs) => LocalizedString;
export const documents_collection_schema_hint: (inputs: Documents_Collection_Schema_HintInputs) => LocalizedString;
export const documents_preview_fallback_title: (inputs: Documents_Preview_Fallback_TitleInputs) => LocalizedString;
export const documents_preview_description: (inputs: Documents_Preview_DescriptionInputs) => LocalizedString;
export const documents_rename_document_title: (inputs: Documents_Rename_Document_TitleInputs) => LocalizedString;
export const documents_loading_document: (inputs: Documents_Loading_DocumentInputs) => LocalizedString;
export const documents_copy_markdown: (inputs: Documents_Copy_MarkdownInputs) => LocalizedString;
export const documents_copy_html: (inputs: Documents_Copy_HtmlInputs) => LocalizedString;
export const documents_copy_json: (inputs: Documents_Copy_JsonInputs) => LocalizedString;
export const documents_copied: (inputs: Documents_CopiedInputs) => LocalizedString;
export const documents_no_json_annotation: (inputs: Documents_No_Json_AnnotationInputs) => LocalizedString;
export const documents_no_markdown_content: (inputs: Documents_No_Markdown_ContentInputs) => LocalizedString;
export const documents_no_preview_available: (inputs: Documents_No_Preview_AvailableInputs) => LocalizedString;
export const documents_close: (inputs: Documents_CloseInputs) => LocalizedString;
export const documents_more: (inputs: Documents_MoreInputs) => LocalizedString;
export const documents_open: (inputs: Documents_OpenInputs) => LocalizedString;
export const documents_share: (inputs: Documents_ShareInputs) => LocalizedString;
export type LocalizedString = import("../runtime.js").LocalizedString;
export type Common_LanguageInputs = {};
export type Common_EnglishInputs = {};
export type Common_RomanianInputs = {};
export type Common_CancelInputs = {};
export type Common_DeleteInputs = {};
export type Common_RetryInputs = {};
export type Common_PreviousInputs = {};
export type Common_NextInputs = {};
export type Common_Rows_Per_PageInputs = {};
export type Common_StrictInputs = {};
export type Common_FlexibleInputs = {};
export type Common_RequiredInputs = {};
export type Common_UnknownInputs = {};
export type Common_ActionsInputs = {};
export type Common_Toggle_ThemeInputs = {};
export type Header_Credits_UnavailableInputs = {};
export type Header_CreditsInputs = {
    count: NonNullable<unknown>;
};
export type Header_Credit_Balance_UnavailableInputs = {
    message: NonNullable<unknown>;
};
export type Nav_AccountInputs = {};
export type Nav_No_Email_AddressInputs = {};
export type Nav_NotificationsInputs = {};
export type Nav_Log_OutInputs = {};
export type Nav_Logout_TitleInputs = {};
export type Nav_Logout_DescriptionInputs = {};
export type Nav_Logout_FailedInputs = {};
export type Nav_Account_LinkedInputs = {
    provider: NonNullable<unknown>;
};
export type Nav_Account_Link_ConflictInputs = {
    provider: NonNullable<unknown>;
};
export type Nav_Account_Link_DeniedInputs = {
    provider: NonNullable<unknown>;
};
export type Nav_Account_Link_Not_ConfiguredInputs = {
    provider: NonNullable<unknown>;
};
export type Nav_Account_Link_Sign_In_AgainInputs = {};
export type Nav_Account_Link_FailedInputs = {
    provider: NonNullable<unknown>;
};
export type Nav_DashboardInputs = {};
export type Nav_SchemasInputs = {};
export type Nav_New_SchemaInputs = {};
export type Nav_Edit_SchemaInputs = {};
export type Nav_JobsInputs = {};
export type Nav_New_JobInputs = {};
export type Nav_BillingInputs = {};
export type Nav_Billing_OrdersInputs = {};
export type Nav_Credit_Usage_HistoryInputs = {};
export type Nav_Developer_SettingsInputs = {};
export type Nav_Get_HelpInputs = {};
export type Nav_Quick_OcrInputs = {};
export type Nav_Create_Quick_Ocr_JobInputs = {};
export type Nav_Create_SchemaInputs = {};
export type Nav_Create_JobInputs = {};
export type Dashboard_Metric_Documents_ProcessedInputs = {};
export type Dashboard_Page_DescriptionInputs = {};
export type Dashboard_RefreshingInputs = {};
export type Dashboard_Loading_TitleInputs = {};
export type Dashboard_Loading_DescriptionInputs = {};
export type Dashboard_Warning_TitleInputs = {};
export type Dashboard_Unavailable_TitleInputs = {};
export type Dashboard_Unavailable_DefaultInputs = {};
export type Dashboard_Metric_Pages_ProcessedInputs = {};
export type Dashboard_Metric_Completion_RateInputs = {};
export type Dashboard_Metric_Credits_SpentInputs = {};
export type Dashboard_Jobs_In_Progress_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Jobs_In_Progress_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Pages_CompletedInputs = {};
export type Dashboard_Completion_SummaryInputs = {
    completed: NonNullable<unknown>;
    failed: NonNullable<unknown>;
};
export type Dashboard_Credits_Available_ShortInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Metrics_AriaInputs = {};
export type Dashboard_Documents_Processed_TitleInputs = {};
export type Dashboard_Chart_Documents_LabelInputs = {};
export type Dashboard_Select_RangeInputs = {};
export type Dashboard_Range_7dInputs = {};
export type Dashboard_Range_30dInputs = {};
export type Dashboard_Range_90dInputs = {};
export type Dashboard_Recent_Documents_TitleInputs = {};
export type Dashboard_Recent_Documents_DescriptionInputs = {};
export type Dashboard_ViewInputs = {};
export type Dashboard_No_Saved_SchemaInputs = {};
export type Dashboard_Pages_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Pages_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_No_Completed_DocumentsInputs = {};
export type Dashboard_Schema_Throughput_TitleInputs = {};
export type Dashboard_Schema_Throughput_DescriptionInputs = {};
export type Dashboard_Documents_Processed_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Documents_Processed_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_No_Schema_ThroughputInputs = {};
export type Dashboard_Datasets_TitleInputs = {};
export type Dashboard_Total_Datasets_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Total_Datasets_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Fields_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Fields_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_No_DatasetsInputs = {};
export type Dashboard_Credits_TitleInputs = {};
export type Dashboard_Credits_DescriptionInputs = {};
export type Dashboard_Low_CreditInputs = {};
export type Dashboard_Available_CreditsInputs = {};
export type Dashboard_Credits_Spent_In_RangeInputs = {};
export type Dashboard_BillingInputs = {};
export type Dashboard_Onboarding_TitleInputs = {};
export type Dashboard_Onboarding_DescriptionInputs = {};
export type Dashboard_New_Ocr_JobInputs = {};
export type Dashboard_Credits_OneInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Credits_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Dashboard_Step_SchemaInputs = {};
export type Dashboard_Step_Ocr_JobInputs = {};
export type Dashboard_Step_DatasetInputs = {};
export type Dashboard_Step_Api_KeyInputs = {};
export type Dashboard_Step_WebhookInputs = {};
export type Dashboard_Step_ReadyInputs = {};
export type Dashboard_Step_OpenInputs = {};
export type Admin_Nav_UsersInputs = {};
export type Admin_Nav_UserInputs = {};
export type Admin_Nav_InvoicesInputs = {};
export type Admin_Nav_OrdersInputs = {};
export type Admin_Nav_Json_RecipesInputs = {};
export type Admin_Nav_AdminInputs = {};
export type Admin_User_FallbackInputs = {};
export type Sidebar_SyncraInputs = {};
export type Sidebar_Syncra_AdminInputs = {};
export type Sidebar_User_SpaceInputs = {};
export type Sidebar_Admin_PortalInputs = {};
export type Sidebar_Switch_SpaceInputs = {};
export type Schemas_New_TitleInputs = {};
export type Schemas_LibraryInputs = {};
export type Schemas_New_DescriptionInputs = {};
export type Schemas_Edit_TitleInputs = {};
export type Schemas_Edit_DescriptionInputs = {};
export type Schemas_Save_SchemaInputs = {};
export type Schemas_Save_ChangesInputs = {};
export type Schemas_Saved_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Saved_Success_With_IdInputs = {
    name: NonNullable<unknown>;
    id: NonNullable<unknown>;
};
export type Schemas_Saved_FeedbackInputs = {
    name: NonNullable<unknown>;
    id: NonNullable<unknown>;
};
export type Schemas_Empty_Schema_ErrorInputs = {};
export type Schemas_Delete_Single_TitleInputs = {};
export type Schemas_Delete_Single_DescriptionInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Delete_Bulk_Title_OneInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Delete_Bulk_Title_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Delete_Bulk_Description_OneInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Delete_Bulk_Description_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Select_All_On_PageInputs = {};
export type Schemas_Select_SchemaInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Name_ColumnInputs = {};
export type Schemas_Id_ColumnInputs = {};
export type Schemas_Id_LabelInputs = {};
export type Schemas_Copy_IdInputs = {};
export type Schemas_Copy_Id_AriaInputs = {
    id: NonNullable<unknown>;
};
export type Schemas_Copy_Id_SuccessInputs = {};
export type Schemas_Copy_Id_ErrorInputs = {};
export type Schemas_Strict_Mode_ColumnInputs = {};
export type Schemas_Created_ColumnInputs = {};
export type Schemas_Updated_ColumnInputs = {};
export type Schemas_New_SchemaInputs = {};
export type Schemas_No_Schemas_FoundInputs = {};
export type Schemas_Empty_BodyInputs = {};
export type Schemas_Create_SchemaInputs = {};
export type Schemas_Showing_Schemas_OneInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Showing_Schemas_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_No_Schemas_To_ShowInputs = {};
export type Schemas_Selected_Count_OneInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_Selected_Count_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Schemas_DeletingInputs = {};
export type Schemas_No_DescriptionInputs = {};
export type Schemas_Sort_Created_AscendingInputs = {};
export type Schemas_Sort_Created_DescendingInputs = {};
export type Schemas_Edit_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Create_Job_WithInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Clone_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Delete_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Schemas_Loading_SchemaInputs = {};
export type Schemas_Not_Found_TitleInputs = {};
export type Schemas_Not_Found_BodyInputs = {};
export type Schemas_View_SchemasInputs = {};
export type Schemas_Could_Not_LoadInputs = {};
export type Schemas_Editor_BadgeInputs = {};
export type Schemas_General_SettingsInputs = {};
export type Schemas_Schema_Name_LabelInputs = {};
export type Schemas_Schema_Name_PlaceholderInputs = {};
export type Schemas_Description_LabelInputs = {};
export type Schemas_Description_PlaceholderInputs = {};
export type Schemas_Strict_ModeInputs = {};
export type Schemas_Flexible_ModeInputs = {};
export type Schemas_Strict_Mode_DescriptionInputs = {};
export type Schemas_Structure_DesignerInputs = {};
export type Schemas_Visual_Node_DesignerInputs = {};
export type Schemas_Validation_Name_RequiredInputs = {};
export type Schemas_Validation_Name_Too_LongInputs = {};
export type Schemas_Validation_Schema_ObjectInputs = {};
export type Schemas_CloneInputs = {};
export type Schemas_CloningInputs = {};
export type Schemas_SavingInputs = {};
export type Json_Recipes_TitleInputs = {};
export type Json_Recipes_DescriptionInputs = {};
export type Json_Recipes_New_RecipeInputs = {};
export type Json_Recipes_No_Recipes_FoundInputs = {};
export type Json_Recipes_Empty_BodyInputs = {};
export type Json_Recipes_LoadingInputs = {};
export type Json_Recipes_Loading_RecipeInputs = {};
export type Json_Recipes_Counter_ColumnInputs = {};
export type Json_Recipes_Created_ColumnInputs = {};
export type Json_Recipes_Updated_ColumnInputs = {};
export type Json_Recipes_Json_Fields_ColumnInputs = {};
export type Json_Recipes_Sort_Created_AscendingInputs = {};
export type Json_Recipes_Sort_Created_DescendingInputs = {};
export type Json_Recipes_Showing_OneInputs = {
    count: NonNullable<unknown>;
};
export type Json_Recipes_Showing_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Json_Recipes_No_Recipes_To_ShowInputs = {};
export type Json_Recipes_Edit_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipes_Delete_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipes_New_TitleInputs = {};
export type Json_Recipes_New_DescriptionInputs = {};
export type Json_Recipes_Edit_TitleInputs = {};
export type Json_Recipes_Edit_DescriptionInputs = {};
export type Json_Recipes_Save_RecipeInputs = {};
export type Json_Recipes_Save_ChangesInputs = {};
export type Json_Recipes_Created_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipes_Saved_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipes_Deleted_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipes_Delete_ConfirmInputs = {};
export type Json_Recipes_Not_Found_TitleInputs = {};
export type Json_Recipes_Not_Found_BodyInputs = {};
export type Json_Recipes_View_RecipesInputs = {};
export type Json_Recipes_Could_Not_LoadInputs = {};
export type Json_Recipes_Editor_BadgeInputs = {};
export type Json_Recipes_General_SettingsInputs = {};
export type Json_Recipes_Title_LabelInputs = {};
export type Json_Recipes_Title_PlaceholderInputs = {};
export type Json_Recipes_Description_LabelInputs = {};
export type Json_Recipes_Description_PlaceholderInputs = {};
export type Json_Recipes_Structure_DesignerInputs = {};
export type Json_Recipes_Visual_Node_DesignerInputs = {};
export type Json_Recipes_Category_LabelInputs = {};
export type Json_Recipes_OthersInputs = {};
export type Json_Recipes_Manage_CategoriesInputs = {};
export type Json_Recipes_Validation_Title_RequiredInputs = {};
export type Json_Recipes_Validation_Title_Too_LongInputs = {};
export type Json_Recipes_Validation_Json_ObjectInputs = {};
export type Json_Recipes_SavingInputs = {};
export type Json_Recipes_DeletingInputs = {};
export type Json_Recipe_Categories_TitleInputs = {};
export type Json_Recipe_Categories_DescriptionInputs = {};
export type Json_Recipe_Categories_Title_En_LabelInputs = {};
export type Json_Recipe_Categories_Title_Ro_LabelInputs = {};
export type Json_Recipe_Categories_Create_CategoryInputs = {};
export type Json_Recipe_Categories_Save_CategoryInputs = {};
export type Json_Recipe_Categories_Edit_TitleInputs = {};
export type Json_Recipe_Categories_Delete_ConfirmInputs = {};
export type Json_Recipe_Categories_LoadingInputs = {};
export type Json_Recipe_Categories_Could_Not_LoadInputs = {};
export type Json_Recipe_Categories_Empty_TitleInputs = {};
export type Json_Recipe_Categories_Empty_BodyInputs = {};
export type Json_Recipe_Categories_Created_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipe_Categories_Saved_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipe_Categories_Deleted_SuccessInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipe_Categories_Validation_Titles_RequiredInputs = {};
export type Json_Recipe_Categories_Validation_Titles_Too_LongInputs = {};
export type Json_Recipe_Categories_Edit_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Json_Recipe_Categories_Delete_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Ocr_Recipes_NavInputs = {};
export type Ocr_Recipes_TitleInputs = {};
export type Ocr_Recipes_Meta_DescriptionInputs = {};
export type Ocr_Recipes_EyebrowInputs = {};
export type Ocr_Recipes_Hero_TitleInputs = {};
export type Ocr_Recipes_Hero_DescriptionInputs = {};
export type Ocr_Recipes_Search_LabelInputs = {};
export type Ocr_Recipes_Search_PlaceholderInputs = {};
export type Ocr_Recipes_Category_FilterInputs = {};
export type Ocr_Recipes_All_CategoriesInputs = {};
export type Ocr_Recipes_Sort_LabelInputs = {};
export type Ocr_Recipes_Sort_PopularInputs = {};
export type Ocr_Recipes_Sort_NewestInputs = {};
export type Ocr_Recipes_Sort_AzInputs = {};
export type Ocr_Recipes_Showing_OneInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Showing_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_No_Matches_TitleInputs = {};
export type Ocr_Recipes_No_Matches_BodyInputs = {};
export type Ocr_Recipes_OthersInputs = {};
export type Ocr_Recipes_Fields_OneInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Fields_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Required_OneInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Required_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Deploys_OneInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Deploys_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Ocr_Recipes_Json_FieldsInputs = {};
export type Ocr_Recipes_System_RecipeInputs = {};
export type Ocr_Recipes_Strict_SchemaInputs = {};
export type Ocr_Recipes_RequiredInputs = {};
export type Ocr_Recipes_Preview_JsonInputs = {};
export type Ocr_Recipes_No_FieldsInputs = {};
export type Ocr_Recipes_Clone_RecipeInputs = {};
export type Ocr_Recipes_Clone_AriaInputs = {
    name: NonNullable<unknown>;
};
export type Ocr_Recipes_Log_In_To_CloneInputs = {};
export type Ocr_Recipes_Clone_FailedInputs = {};
export type Ocr_Recipes_Load_FailedInputs = {};
export type Jobs_Page_TitleInputs = {};
export type Jobs_Missing_Schema_IdInputs = {};
export type Jobs_Missing_Job_IdInputs = {};
export type Jobs_Delete_Bulk_Title_OneInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Delete_Bulk_Title_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Delete_Bulk_Description_OneInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Delete_Bulk_Description_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Delete_Single_TitleInputs = {};
export type Jobs_Delete_Single_DescriptionInputs = {
    name: NonNullable<unknown>;
};
export type Jobs_Status_QueuedInputs = {};
export type Jobs_Status_PendingInputs = {};
export type Jobs_Status_ProcessingInputs = {};
export type Jobs_Status_CompletedInputs = {};
export type Jobs_Status_FailedInputs = {};
export type Jobs_Inline_SchemaInputs = {};
export type Jobs_No_SchemaInputs = {};
export type Jobs_SchemaInputs = {};
export type Jobs_Select_All_On_PageInputs = {};
export type Jobs_Select_JobInputs = {
    name: NonNullable<unknown>;
};
export type Jobs_Filename_ColumnInputs = {};
export type Jobs_Status_ColumnInputs = {};
export type Jobs_Created_ColumnInputs = {};
export type Jobs_File_Size_ColumnInputs = {};
export type Jobs_Pages_ColumnInputs = {};
export type Jobs_New_JobInputs = {};
export type Jobs_No_Jobs_FoundInputs = {};
export type Jobs_Empty_BodyInputs = {};
export type Jobs_Showing_Jobs_OneInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Showing_Jobs_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_No_Jobs_To_ShowInputs = {};
export type Jobs_Selected_Count_OneInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_Selected_Count_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Jobs_DeletingInputs = {};
export type Jobs_Delete_JobInputs = {
    name: NonNullable<unknown>;
};
export type Jobs_Saved_Extraction_SchemaInputs = {};
export type Jobs_Inline_Schema_DescriptionInputs = {};
export type Jobs_Extraction_Schema_DetailsInputs = {};
export type New_Job_Missing_Document_IdInputs = {};
export type New_Job_Failed_CreateInputs = {};
export type New_Job_Insufficient_Credits_BuyInputs = {};
export type New_Job_Failed_Load_DocumentInputs = {};
export type New_Job_Invalid_Document_ResponseInputs = {};
export type New_Job_Failed_Load_SchemasInputs = {};
export type New_Job_Invalid_Schema_ResponseInputs = {};
export type New_Job_Invalid_Job_ResponseInputs = {};
export type New_Job_Failed_Load_JobInputs = {};
export type New_Job_Failed_Poll_JobInputs = {};
export type New_Job_Select_SchemaInputs = {};
export type New_Job_Select_Schema_PlaceholderInputs = {};
export type New_Job_Configure_Payload_FormatInputs = {};
export type New_Job_Upload_DocumentsInputs = {};
export type New_Job_Files_Selected_OneInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Files_Selected_OtherInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Drag_Or_Browse_FilesInputs = {};
export type New_Job_Run_MonitorInputs = {};
export type New_Job_Processing_BatchInputs = {};
export type New_Job_Start_Extraction_PipelineInputs = {};
export type New_Job_Select_Extraction_SchemaInputs = {};
export type New_Job_Select_Schema_DescriptionInputs = {};
export type New_Job_Select_Extraction_Schema_AriaInputs = {};
export type New_Job_Search_SchemasInputs = {};
export type New_Job_Loading_SchemasInputs = {};
export type New_Job_No_Schemas_FoundInputs = {};
export type New_Job_No_Schema_Ocr_OnlyInputs = {};
export type New_Job_No_Schema_DescriptionInputs = {};
export type New_Job_No_Personal_SchemasInputs = {};
export type New_Job_Create_OneInputs = {};
export type New_Job_Selected_Schema_HelpInputs = {};
export type New_Job_No_Schema_Selected_HelpInputs = {};
export type New_Job_Target_Mapped_FieldsInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_No_Fields_DefinedInputs = {};
export type New_Job_Ocr_Only_Mode_ActiveInputs = {};
export type New_Job_Ocr_Only_Mode_BodyInputs = {};
export type New_Job_Upload_Documents_DescriptionInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Dropzone_TitleInputs = {};
export type New_Job_Dropzone_DescriptionInputs = {
    size: NonNullable<unknown>;
};
export type New_Job_Browse_FilesInputs = {};
export type New_Job_Pending_Upload_QueueInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Clear_AllInputs = {};
export type New_Job_Remove_FileInputs = {};
export type New_Job_Extraction_Queue_ResultsInputs = {};
export type New_Job_File_Count_OneInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_File_Count_OtherInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_TotalInputs = {
    label: NonNullable<unknown>;
};
export type New_Job_Active_Batch_StatusInputs = {};
export type New_Job_Active_Batch_DescriptionInputs = {};
export type New_Job_ProgressInputs = {
    progress: NonNullable<unknown>;
};
export type New_Job_Total_FilesInputs = {};
export type New_Job_CompletedInputs = {};
export type New_Job_ProcessingInputs = {};
export type New_Job_FailedInputs = {};
export type New_Job_No_Active_Extraction_JobsInputs = {};
export type New_Job_No_Active_Extraction_Jobs_BodyInputs = {};
export type New_Job_Preview_DocumentInputs = {};
export type New_Job_Preview_UnavailableInputs = {};
export type New_Job_Remove_Failed_JobInputs = {};
export type New_Job_Queueing_DocumentsInputs = {};
export type New_Job_Extracting_ContentInputs = {};
export type New_Job_Run_Extraction_OneInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Run_Extraction_OtherInputs = {
    count: NonNullable<unknown>;
};
export type New_Job_Insufficient_Credits_DocumentInputs = {};
export type New_Job_Processing_FailedInputs = {};
export type New_Job_ProcessedInputs = {};
export type New_Job_Document_IdInputs = {
    id: NonNullable<unknown>;
};
export type New_Job_Creating_JobInputs = {};
export type New_Job_Queued_ProcessingInputs = {};
export type New_Job_Extracting_EntitiesInputs = {};
export type Common_ApplyInputs = {};
export type Common_ClearInputs = {};
export type Common_SavingInputs = {};
export type Common_LoadingInputs = {};
export type Common_RefreshInputs = {};
export type Common_ConnectedInputs = {};
export type Common_ConnectInputs = {};
export type Common_DownloadInputs = {};
export type Common_TodayInputs = {};
export type Common_This_WeekInputs = {};
export type Common_This_MonthInputs = {};
export type Common_AnyInputs = {};
export type Billing_UnavailableInputs = {};
export type Billing_Credit_Blocks_ErrorInputs = {};
export type Billing_Checkout_UnavailableInputs = {};
export type Billing_Payment_Received_TitleInputs = {};
export type Billing_Payment_Received_BodyInputs = {};
export type Billing_Checkout_Canceled_TitleInputs = {};
export type Billing_Checkout_Canceled_BodyInputs = {};
export type Billing_Available_BalanceInputs = {};
export type Billing_ConversionInputs = {};
export type Billing_Conversion_RateInputs = {};
export type Billing_Balance_Checked_UploadInputs = {};
export type Billing_Debited_After_SuccessInputs = {};
export type Billing_Secure_Stripe_CheckoutInputs = {};
export type Billing_Purchase_CreditsInputs = {};
export type Billing_Credits_To_PurchaseInputs = {};
export type Billing_Volume_Discount_TiersInputs = {};
export type Billing_Total_To_PayInputs = {};
export type Billing_Base_PriceInputs = {};
export type Billing_Volume_DiscountInputs = {};
export type Billing_Starting_CheckoutInputs = {};
export type Billing_Secure_CheckoutInputs = {};
export type Billing_Buy_CreditsInputs = {};
export type Billing_Orders_Page_TitleInputs = {};
export type Billing_Orders_Order_Date_FilterInputs = {};
export type Billing_Orders_Amount_ColumnInputs = {};
export type Billing_Orders_Credits_ColumnInputs = {};
export type Billing_Orders_Status_ColumnInputs = {};
export type Billing_Orders_Payment_Datetime_ColumnInputs = {};
export type Billing_Orders_Invoice_ColumnInputs = {};
export type Billing_Orders_PresetsInputs = {};
export type Billing_Orders_Filter_StatusInputs = {};
export type Billing_Orders_All_OrdersInputs = {};
export type Billing_Orders_Clear_FiltersInputs = {};
export type Billing_Orders_Clear_Filters_ActionInputs = {};
export type Billing_Orders_No_Orders_FoundInputs = {};
export type Billing_Orders_No_Orders_YetInputs = {};
export type Billing_Orders_No_Orders_MatchInputs = {};
export type Billing_Orders_Empty_BodyInputs = {};
export type Billing_Orders_Showing_OneInputs = {
    count: NonNullable<unknown>;
};
export type Billing_Orders_Showing_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Billing_Orders_None_To_ShowInputs = {};
export type Billing_Orders_Sort_Order_Date_AscendingInputs = {};
export type Billing_Orders_Sort_Order_Date_DescendingInputs = {};
export type Billing_Order_Status_PendingInputs = {};
export type Billing_Order_Status_PaidInputs = {};
export type Billing_Order_Status_FailedInputs = {};
export type Billing_Order_Status_RefundedInputs = {};
export type Billing_Order_Status_CanceledInputs = {};
export type Billing_Orders_Invoice_Pdf_TitleInputs = {
    invoice: NonNullable<unknown>;
};
export type Billing_Orders_Invoice_Preview_TitleInputs = {
    invoice: NonNullable<unknown>;
};
export type Billing_Orders_Invoice_Preview_DescriptionInputs = {};
export type Billing_Orders_Invoice_Iframe_TitleInputs = {
    invoice: NonNullable<unknown>;
};
export type Billing_Orders_Download_InvoiceInputs = {};
export type Credit_Usage_Page_TitleInputs = {};
export type Credit_Usage_Date_Range_FilterInputs = {};
export type Credit_Usage_Created_ColumnInputs = {};
export type Credit_Usage_Type_ColumnInputs = {};
export type Credit_Usage_Credits_ColumnInputs = {};
export type Credit_Usage_Related_Id_ColumnInputs = {};
export type Credit_Usage_Filter_TypeInputs = {};
export type Credit_Usage_All_ActivityInputs = {};
export type Credit_Usage_Type_PurchaseInputs = {};
export type Credit_Usage_Type_DebitInputs = {};
export type Credit_Usage_No_Usage_FoundInputs = {};
export type Credit_Usage_No_Usage_YetInputs = {};
export type Credit_Usage_No_Usage_MatchInputs = {};
export type Credit_Usage_Empty_BodyInputs = {};
export type Credit_Usage_Showing_OneInputs = {
    count: NonNullable<unknown>;
};
export type Credit_Usage_Showing_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Credit_Usage_None_To_ShowInputs = {};
export type Credit_Usage_Sort_Created_AscendingInputs = {};
export type Credit_Usage_Sort_Created_DescendingInputs = {};
export type Account_Settings_TitleInputs = {};
export type Account_Settings_DescriptionInputs = {};
export type Account_Settings_Nav_LabelInputs = {};
export type Account_Settings_Account_FallbackInputs = {};
export type Account_Settings_No_Email_AddressInputs = {};
export type Account_Settings_GeneralInputs = {};
export type Account_Settings_SecurityInputs = {};
export type Account_Settings_SessionsInputs = {};
export type Account_Settings_Linked_AccountsInputs = {};
export type Account_Settings_Update_ErrorInputs = {};
export type Account_Settings_Save_ErrorInputs = {};
export type Account_Settings_Revoke_Session_TitleInputs = {};
export type Account_Settings_Revoke_Session_DescriptionInputs = {
    session: NonNullable<unknown>;
};
export type Account_Settings_RevokeInputs = {};
export type Account_Settings_Session_RevokedInputs = {};
export type Account_Settings_Unlink_Provider_TitleInputs = {
    provider: NonNullable<unknown>;
};
export type Account_Settings_Unlink_Provider_DescriptionInputs = {
    provider: NonNullable<unknown>;
};
export type Account_Settings_UnlinkInputs = {};
export type Account_Settings_Linked_Account_RemovedInputs = {};
export type Account_Settings_Avatar_SavedInputs = {};
export type Account_Settings_Name_SavedInputs = {};
export type Account_Settings_Email_SavedInputs = {};
export type Account_Settings_Language_SavedInputs = {};
export type Account_Settings_Password_UpdatedInputs = {};
export type Account_Settings_Current_SessionInputs = {};
export type Account_Settings_Browser_SessionInputs = {};
export type Account_Settings_Session_Created_AtInputs = {
    date: NonNullable<unknown>;
};
export type Account_Settings_Session_Ip_Created_AtInputs = {
    ip: NonNullable<unknown>;
    date: NonNullable<unknown>;
};
export type Account_Settings_UnknownInputs = {};
export type Account_Settings_AvatarInputs = {};
export type Account_Settings_Avatar_DescriptionInputs = {};
export type Account_Settings_Avatar_UploadingInputs = {};
export type Account_Settings_Avatar_UploadInputs = {};
export type Account_Settings_Avatar_File_HintInputs = {};
export type Account_Settings_Crop_AvatarInputs = {};
export type Account_Settings_Crop_Avatar_DescriptionInputs = {};
export type Account_Settings_Display_NameInputs = {};
export type Account_Settings_Email_AddressInputs = {};
export type Account_Settings_LanguageInputs = {};
export type Account_Settings_Save_NameInputs = {};
export type Account_Settings_Save_EmailInputs = {};
export type Account_Settings_Save_LanguageInputs = {};
export type Account_Settings_Save_PasswordInputs = {};
export type Account_Settings_New_PasswordInputs = {};
export type Account_Settings_Confirm_PasswordInputs = {};
export type Account_Settings_Security_DescriptionInputs = {};
export type Account_Settings_Sessions_DescriptionInputs = {};
export type Account_Settings_Loading_SessionsInputs = {};
export type Account_Settings_No_SessionsInputs = {};
export type Account_Settings_CurrentInputs = {};
export type Account_Settings_ExpiresInputs = {
    date: NonNullable<unknown>;
};
export type Account_Settings_Current_Session_Cannot_RevokeInputs = {};
export type Account_Settings_Revoke_SessionInputs = {};
export type Account_Settings_RevokingInputs = {};
export type Account_Settings_Linked_Accounts_DescriptionInputs = {};
export type Account_Settings_Loading_Linked_AccountsInputs = {};
export type Account_Settings_No_Sign_In_MethodsInputs = {};
export type Account_Settings_Email_PasswordInputs = {};
export type Account_Settings_Password_EnabledInputs = {
    email: NonNullable<unknown>;
};
export type Account_Settings_Add_PasswordInputs = {};
export type Account_Settings_Set_PasswordInputs = {};
export type Account_Settings_Provider_Google_DescriptionInputs = {};
export type Account_Settings_Provider_Github_DescriptionInputs = {};
export type Account_Settings_Linked_AtInputs = {
    date: NonNullable<unknown>;
};
export type Account_Settings_UnlinkingInputs = {};
export type Account_Settings_Unavailable_TitleInputs = {};
export type Account_Settings_Unavailable_BodyInputs = {};
export type Billing_Profile_TitleInputs = {};
export type Billing_Profile_DescriptionInputs = {};
export type Billing_Profile_Load_ErrorInputs = {};
export type Billing_Profile_Save_ErrorInputs = {};
export type Billing_Profile_SavedInputs = {};
export type Billing_Profile_Company_NameInputs = {};
export type Billing_Profile_Full_NameInputs = {};
export type Billing_Profile_Error_TitleInputs = {};
export type Billing_Profile_LoadingInputs = {};
export type Billing_Profile_Loading_BodyInputs = {};
export type Billing_Profile_Failed_LoadInputs = {};
export type Billing_Profile_Retry_LoadingInputs = {};
export type Billing_Profile_Billing_EntityInputs = {};
export type Billing_Profile_Entity_DescriptionInputs = {};
export type Billing_Profile_IndividualInputs = {};
export type Billing_Profile_CompanyInputs = {};
export type Billing_Profile_General_DetailsInputs = {};
export type Billing_Profile_Billing_EmailInputs = {};
export type Billing_Profile_Billing_AddressInputs = {};
export type Billing_Profile_Address_Line1Inputs = {};
export type Billing_Profile_Address_Line2Inputs = {};
export type Billing_Profile_CityInputs = {};
export type Billing_Profile_Region_StateInputs = {};
export type Billing_Profile_CountryInputs = {};
export type Billing_Profile_Postal_CodeInputs = {};
export type Billing_Profile_Company_DetailsInputs = {};
export type Billing_Profile_Fiscal_CodeInputs = {};
export type Billing_Profile_Registration_NumberInputs = {};
export type Billing_Profile_Save_ButtonInputs = {};
export type Datasets_Page_TitleInputs = {};
export type Datasets_Detail_Page_TitleInputs = {};
export type Datasets_Name_ColumnInputs = {};
export type Datasets_Schema_ColumnInputs = {};
export type Datasets_Fields_ColumnInputs = {};
export type Datasets_Created_ColumnInputs = {};
export type Datasets_Actions_ColumnInputs = {};
export type Datasets_Sort_Created_AscendingInputs = {};
export type Datasets_Sort_Created_DescendingInputs = {};
export type Datasets_RetryInputs = {};
export type Datasets_OpenInputs = {};
export type Datasets_No_Datasets_FoundInputs = {};
export type Datasets_Showing_Datasets_OneInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_Showing_Datasets_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_No_Datasets_To_ShowInputs = {};
export type Datasets_Rows_Per_PageInputs = {};
export type Datasets_Previous_PageInputs = {};
export type Datasets_Next_PageInputs = {};
export type Datasets_Field_Count_OneInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_Field_Count_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_Date_RangeInputs = {};
export type Datasets_Any_DateInputs = {};
export type Datasets_Date_Range_ValueInputs = {
    start: NonNullable<unknown>;
    end: NonNullable<unknown>;
};
export type Datasets_PresetsInputs = {};
export type Datasets_TodayInputs = {};
export type Datasets_This_WeekInputs = {};
export type Datasets_This_MonthInputs = {};
export type Datasets_ClearInputs = {};
export type Datasets_ApplyInputs = {};
export type Datasets_Document_Id_ColumnInputs = {};
export type Datasets_Filename_ColumnInputs = {};
export type Datasets_Not_Found_TitleInputs = {};
export type Datasets_Not_Found_BodyInputs = {};
export type Datasets_View_DatasetsInputs = {};
export type Datasets_Preview_DocumentInputs = {
    documentId: NonNullable<unknown>;
};
export type Datasets_No_Documents_ExtractedInputs = {};
export type Datasets_Showing_Rows_OneInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_Showing_Rows_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_No_Rows_To_ShowInputs = {};
export type Datasets_Export_CsvInputs = {};
export type Datasets_Export_XlsxInputs = {};
export type Datasets_Failed_ExportInputs = {};
export type Datasets_Invalid_DateInputs = {};
export type Datasets_Missing_Document_IdInputs = {};
export type Datasets_Add_DatasetInputs = {};
export type Datasets_All_DatasetsInputs = {};
export type Datasets_Retry_DatasetsInputs = {};
export type Datasets_No_DatasetsInputs = {};
export type Datasets_Dataset_ActionsInputs = {};
export type Datasets_EditInputs = {};
export type Datasets_DeleteInputs = {};
export type Datasets_Delete_FailedInputs = {};
export type Datasets_Delete_Confirm_TitleInputs = {};
export type Datasets_Delete_Confirm_DescriptionInputs = {
    name: NonNullable<unknown>;
};
export type Datasets_Dialog_Title_NewInputs = {};
export type Datasets_Dialog_Title_EditInputs = {};
export type Datasets_Save_ChangesInputs = {};
export type Datasets_Create_DatasetInputs = {};
export type Datasets_Selected_SchemaInputs = {};
export type Datasets_Loading_SchemasInputs = {};
export type Datasets_Select_SchemaInputs = {};
export type Datasets_No_Fields_SelectedInputs = {};
export type Datasets_One_Field_SelectedInputs = {};
export type Datasets_Fields_SelectedInputs = {
    count: NonNullable<unknown>;
};
export type Datasets_Collapse_FieldInputs = {
    label: NonNullable<unknown>;
};
export type Datasets_Expand_FieldInputs = {
    label: NonNullable<unknown>;
};
export type Datasets_Select_FieldInputs = {
    label: NonNullable<unknown>;
};
export type Datasets_Name_PlaceholderInputs = {};
export type Datasets_Search_SchemasInputs = {};
export type Datasets_No_Schemas_FoundInputs = {};
export type Datasets_No_FieldsInputs = {};
export type Datasets_CancelInputs = {};
export type Datasets_Json_BadgeInputs = {};
export type Documents_Page_TitleInputs = {};
export type Documents_New_Ocr_JobInputs = {};
export type Documents_Search_Filename_PlaceholderInputs = {};
export type Documents_Search_FilenameInputs = {};
export type Documents_Date_RangeInputs = {};
export type Documents_Any_DateInputs = {};
export type Documents_Date_Range_ValueInputs = {
    start: NonNullable<unknown>;
    end: NonNullable<unknown>;
};
export type Documents_PresetsInputs = {};
export type Documents_TodayInputs = {};
export type Documents_This_WeekInputs = {};
export type Documents_This_MonthInputs = {};
export type Documents_ClearInputs = {};
export type Documents_ApplyInputs = {};
export type Documents_Filter_By_CollectionInputs = {};
export type Documents_Filter_By_SchemaInputs = {};
export type Documents_Unknown_CollectionInputs = {};
export type Documents_All_CollectionsInputs = {};
export type Documents_All_SchemasInputs = {};
export type Documents_Missing_Document_IdInputs = {};
export type Documents_Failed_Load_DocumentsInputs = {};
export type Documents_Failed_Load_DocumentInputs = {};
export type Documents_Failed_Delete_DocumentInputs = {};
export type Documents_Failed_Update_DocumentInputs = {};
export type Documents_Failed_Delete_DocumentsInputs = {};
export type Documents_Failed_Move_DocumentsInputs = {};
export type Documents_Failed_DownloadInputs = {};
export type Documents_Invalid_DateInputs = {};
export type Documents_Select_All_On_PageInputs = {};
export type Documents_Select_DocumentInputs = {
    name: NonNullable<unknown>;
};
export type Documents_Filename_ColumnInputs = {};
export type Documents_Collections_ColumnInputs = {};
export type Documents_Pages_ColumnInputs = {};
export type Documents_Created_ColumnInputs = {};
export type Documents_File_Size_ColumnInputs = {};
export type Documents_Sort_Created_AscendingInputs = {};
export type Documents_Sort_Created_DescendingInputs = {};
export type Documents_Collection_Not_Found_TitleInputs = {};
export type Documents_Collection_Not_Found_BodyInputs = {};
export type Documents_View_All_DocumentsInputs = {};
export type Documents_RetryInputs = {};
export type Documents_No_Documents_FoundInputs = {};
export type Documents_Empty_Filtered_BodyInputs = {};
export type Documents_Empty_Unfiltered_BodyInputs = {};
export type Documents_Clear_FiltersInputs = {};
export type Documents_Process_First_DocumentInputs = {};
export type Documents_Showing_Documents_OneInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Showing_Documents_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_No_Documents_To_ShowInputs = {};
export type Documents_Rows_Per_PageInputs = {};
export type Documents_PreviousInputs = {};
export type Documents_NextInputs = {};
export type Documents_DeleteInputs = {};
export type Documents_Delete_Single_TitleInputs = {};
export type Documents_Delete_Single_DescriptionInputs = {
    name: NonNullable<unknown>;
};
export type Documents_Delete_Bulk_Title_OneInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Delete_Bulk_Title_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Delete_Bulk_Description_OneInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Delete_Bulk_Description_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Selected_Count_OneInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Selected_Count_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Download_SelectedInputs = {};
export type Documents_DownloadInputs = {};
export type Documents_DownloadingInputs = {};
export type Documents_MoveInputs = {};
export type Documents_MovingInputs = {};
export type Documents_DeletingInputs = {};
export type Documents_Open_Actions_ForInputs = {
    name: NonNullable<unknown>;
};
export type Documents_PreviewInputs = {};
export type Documents_RenameInputs = {};
export type Documents_Failed_RenameInputs = {};
export type Documents_Rename_FileInputs = {
    name: NonNullable<unknown>;
};
export type Documents_Preview_FileInputs = {
    name: NonNullable<unknown>;
};
export type Documents_Download_Dialog_Title_OneInputs = {};
export type Documents_Download_Dialog_Title_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Selected_DocumentsInputs = {};
export type Documents_Format_MarkdownInputs = {};
export type Documents_Format_HtmlInputs = {};
export type Documents_Format_JsonInputs = {};
export type Documents_Preparing_DownloadInputs = {};
export type Documents_No_Collections_SelectedInputs = {};
export type Documents_One_Collection_SelectedInputs = {};
export type Documents_Collections_SelectedInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Remove_From_AllInputs = {};
export type Documents_Move_DocumentsInputs = {};
export type Documents_Move_Description_OneInputs = {};
export type Documents_Move_Description_OtherInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Collections_LabelInputs = {};
export type Documents_Search_CollectionsInputs = {};
export type Documents_Loading_CollectionsInputs = {};
export type Documents_No_Collections_FoundInputs = {};
export type Documents_CancelInputs = {};
export type Documents_Collections_Nav_LabelInputs = {};
export type Documents_Add_CollectionInputs = {};
export type Documents_All_DocumentsInputs = {};
export type Documents_Retry_CollectionsInputs = {};
export type Documents_No_CollectionsInputs = {};
export type Documents_Collection_ActionsInputs = {};
export type Documents_EditInputs = {};
export type Documents_Delete_FailedInputs = {};
export type Documents_Delete_Collection_TitleInputs = {};
export type Documents_Delete_Collection_DescriptionInputs = {
    name: NonNullable<unknown>;
};
export type Documents_Collection_Dialog_Title_NewInputs = {};
export type Documents_Collection_Dialog_Title_EditInputs = {};
export type Documents_Collection_Dialog_Description_NewInputs = {};
export type Documents_Collection_Dialog_Description_EditInputs = {};
export type Documents_Save_ChangesInputs = {};
export type Documents_Create_CollectionInputs = {};
export type Documents_Name_ColumnInputs = {};
export type Documents_Collection_Name_PlaceholderInputs = {};
export type Documents_Schemas_LabelInputs = {};
export type Documents_No_Schemas_SelectedInputs = {};
export type Documents_One_Schema_SelectedInputs = {};
export type Documents_Schemas_SelectedInputs = {
    count: NonNullable<unknown>;
};
export type Documents_Search_SchemasInputs = {};
export type Documents_Loading_SchemasInputs = {};
export type Documents_No_Schemas_FoundInputs = {};
export type Documents_Collection_Schema_HintInputs = {};
export type Documents_Preview_Fallback_TitleInputs = {};
export type Documents_Preview_DescriptionInputs = {};
export type Documents_Rename_Document_TitleInputs = {};
export type Documents_Loading_DocumentInputs = {};
export type Documents_Copy_MarkdownInputs = {};
export type Documents_Copy_HtmlInputs = {};
export type Documents_Copy_JsonInputs = {};
export type Documents_CopiedInputs = {};
export type Documents_No_Json_AnnotationInputs = {};
export type Documents_No_Markdown_ContentInputs = {};
export type Documents_No_Preview_AvailableInputs = {};
export type Documents_CloseInputs = {};
export type Documents_MoreInputs = {};
export type Documents_OpenInputs = {};
export type Documents_ShareInputs = {};
