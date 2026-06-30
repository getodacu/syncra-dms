/* eslint-disable */
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


export const common_language = /** @type {(inputs: Common_LanguageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Limbă`)
};

export const common_english = /** @type {(inputs: Common_EnglishInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Engleză`)
};

export const common_romanian = /** @type {(inputs: Common_RomanianInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Română`)
};

export const common_cancel = /** @type {(inputs: Common_CancelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anulează`)
};

export const common_delete = /** @type {(inputs: Common_DeleteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const common_retry = /** @type {(inputs: Common_RetryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă`)
};

export const common_previous = /** @type {(inputs: Common_PreviousInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anterior`)
};

export const common_next = /** @type {(inputs: Common_NextInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Următor`)
};

export const common_rows_per_page = /** @type {(inputs: Common_Rows_Per_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rânduri pe pagină`)
};

export const common_strict = /** @type {(inputs: Common_StrictInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Strict`)
};

export const common_flexible = /** @type {(inputs: Common_FlexibleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Flexibil`)
};

export const common_required = /** @type {(inputs: Common_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Obligatoriu`)
};

export const common_unknown = /** @type {(inputs: Common_UnknownInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Necunoscut`)
};

export const common_actions = /** @type {(inputs: Common_ActionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Acțiuni`)
};

export const common_toggle_theme = /** @type {(inputs: Common_Toggle_ThemeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schimbă tema`)
};

export const header_credits_unavailable = /** @type {(inputs: Header_Credits_UnavailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creditele nu sunt disponibile`)
};

export const header_credits = /** @type {(inputs: Header_CreditsInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} credite`)
};

export const header_credit_balance_unavailable = /** @type {(inputs: Header_Credit_Balance_UnavailableInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Soldul de credite nu este disponibil: ${i?.message}`)
};

export const nav_account = /** @type {(inputs: Nav_AccountInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cont`)
};

export const nav_no_email_address = /** @type {(inputs: Nav_No_Email_AddressInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio adresă de email`)
};

export const nav_notifications = /** @type {(inputs: Nav_NotificationsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Notificări`)
};

export const nav_log_out = /** @type {(inputs: Nav_Log_OutInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deconectare`)
};

export const nav_logout_title = /** @type {(inputs: Nav_Logout_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Te deconectezi?`)
};

export const nav_logout_description = /** @type {(inputs: Nav_Logout_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sigur vrei să te deconectezi din Syncra?`)
};

export const nav_logout_failed = /** @type {(inputs: Nav_Logout_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deconectarea a eșuat. Încearcă din nou.`)
};

export const nav_account_linked = /** @type {(inputs: Nav_Account_LinkedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Contul ${i?.provider} a fost conectat.`)
};

export const nav_account_link_conflict = /** @type {(inputs: Nav_Account_Link_ConflictInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.provider} este deja conectat la alt cont.`)
};

export const nav_account_link_denied = /** @type {(inputs: Nav_Account_Link_DeniedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Conectarea ${i?.provider} a fost anulată.`)
};

export const nav_account_link_not_configured = /** @type {(inputs: Nav_Account_Link_Not_ConfiguredInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Conectarea ${i?.provider} nu este configurată.`)
};

export const nav_account_link_sign_in_again = /** @type {(inputs: Nav_Account_Link_Sign_In_AgainInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Autentifică-te din nou înainte de a conecta conturi.`)
};

export const nav_account_link_failed = /** @type {(inputs: Nav_Account_Link_FailedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.provider} nu a putut fi conectat.`)
};

export const nav_dashboard = /** @type {(inputs: Nav_DashboardInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dashboard`)
};

export const nav_schemas = /** @type {(inputs: Nav_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Scheme`)
};

export const nav_new_schema = /** @type {(inputs: Nav_New_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă nouă`)
};

export const nav_edit_schema = /** @type {(inputs: Nav_Edit_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează schema`)
};

export const nav_jobs = /** @type {(inputs: Nav_JobsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Joburi`)
};

export const nav_new_job = /** @type {(inputs: Nav_New_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Job nou`)
};

export const nav_billing = /** @type {(inputs: Nav_BillingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Facturare`)
};

export const nav_billing_orders = /** @type {(inputs: Nav_Billing_OrdersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Comenzi de facturare`)
};

export const nav_credit_usage_history = /** @type {(inputs: Nav_Credit_Usage_HistoryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Istoric utilizare credite`)
};

export const nav_developer_settings = /** @type {(inputs: Nav_Developer_SettingsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setări dezvoltator`)
};

export const nav_get_help = /** @type {(inputs: Nav_Get_HelpInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ajutor`)
};

export const nav_quick_ocr = /** @type {(inputs: Nav_Quick_OcrInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`OCR rapid`)
};

export const nav_create_quick_ocr_job = /** @type {(inputs: Nav_Create_Quick_Ocr_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează procesare OCR rapidă`)
};

export const nav_create_schema = /** @type {(inputs: Nav_Create_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează schemă`)
};

export const nav_create_job = /** @type {(inputs: Nav_Create_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează job`)
};

export const dashboard_metric_documents_processed = /** @type {(inputs: Dashboard_Metric_Documents_ProcessedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente procesate`)
};

export const dashboard_page_description = /** @type {(inputs: Dashboard_Page_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesare, activitate recentă, dataseturi și credite într-un singur loc.`)
};

export const dashboard_refreshing = /** @type {(inputs: Dashboard_RefreshingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se reîmprospătează`)
};

export const dashboard_loading_title = /** @type {(inputs: Dashboard_Loading_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă dashboardul`)
};

export const dashboard_loading_description = /** @type {(inputs: Dashboard_Loading_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pregătim imaginea de ansamblu a spațiului tău de lucru.`)
};

export const dashboard_warning_title = /** @type {(inputs: Dashboard_Warning_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dashboard încărcat parțial`)
};

export const dashboard_unavailable_title = /** @type {(inputs: Dashboard_Unavailable_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dashboard indisponibil`)
};

export const dashboard_unavailable_default = /** @type {(inputs: Dashboard_Unavailable_DefaultInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Datele dashboardului nu au putut fi încărcate.`)
};

export const dashboard_metric_pages_processed = /** @type {(inputs: Dashboard_Metric_Pages_ProcessedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagini procesate`)
};

export const dashboard_metric_completion_rate = /** @type {(inputs: Dashboard_Metric_Completion_RateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rată de finalizare`)
};

export const dashboard_metric_credits_spent = /** @type {(inputs: Dashboard_Metric_Credits_SpentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite consumate`)
};

export const dashboard_jobs_in_progress_one = /** @type {(inputs: Dashboard_Jobs_In_Progress_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} job în desfășurare`)
};

export const dashboard_jobs_in_progress_other = /** @type {(inputs: Dashboard_Jobs_In_Progress_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} joburi în desfășurare`)
};

export const dashboard_pages_completed = /** @type {(inputs: Dashboard_Pages_CompletedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagini OCR finalizate`)
};

export const dashboard_completion_summary = /** @type {(inputs: Dashboard_Completion_SummaryInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.completed} finalizate, ${i?.failed} eșuate`)
};

export const dashboard_credits_available_short = /** @type {(inputs: Dashboard_Credits_Available_ShortInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} disponibile`)
};

export const dashboard_metrics_aria = /** @type {(inputs: Dashboard_Metrics_AriaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Metrici dashboard`)
};

export const dashboard_documents_processed_title = /** @type {(inputs: Dashboard_Documents_Processed_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente procesate`)
};

export const dashboard_chart_documents_label = /** @type {(inputs: Dashboard_Chart_Documents_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente`)
};

export const dashboard_select_range = /** @type {(inputs: Dashboard_Select_RangeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează intervalul dashboardului`)
};

export const dashboard_range_7d = /** @type {(inputs: Dashboard_Range_7dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ultimele 7 zile`)
};

export const dashboard_range_30d = /** @type {(inputs: Dashboard_Range_30dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ultimele 30 de zile`)
};

export const dashboard_range_90d = /** @type {(inputs: Dashboard_Range_90dInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ultimele 90 de zile`)
};

export const dashboard_recent_documents_title = /** @type {(inputs: Dashboard_Recent_Documents_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente recente`)
};

export const dashboard_recent_documents_description = /** @type {(inputs: Dashboard_Recent_Documents_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cele mai recente documente OCR finalizate`)
};

export const dashboard_view = /** @type {(inputs: Dashboard_ViewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Vezi`)
};

export const dashboard_no_saved_schema = /** @type {(inputs: Dashboard_No_Saved_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fără schemă salvată`)
};

export const dashboard_pages_one = /** @type {(inputs: Dashboard_Pages_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} pagină`)
};

export const dashboard_pages_other = /** @type {(inputs: Dashboard_Pages_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} pagini`)
};

export const dashboard_no_completed_documents = /** @type {(inputs: Dashboard_No_Completed_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există încă documente finalizate`)
};

export const dashboard_schema_throughput_title = /** @type {(inputs: Dashboard_Schema_Throughput_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesare pe schemă`)
};

export const dashboard_schema_throughput_description = /** @type {(inputs: Dashboard_Schema_Throughput_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente finalizate pe schemă`)
};

export const dashboard_documents_processed_one = /** @type {(inputs: Dashboard_Documents_Processed_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} document procesat`)
};

export const dashboard_documents_processed_other = /** @type {(inputs: Dashboard_Documents_Processed_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} documente procesate`)
};

export const dashboard_no_schema_throughput = /** @type {(inputs: Dashboard_No_Schema_ThroughputInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există procesare pe schemă în acest interval`)
};

export const dashboard_datasets_title = /** @type {(inputs: Dashboard_Datasets_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dataseturi`)
};

export const dashboard_total_datasets_one = /** @type {(inputs: Dashboard_Total_Datasets_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} dataset în total`)
};

export const dashboard_total_datasets_other = /** @type {(inputs: Dashboard_Total_Datasets_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} dataseturi în total`)
};

export const dashboard_fields_one = /** @type {(inputs: Dashboard_Fields_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmp`)
};

export const dashboard_fields_other = /** @type {(inputs: Dashboard_Fields_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmpuri`)
};

export const dashboard_no_datasets = /** @type {(inputs: Dashboard_No_DatasetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există încă dataseturi`)
};

export const dashboard_credits_title = /** @type {(inputs: Dashboard_Credits_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite`)
};

export const dashboard_credits_description = /** @type {(inputs: Dashboard_Credits_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sold și utilizare în interval`)
};

export const dashboard_low_credit = /** @type {(inputs: Dashboard_Low_CreditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite puține`)
};

export const dashboard_available_credits = /** @type {(inputs: Dashboard_Available_CreditsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite disponibile`)
};

export const dashboard_credits_spent_in_range = /** @type {(inputs: Dashboard_Credits_Spent_In_RangeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite consumate în intervalul selectat`)
};

export const dashboard_billing = /** @type {(inputs: Dashboard_BillingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Facturare`)
};

export const dashboard_onboarding_title = /** @type {(inputs: Dashboard_Onboarding_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Începe procesarea documentelor`)
};

export const dashboard_onboarding_description = /** @type {(inputs: Dashboard_Onboarding_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează o schemă, rulează OCR, apoi transformă rezultatele în dataseturi.`)
};

export const dashboard_new_ocr_job = /** @type {(inputs: Dashboard_New_Ocr_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Job OCR nou`)
};

export const dashboard_credits_one = /** @type {(inputs: Dashboard_Credits_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} credit`)
};

export const dashboard_credits_other = /** @type {(inputs: Dashboard_Credits_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} credite`)
};

export const dashboard_step_schema = /** @type {(inputs: Dashboard_Step_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă`)
};

export const dashboard_step_ocr_job = /** @type {(inputs: Dashboard_Step_Ocr_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Job OCR`)
};

export const dashboard_step_dataset = /** @type {(inputs: Dashboard_Step_DatasetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dataset`)
};

export const dashboard_step_api_key = /** @type {(inputs: Dashboard_Step_Api_KeyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cheie API`)
};

export const dashboard_step_webhook = /** @type {(inputs: Dashboard_Step_WebhookInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Webhook`)
};

export const dashboard_step_ready = /** @type {(inputs: Dashboard_Step_ReadyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Gata`)
};

export const dashboard_step_open = /** @type {(inputs: Dashboard_Step_OpenInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deschide`)
};

export const admin_nav_users = /** @type {(inputs: Admin_Nav_UsersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Utilizatori`)
};

export const admin_nav_user = /** @type {(inputs: Admin_Nav_UserInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Utilizator`)
};

export const admin_nav_invoices = /** @type {(inputs: Admin_Nav_InvoicesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Facturi`)
};

export const admin_nav_orders = /** @type {(inputs: Admin_Nav_OrdersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Comenzi`)
};

export const admin_nav_json_recipes = /** @type {(inputs: Admin_Nav_Json_RecipesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețete`)
};

export const admin_nav_admin = /** @type {(inputs: Admin_Nav_AdminInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Admin`)
};

export const admin_user_fallback = /** @type {(inputs: Admin_User_FallbackInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Admin`)
};

export const sidebar_syncra = /** @type {(inputs: Sidebar_SyncraInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Syncra`)
};

export const sidebar_syncra_admin = /** @type {(inputs: Sidebar_Syncra_AdminInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Syncra Admin`)
};

export const sidebar_user_space = /** @type {(inputs: Sidebar_User_SpaceInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Spațiu utilizator`)
};

export const sidebar_admin_portal = /** @type {(inputs: Sidebar_Admin_PortalInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Portal admin`)
};

export const sidebar_switch_space = /** @type {(inputs: Sidebar_Switch_SpaceInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schimbă spațiul`)
};

export const schemas_new_title = /** @type {(inputs: Schemas_New_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă nouă`)
};

export const schemas_library = /** @type {(inputs: Schemas_LibraryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Bibliotecă`)
};

export const schemas_new_description = /** @type {(inputs: Schemas_New_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Definește metadatele și structura schemei.`)
};

export const schemas_edit_title = /** @type {(inputs: Schemas_Edit_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează schema`)
};

export const schemas_edit_description = /** @type {(inputs: Schemas_Edit_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizează metadatele și structura schemei.`)
};

export const schemas_save_schema = /** @type {(inputs: Schemas_Save_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează schema`)
};

export const schemas_save_changes = /** @type {(inputs: Schemas_Save_ChangesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează modificările`)
};

export const schemas_saved_success = /** @type {(inputs: Schemas_Saved_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Schema ${i?.name} a fost salvată.`)
};

export const schemas_saved_success_with_id = /** @type {(inputs: Schemas_Saved_Success_With_IdInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Schema ${i?.name} (${i?.id}) a fost salvată.`)
};

export const schemas_saved_feedback = /** @type {(inputs: Schemas_Saved_FeedbackInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.name} (${i?.id}) salvată`)
};

export const schemas_empty_schema_error = /** @type {(inputs: Schemas_Empty_Schema_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema trebuie să includă cel puțin un câmp.`)
};

export const schemas_delete_single_title = /** @type {(inputs: Schemas_Delete_Single_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi schema?`)
};

export const schemas_delete_single_description = /** @type {(inputs: Schemas_Delete_Single_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi „${i?.name}”? Această acțiune nu poate fi anulată.`)
};

export const schemas_delete_bulk_title_one = /** @type {(inputs: Schemas_Delete_Bulk_Title_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} schemă?`)
};

export const schemas_delete_bulk_title_other = /** @type {(inputs: Schemas_Delete_Bulk_Title_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} scheme?`)
};

export const schemas_delete_bulk_description_one = /** @type {(inputs: Schemas_Delete_Bulk_Description_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} schemă selectată? Această acțiune nu poate fi anulată.`)
};

export const schemas_delete_bulk_description_other = /** @type {(inputs: Schemas_Delete_Bulk_Description_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} scheme selectate? Această acțiune nu poate fi anulată.`)
};

export const schemas_select_all_on_page = /** @type {(inputs: Schemas_Select_All_On_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează toate schemele de pe această pagină`)
};

export const schemas_select_schema = /** @type {(inputs: Schemas_Select_SchemaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Selectează ${i?.name}`)
};

export const schemas_name_column = /** @type {(inputs: Schemas_Name_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume`)
};

export const schemas_id_column = /** @type {(inputs: Schemas_Id_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID`)
};

export const schemas_id_label = /** @type {(inputs: Schemas_Id_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID schemă`)
};

export const schemas_copy_id = /** @type {(inputs: Schemas_Copy_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Copiază ID-ul`)
};

export const schemas_copy_id_aria = /** @type {(inputs: Schemas_Copy_Id_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Copiază ID-ul schemei ${i?.id}`)
};

export const schemas_copy_id_success = /** @type {(inputs: Schemas_Copy_Id_SuccessInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul schemei a fost copiat.`)
};

export const schemas_copy_id_error = /** @type {(inputs: Schemas_Copy_Id_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul schemei nu a putut fi copiat.`)
};

export const schemas_strict_mode_column = /** @type {(inputs: Schemas_Strict_Mode_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mod strict`)
};

export const schemas_created_column = /** @type {(inputs: Schemas_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creată`)
};

export const schemas_updated_column = /** @type {(inputs: Schemas_Updated_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizată`)
};

export const schemas_new_schema = /** @type {(inputs: Schemas_New_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă nouă`)
};

export const schemas_no_schemas_found = /** @type {(inputs: Schemas_No_Schemas_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit scheme`)
};

export const schemas_empty_body = /** @type {(inputs: Schemas_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează o schemă pentru a defini câmpuri structurate pentru extragerea din documente.`)
};

export const schemas_create_schema = /** @type {(inputs: Schemas_Create_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează schemă`)
};

export const schemas_showing_schemas_one = /** @type {(inputs: Schemas_Showing_Schemas_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} schemă pe această pagină.`)
};

export const schemas_showing_schemas_other = /** @type {(inputs: Schemas_Showing_Schemas_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} scheme pe această pagină.`)
};

export const schemas_no_schemas_to_show = /** @type {(inputs: Schemas_No_Schemas_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există scheme de afișat.`)
};

export const schemas_selected_count_one = /** @type {(inputs: Schemas_Selected_Count_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectată`)
};

export const schemas_selected_count_other = /** @type {(inputs: Schemas_Selected_Count_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectate`)
};

export const schemas_deleting = /** @type {(inputs: Schemas_DeletingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se șterge...`)
};

export const schemas_no_description = /** @type {(inputs: Schemas_No_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fără descriere`)
};

export const schemas_sort_created_ascending = /** @type {(inputs: Schemas_Sort_Created_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării ascendent`)
};

export const schemas_sort_created_descending = /** @type {(inputs: Schemas_Sort_Created_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării descendent`)
};

export const schemas_edit_aria = /** @type {(inputs: Schemas_Edit_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Editează ${i?.name}`)
};

export const schemas_create_job_with = /** @type {(inputs: Schemas_Create_Job_WithInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Creează job cu ${i?.name}`)
};

export const schemas_clone_aria = /** @type {(inputs: Schemas_Clone_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Clonează ${i?.name}`)
};

export const schemas_delete_aria = /** @type {(inputs: Schemas_Delete_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Șterge ${i?.name}`)
};

export const schemas_loading_schema = /** @type {(inputs: Schemas_Loading_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă schema...`)
};

export const schemas_not_found_title = /** @type {(inputs: Schemas_Not_Found_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema nu a fost găsită`)
};

export const schemas_not_found_body = /** @type {(inputs: Schemas_Not_Found_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Această schemă nu există.`)
};

export const schemas_view_schemas = /** @type {(inputs: Schemas_View_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Vezi schemele`)
};

export const schemas_could_not_load = /** @type {(inputs: Schemas_Could_Not_LoadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema nu a putut fi încărcată`)
};

export const schemas_editor_badge = /** @type {(inputs: Schemas_Editor_BadgeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editor schemă`)
};

export const schemas_general_settings = /** @type {(inputs: Schemas_General_SettingsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setări generale`)
};

export const schemas_schema_name_label = /** @type {(inputs: Schemas_Schema_Name_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume schemă`)
};

export const schemas_schema_name_placeholder = /** @type {(inputs: Schemas_Schema_Name_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume schemă`)
};

export const schemas_description_label = /** @type {(inputs: Schemas_Description_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descriere`)
};

export const schemas_description_placeholder = /** @type {(inputs: Schemas_Description_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adaugă context sau instrucțiuni opționale pentru această schemă...`)
};

export const schemas_strict_mode = /** @type {(inputs: Schemas_Strict_ModeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mod strict`)
};

export const schemas_flexible_mode = /** @type {(inputs: Schemas_Flexible_ModeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mod flexibil`)
};

export const schemas_strict_mode_description = /** @type {(inputs: Schemas_Strict_Mode_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Respinge câmpurile care nu sunt declarate explicit în această schemă. Recomandat pentru extragerea structurată de entități.`)
};

export const schemas_structure_designer = /** @type {(inputs: Schemas_Structure_DesignerInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Designer structură`)
};

export const schemas_visual_node_designer = /** @type {(inputs: Schemas_Visual_Node_DesignerInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Designer vizual noduri`)
};

export const schemas_validation_name_required = /** @type {(inputs: Schemas_Validation_Name_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Numele este obligatoriu.`)
};

export const schemas_validation_name_too_long = /** @type {(inputs: Schemas_Validation_Name_Too_LongInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Numele trebuie să aibă cel mult 160 de caractere.`)
};

export const schemas_validation_schema_object = /** @type {(inputs: Schemas_Validation_Schema_ObjectInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema trebuie să fie un obiect JSON.`)
};

export const schemas_clone = /** @type {(inputs: Schemas_CloneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Clonează`)
};

export const schemas_cloning = /** @type {(inputs: Schemas_CloningInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se clonează...`)
};

export const schemas_saving = /** @type {(inputs: Schemas_SavingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se salvează...`)
};

export const json_recipes_title = /** @type {(inputs: Json_Recipes_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețete JSON`)
};

export const json_recipes_description = /** @type {(inputs: Json_Recipes_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șabloane administrate pentru publicarea schemelor de extragere.`)
};

export const json_recipes_new_recipe = /** @type {(inputs: Json_Recipes_New_RecipeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețetă nouă`)
};

export const json_recipes_no_recipes_found = /** @type {(inputs: Json_Recipes_No_Recipes_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit rețete`)
};

export const json_recipes_empty_body = /** @type {(inputs: Json_Recipes_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează o rețetă pentru a face disponibilă o schemă reutilizabilă.`)
};

export const json_recipes_loading = /** @type {(inputs: Json_Recipes_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă rețetele...`)
};

export const json_recipes_loading_recipe = /** @type {(inputs: Json_Recipes_Loading_RecipeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă rețeta...`)
};

export const json_recipes_counter_column = /** @type {(inputs: Json_Recipes_Counter_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Publicări`)
};

export const json_recipes_created_column = /** @type {(inputs: Json_Recipes_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creată`)
};

export const json_recipes_updated_column = /** @type {(inputs: Json_Recipes_Updated_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizată`)
};

export const json_recipes_json_fields_column = /** @type {(inputs: Json_Recipes_Json_Fields_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Câmpuri`)
};

export const json_recipes_sort_created_ascending = /** @type {(inputs: Json_Recipes_Sort_Created_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării ascendent`)
};

export const json_recipes_sort_created_descending = /** @type {(inputs: Json_Recipes_Sort_Created_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării descendent`)
};

export const json_recipes_showing_one = /** @type {(inputs: Json_Recipes_Showing_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rețetă pe această pagină.`)
};

export const json_recipes_showing_other = /** @type {(inputs: Json_Recipes_Showing_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rețete pe această pagină.`)
};

export const json_recipes_no_recipes_to_show = /** @type {(inputs: Json_Recipes_No_Recipes_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există rețete de afișat.`)
};

export const json_recipes_edit_aria = /** @type {(inputs: Json_Recipes_Edit_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Editează ${i?.name}`)
};

export const json_recipes_delete_aria = /** @type {(inputs: Json_Recipes_Delete_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Șterge ${i?.name}`)
};

export const json_recipes_new_title = /** @type {(inputs: Json_Recipes_New_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețetă JSON nouă`)
};

export const json_recipes_new_description = /** @type {(inputs: Json_Recipes_New_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Definește metadatele rețetei și structura JSON Schema.`)
};

export const json_recipes_edit_title = /** @type {(inputs: Json_Recipes_Edit_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează rețeta JSON`)
};

export const json_recipes_edit_description = /** @type {(inputs: Json_Recipes_Edit_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizează metadatele rețetei și structura JSON Schema.`)
};

export const json_recipes_save_recipe = /** @type {(inputs: Json_Recipes_Save_RecipeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează rețeta`)
};

export const json_recipes_save_changes = /** @type {(inputs: Json_Recipes_Save_ChangesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează modificările`)
};

export const json_recipes_created_success = /** @type {(inputs: Json_Recipes_Created_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Rețeta ${i?.name} a fost creată.`)
};

export const json_recipes_saved_success = /** @type {(inputs: Json_Recipes_Saved_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Rețeta ${i?.name} a fost salvată.`)
};

export const json_recipes_deleted_success = /** @type {(inputs: Json_Recipes_Deleted_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Rețeta ${i?.name} a fost ștearsă.`)
};

export const json_recipes_delete_confirm = /** @type {(inputs: Json_Recipes_Delete_ConfirmInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi această rețetă? Schemele publicate rămân neschimbate.`)
};

export const json_recipes_not_found_title = /** @type {(inputs: Json_Recipes_Not_Found_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețeta nu a fost găsită`)
};

export const json_recipes_not_found_body = /** @type {(inputs: Json_Recipes_Not_Found_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Această rețetă JSON nu există.`)
};

export const json_recipes_view_recipes = /** @type {(inputs: Json_Recipes_View_RecipesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Vezi rețetele`)
};

export const json_recipes_could_not_load = /** @type {(inputs: Json_Recipes_Could_Not_LoadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețeta nu a putut fi încărcată`)
};

export const json_recipes_editor_badge = /** @type {(inputs: Json_Recipes_Editor_BadgeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editor rețetă`)
};

export const json_recipes_general_settings = /** @type {(inputs: Json_Recipes_General_SettingsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setări rețetă`)
};

export const json_recipes_title_label = /** @type {(inputs: Json_Recipes_Title_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlu`)
};

export const json_recipes_title_placeholder = /** @type {(inputs: Json_Recipes_Title_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlul rețetei`)
};

export const json_recipes_description_label = /** @type {(inputs: Json_Recipes_Description_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descriere`)
};

export const json_recipes_description_placeholder = /** @type {(inputs: Json_Recipes_Description_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descrie când ar trebui folosită această rețetă...`)
};

export const json_recipes_structure_designer = /** @type {(inputs: Json_Recipes_Structure_DesignerInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Designer structură`)
};

export const json_recipes_visual_node_designer = /** @type {(inputs: Json_Recipes_Visual_Node_DesignerInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Designer vizual noduri`)
};

export const json_recipes_category_label = /** @type {(inputs: Json_Recipes_Category_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Categorie`)
};

export const json_recipes_others = /** @type {(inputs: Json_Recipes_OthersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Altele`)
};

export const json_recipes_manage_categories = /** @type {(inputs: Json_Recipes_Manage_CategoriesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Gestionează categoriile`)
};

export const json_recipes_validation_title_required = /** @type {(inputs: Json_Recipes_Validation_Title_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlul este obligatoriu.`)
};

export const json_recipes_validation_title_too_long = /** @type {(inputs: Json_Recipes_Validation_Title_Too_LongInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlul trebuie să aibă cel mult 160 de caractere.`)
};

export const json_recipes_validation_json_object = /** @type {(inputs: Json_Recipes_Validation_Json_ObjectInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`JSON-ul rețetei trebuie să fie un obiect JSON.`)
};

export const json_recipes_saving = /** @type {(inputs: Json_Recipes_SavingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se salvează...`)
};

export const json_recipes_deleting = /** @type {(inputs: Json_Recipes_DeletingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se șterge...`)
};

export const json_recipe_categories_title = /** @type {(inputs: Json_Recipe_Categories_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Categorii pentru rețete JSON`)
};

export const json_recipe_categories_description = /** @type {(inputs: Json_Recipe_Categories_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Definește etichete localizate pentru gruparea rețetelor JSON.`)
};

export const json_recipe_categories_title_en_label = /** @type {(inputs: Json_Recipe_Categories_Title_En_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlu în engleză`)
};

export const json_recipe_categories_title_ro_label = /** @type {(inputs: Json_Recipe_Categories_Title_Ro_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlu în română`)
};

export const json_recipe_categories_create_category = /** @type {(inputs: Json_Recipe_Categories_Create_CategoryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează categoria`)
};

export const json_recipe_categories_save_category = /** @type {(inputs: Json_Recipe_Categories_Save_CategoryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează categoria`)
};

export const json_recipe_categories_edit_title = /** @type {(inputs: Json_Recipe_Categories_Edit_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează categoria`)
};

export const json_recipe_categories_delete_confirm = /** @type {(inputs: Json_Recipe_Categories_Delete_ConfirmInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi această categorie? Rețetele trebuie mutate înainte de ștergere.`)
};

export const json_recipe_categories_loading = /** @type {(inputs: Json_Recipe_Categories_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă categoriile...`)
};

export const json_recipe_categories_could_not_load = /** @type {(inputs: Json_Recipe_Categories_Could_Not_LoadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Categoriile nu au putut fi încărcate`)
};

export const json_recipe_categories_empty_title = /** @type {(inputs: Json_Recipe_Categories_Empty_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există categorii încă`)
};

export const json_recipe_categories_empty_body = /** @type {(inputs: Json_Recipe_Categories_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețetele fără categorie vor apărea la Altele.`)
};

export const json_recipe_categories_created_success = /** @type {(inputs: Json_Recipe_Categories_Created_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Categoria ${i?.name} a fost creată.`)
};

export const json_recipe_categories_saved_success = /** @type {(inputs: Json_Recipe_Categories_Saved_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Categoria ${i?.name} a fost salvată.`)
};

export const json_recipe_categories_deleted_success = /** @type {(inputs: Json_Recipe_Categories_Deleted_SuccessInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Categoria ${i?.name} a fost ștearsă.`)
};

export const json_recipe_categories_validation_titles_required = /** @type {(inputs: Json_Recipe_Categories_Validation_Titles_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlurile în engleză și română sunt obligatorii.`)
};

export const json_recipe_categories_validation_titles_too_long = /** @type {(inputs: Json_Recipe_Categories_Validation_Titles_Too_LongInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Titlurile trebuie să aibă cel mult 160 de caractere.`)
};

export const json_recipe_categories_edit_aria = /** @type {(inputs: Json_Recipe_Categories_Edit_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Editează ${i?.name}`)
};

export const json_recipe_categories_delete_aria = /** @type {(inputs: Json_Recipe_Categories_Delete_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Șterge ${i?.name}`)
};

export const ocr_recipes_nav = /** @type {(inputs: Ocr_Recipes_NavInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețete OCR`)
};

export const ocr_recipes_title = /** @type {(inputs: Ocr_Recipes_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețete OCR`)
};

export const ocr_recipes_meta_description = /** @type {(inputs: Ocr_Recipes_Meta_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Răsfoiește rețete OCR JSON de sistem și clonează-le în schemele tale Syncra.`)
};

export const ocr_recipes_eyebrow = /** @type {(inputs: Ocr_Recipes_EyebrowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șabloane de extracție de sistem`)
};

export const ocr_recipes_hero_title = /** @type {(inputs: Ocr_Recipes_Hero_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pornește de la o rețetă OCR testată`)
};

export const ocr_recipes_hero_description = /** @type {(inputs: Ocr_Recipes_Hero_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Răsfoiește rețete JSON Schema reutilizabile pentru documente românești uzuale. Clonează o rețetă în spațiul tău, apoi ajusteaz-o în editorul de scheme înainte să rulezi joburi OCR.`)
};

export const ocr_recipes_search_label = /** @type {(inputs: Ocr_Recipes_Search_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută rețete`)
};

export const ocr_recipes_search_placeholder = /** @type {(inputs: Ocr_Recipes_Search_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută după rețetă, câmp, tip sau descriere`)
};

export const ocr_recipes_category_filter = /** @type {(inputs: Ocr_Recipes_Category_FilterInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Filtrează după categorie`)
};

export const ocr_recipes_all_categories = /** @type {(inputs: Ocr_Recipes_All_CategoriesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate categoriile`)
};

export const ocr_recipes_sort_label = /** @type {(inputs: Ocr_Recipes_Sort_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează rețetele`)
};

export const ocr_recipes_sort_popular = /** @type {(inputs: Ocr_Recipes_Sort_PopularInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cele mai clonate`)
};

export const ocr_recipes_sort_newest = /** @type {(inputs: Ocr_Recipes_Sort_NewestInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cele mai noi`)
};

export const ocr_recipes_sort_az = /** @type {(inputs: Ocr_Recipes_Sort_AzInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`A-Z`)
};

export const ocr_recipes_showing_one = /** @type {(inputs: Ocr_Recipes_Showing_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rețetă`)
};

export const ocr_recipes_showing_other = /** @type {(inputs: Ocr_Recipes_Showing_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rețete`)
};

export const ocr_recipes_no_matches_title = /** @type {(inputs: Ocr_Recipes_No_Matches_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio rețetă nu se potrivește căutării`)
};

export const ocr_recipes_no_matches_body = /** @type {(inputs: Ocr_Recipes_No_Matches_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Golește câmpul de căutare pentru a vedea toate rețetele de sistem.`)
};

export const ocr_recipes_others = /** @type {(inputs: Ocr_Recipes_OthersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Altele`)
};

export const ocr_recipes_fields_one = /** @type {(inputs: Ocr_Recipes_Fields_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmp`)
};

export const ocr_recipes_fields_other = /** @type {(inputs: Ocr_Recipes_Fields_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmpuri`)
};

export const ocr_recipes_required_one = /** @type {(inputs: Ocr_Recipes_Required_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} obligatoriu`)
};

export const ocr_recipes_required_other = /** @type {(inputs: Ocr_Recipes_Required_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} obligatorii`)
};

export const ocr_recipes_deploys_one = /** @type {(inputs: Ocr_Recipes_Deploys_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} clonare`)
};

export const ocr_recipes_deploys_other = /** @type {(inputs: Ocr_Recipes_Deploys_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} clonări`)
};

export const ocr_recipes_json_fields = /** @type {(inputs: Ocr_Recipes_Json_FieldsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Câmpuri JSON`)
};

export const ocr_recipes_system_recipe = /** @type {(inputs: Ocr_Recipes_System_RecipeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețetă de sistem`)
};

export const ocr_recipes_strict_schema = /** @type {(inputs: Ocr_Recipes_Strict_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`JSON Schema strictă`)
};

export const ocr_recipes_required = /** @type {(inputs: Ocr_Recipes_RequiredInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Obligatoriu`)
};

export const ocr_recipes_preview_json = /** @type {(inputs: Ocr_Recipes_Preview_JsonInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizează JSON`)
};

export const ocr_recipes_no_fields = /** @type {(inputs: Ocr_Recipes_No_FieldsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu sunt definite câmpuri JSON.`)
};

export const ocr_recipes_clone_recipe = /** @type {(inputs: Ocr_Recipes_Clone_RecipeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Clonează rețeta`)
};

export const ocr_recipes_clone_aria = /** @type {(inputs: Ocr_Recipes_Clone_AriaInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Clonează ${i?.name}`)
};

export const ocr_recipes_log_in_to_clone = /** @type {(inputs: Ocr_Recipes_Log_In_To_CloneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Autentifică-te ca să clonezi`)
};

export const ocr_recipes_clone_failed = /** @type {(inputs: Ocr_Recipes_Clone_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețeta nu a putut fi clonată.`)
};

export const ocr_recipes_load_failed = /** @type {(inputs: Ocr_Recipes_Load_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rețetele OCR nu au putut fi încărcate.`)
};

export const jobs_page_title = /** @type {(inputs: Jobs_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Joburi`)
};

export const jobs_missing_schema_id = /** @type {(inputs: Jobs_Missing_Schema_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul schemei lipsește`)
};

export const jobs_missing_job_id = /** @type {(inputs: Jobs_Missing_Job_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul jobului lipsește`)
};

export const jobs_delete_bulk_title_one = /** @type {(inputs: Jobs_Delete_Bulk_Title_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} job?`)
};

export const jobs_delete_bulk_title_other = /** @type {(inputs: Jobs_Delete_Bulk_Title_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} joburi?`)
};

export const jobs_delete_bulk_description_one = /** @type {(inputs: Jobs_Delete_Bulk_Description_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Elimini ${i?.count} job selectat din listă. Documentele generate rămân disponibile.`)
};

export const jobs_delete_bulk_description_other = /** @type {(inputs: Jobs_Delete_Bulk_Description_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Elimini ${i?.count} joburi selectate din listă. Documentele generate rămân disponibile.`)
};

export const jobs_delete_single_title = /** @type {(inputs: Jobs_Delete_Single_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi jobul?`)
};

export const jobs_delete_single_description = /** @type {(inputs: Jobs_Delete_Single_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Elimini „${i?.name}” din lista de joburi. Documentele generate rămân disponibile.`)
};

export const jobs_status_queued = /** @type {(inputs: Jobs_Status_QueuedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În coadă`)
};

export const jobs_status_pending = /** @type {(inputs: Jobs_Status_PendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În așteptare`)
};

export const jobs_status_processing = /** @type {(inputs: Jobs_Status_ProcessingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În procesare`)
};

export const jobs_status_completed = /** @type {(inputs: Jobs_Status_CompletedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Finalizat`)
};

export const jobs_status_failed = /** @type {(inputs: Jobs_Status_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Eșuat`)
};

export const jobs_inline_schema = /** @type {(inputs: Jobs_Inline_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă inline`)
};

export const jobs_no_schema = /** @type {(inputs: Jobs_No_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio schemă`)
};

export const jobs_schema = /** @type {(inputs: Jobs_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă`)
};

export const jobs_select_all_on_page = /** @type {(inputs: Jobs_Select_All_On_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează toate joburile de pe această pagină`)
};

export const jobs_select_job = /** @type {(inputs: Jobs_Select_JobInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Selectează ${i?.name}`)
};

export const jobs_filename_column = /** @type {(inputs: Jobs_Filename_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume fișier`)
};

export const jobs_status_column = /** @type {(inputs: Jobs_Status_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stare`)
};

export const jobs_created_column = /** @type {(inputs: Jobs_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creat`)
};

export const jobs_file_size_column = /** @type {(inputs: Jobs_File_Size_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dimensiune fișier`)
};

export const jobs_pages_column = /** @type {(inputs: Jobs_Pages_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagini`)
};

export const jobs_new_job = /** @type {(inputs: Jobs_New_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Job nou`)
};

export const jobs_no_jobs_found = /** @type {(inputs: Jobs_No_Jobs_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit joburi`)
};

export const jobs_empty_body = /** @type {(inputs: Jobs_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pornește un job batch pentru a procesa documente și a urmări progresul aici.`)
};

export const jobs_showing_jobs_one = /** @type {(inputs: Jobs_Showing_Jobs_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} job pe această pagină.`)
};

export const jobs_showing_jobs_other = /** @type {(inputs: Jobs_Showing_Jobs_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} joburi pe această pagină.`)
};

export const jobs_no_jobs_to_show = /** @type {(inputs: Jobs_No_Jobs_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există joburi de afișat.`)
};

export const jobs_selected_count_one = /** @type {(inputs: Jobs_Selected_Count_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectat`)
};

export const jobs_selected_count_other = /** @type {(inputs: Jobs_Selected_Count_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectate`)
};

export const jobs_deleting = /** @type {(inputs: Jobs_DeletingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se șterge...`)
};

export const jobs_delete_job = /** @type {(inputs: Jobs_Delete_JobInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Șterge ${i?.name}`)
};

export const jobs_saved_extraction_schema = /** @type {(inputs: Jobs_Saved_Extraction_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă de extragere salvată`)
};

export const jobs_inline_schema_description = /** @type {(inputs: Jobs_Inline_Schema_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema a fost trimisă direct cu acest job.`)
};

export const jobs_extraction_schema_details = /** @type {(inputs: Jobs_Extraction_Schema_DetailsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Detalii schemă de extragere.`)
};

export const new_job_missing_document_id = /** @type {(inputs: New_Job_Missing_Document_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul documentului lipsește`)
};

export const new_job_failed_create = /** @type {(inputs: New_Job_Failed_CreateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Crearea jobului OCR a eșuat`)
};

export const new_job_insufficient_credits_buy = /** @type {(inputs: New_Job_Insufficient_Credits_BuyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite insuficiente. Cumpără credite pentru a procesa acest document.`)
};

export const new_job_failed_load_document = /** @type {(inputs: New_Job_Failed_Load_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încărcarea documentului a eșuat`)
};

export const new_job_invalid_document_response = /** @type {(inputs: New_Job_Invalid_Document_ResponseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Răspuns document invalid`)
};

export const new_job_failed_load_schemas = /** @type {(inputs: New_Job_Failed_Load_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încărcarea schemelor a eșuat`)
};

export const new_job_invalid_schema_response = /** @type {(inputs: New_Job_Invalid_Schema_ResponseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Răspuns schemă invalid`)
};

export const new_job_invalid_job_response = /** @type {(inputs: New_Job_Invalid_Job_ResponseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Răspuns job OCR invalid`)
};

export const new_job_failed_load_job = /** @type {(inputs: New_Job_Failed_Load_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încărcarea jobului OCR a eșuat`)
};

export const new_job_failed_poll_job = /** @type {(inputs: New_Job_Failed_Poll_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Interogarea jobului OCR a eșuat`)
};

export const new_job_select_schema = /** @type {(inputs: New_Job_Select_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează schema`)
};

export const new_job_select_schema_placeholder = /** @type {(inputs: New_Job_Select_Schema_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează schema`)
};

export const new_job_configure_payload_format = /** @type {(inputs: New_Job_Configure_Payload_FormatInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Configurează formatul payloadului`)
};

export const new_job_upload_documents = /** @type {(inputs: New_Job_Upload_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încarcă documente`)
};

export const new_job_files_selected_one = /** @type {(inputs: New_Job_Files_Selected_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} fișier selectat`)
};

export const new_job_files_selected_other = /** @type {(inputs: New_Job_Files_Selected_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} fișiere selectate`)
};

export const new_job_drag_or_browse_files = /** @type {(inputs: New_Job_Drag_Or_Browse_FilesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Trage sau alege fișiere`)
};

export const new_job_run_monitor = /** @type {(inputs: New_Job_Run_MonitorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rulează și monitorizează`)
};

export const new_job_processing_batch = /** @type {(inputs: New_Job_Processing_BatchInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se procesează batchul...`)
};

export const new_job_start_extraction_pipeline = /** @type {(inputs: New_Job_Start_Extraction_PipelineInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pornește fluxul de extragere`)
};

export const new_job_select_extraction_schema = /** @type {(inputs: New_Job_Select_Extraction_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează schema de extragere`)
};

export const new_job_select_schema_description = /** @type {(inputs: New_Job_Select_Schema_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Alege o schemă pentru câmpuri structurate de extragere AI sau continuă în modul OCR brut.`)
};

export const new_job_select_extraction_schema_aria = /** @type {(inputs: New_Job_Select_Extraction_Schema_AriaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează schema de extragere`)
};

export const new_job_search_schemas = /** @type {(inputs: New_Job_Search_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută scheme...`)
};

export const new_job_loading_schemas = /** @type {(inputs: New_Job_Loading_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă schemele`)
};

export const new_job_no_schemas_found = /** @type {(inputs: New_Job_No_Schemas_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit scheme.`)
};

export const new_job_no_schema_ocr_only = /** @type {(inputs: New_Job_No_Schema_Ocr_OnlyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fără schemă (doar OCR)`)
};

export const new_job_no_schema_description = /** @type {(inputs: New_Job_No_Schema_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesează documente fără extragere structurată.`)
};

export const new_job_no_personal_schemas = /** @type {(inputs: New_Job_No_Personal_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există scheme personale disponibile.`)
};

export const new_job_create_one = /** @type {(inputs: New_Job_Create_OneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează una`)
};

export const new_job_selected_schema_help = /** @type {(inputs: New_Job_Selected_Schema_HelpInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema selectată definește câmpurile structurate care vor fi extrase din fișiere.`)
};

export const new_job_no_schema_selected_help = /** @type {(inputs: New_Job_No_Schema_Selected_HelpInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio schemă selectată. Fișierele vor fi procesate OCR fără extragere structurată.`)
};

export const new_job_target_mapped_fields = /** @type {(inputs: New_Job_Target_Mapped_FieldsInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Câmpuri mapate țintă (${i?.count})`)
};

export const new_job_no_fields_defined = /** @type {(inputs: New_Job_No_Fields_DefinedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există câmpuri definite în această schemă.`)
};

export const new_job_ocr_only_mode_active = /** @type {(inputs: New_Job_Ocr_Only_Mode_ActiveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mod doar OCR activ`)
};

export const new_job_ocr_only_mode_body = /** @type {(inputs: New_Job_Ocr_Only_Mode_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documentele vor fi procesate pentru text OCR de înaltă fidelitate fără conversia câmpurilor în payloaduri structurate. Selectează o schemă mai sus pentru extragere automată de câmpuri.`)
};

export const new_job_upload_documents_description = /** @type {(inputs: New_Job_Upload_Documents_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Selectează fișiere PDF sau imagini pentru extragerea conținutului. Poți încărca simultan până la ${i?.count} fișiere.`)
};

export const new_job_dropzone_title = /** @type {(inputs: New_Job_Dropzone_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Trage fișiere aici sau fă clic pentru încărcare`)
};

export const new_job_dropzone_description = /** @type {(inputs: New_Job_Dropzone_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Acceptă PDF, PNG și JPG până la ${i?.size} per fișier`)
};

export const new_job_browse_files = /** @type {(inputs: New_Job_Browse_FilesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Alege fișiere`)
};

export const new_job_pending_upload_queue = /** @type {(inputs: New_Job_Pending_Upload_QueueInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Coadă încărcări în așteptare (${i?.count})`)
};

export const new_job_clear_all = /** @type {(inputs: New_Job_Clear_AllInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge tot`)
};

export const new_job_remove_file = /** @type {(inputs: New_Job_Remove_FileInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Elimină fișierul`)
};

export const new_job_extraction_queue_results = /** @type {(inputs: New_Job_Extraction_Queue_ResultsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Coadă și rezultate extragere`)
};

export const new_job_file_count_one = /** @type {(inputs: New_Job_File_Count_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} fișier`)
};

export const new_job_file_count_other = /** @type {(inputs: New_Job_File_Count_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} fișiere`)
};

export const new_job_total = /** @type {(inputs: New_Job_TotalInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.label} în total`)
};

export const new_job_active_batch_status = /** @type {(inputs: New_Job_Active_Batch_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stare batch activ`)
};

export const new_job_active_batch_description = /** @type {(inputs: New_Job_Active_Batch_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Monitorizează documentele batchului în timp real.`)
};

export const new_job_progress = /** @type {(inputs: New_Job_ProgressInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Progres: ${i?.progress}%`)
};

export const new_job_total_files = /** @type {(inputs: New_Job_Total_FilesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Total fișiere`)
};

export const new_job_completed = /** @type {(inputs: New_Job_CompletedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Finalizate`)
};

export const new_job_processing = /** @type {(inputs: New_Job_ProcessingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În procesare`)
};

export const new_job_failed = /** @type {(inputs: New_Job_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Eșuate`)
};

export const new_job_no_active_extraction_jobs = /** @type {(inputs: New_Job_No_Active_Extraction_JobsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există joburi active de extragere`)
};

export const new_job_no_active_extraction_jobs_body = /** @type {(inputs: New_Job_No_Active_Extraction_Jobs_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încarcă documente mai sus și selectează o schemă pentru a porni procesul automat de OCR și extragere.`)
};

export const new_job_preview_document = /** @type {(inputs: New_Job_Preview_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizează documentul`)
};

export const new_job_preview_unavailable = /** @type {(inputs: New_Job_Preview_UnavailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizarea documentului nu este disponibilă încă`)
};

export const new_job_remove_failed_job = /** @type {(inputs: New_Job_Remove_Failed_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Elimină jobul eșuat`)
};

export const new_job_queueing_documents = /** @type {(inputs: New_Job_Queueing_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se pun documentele în coadă...`)
};

export const new_job_extracting_content = /** @type {(inputs: New_Job_Extracting_ContentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se extrage conținutul...`)
};

export const new_job_run_extraction_one = /** @type {(inputs: New_Job_Run_Extraction_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Rulează extragerea (${i?.count} fișier)`)
};

export const new_job_run_extraction_other = /** @type {(inputs: New_Job_Run_Extraction_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Rulează extragerea (${i?.count} fișiere)`)
};

export const new_job_insufficient_credits_document = /** @type {(inputs: New_Job_Insufficient_Credits_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite insuficiente pentru acest document.`)
};

export const new_job_processing_failed = /** @type {(inputs: New_Job_Processing_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesarea a eșuat`)
};

export const new_job_processed = /** @type {(inputs: New_Job_ProcessedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesat`)
};

export const new_job_document_id = /** @type {(inputs: New_Job_Document_IdInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Document ${i?.id}`)
};

export const new_job_creating_job = /** @type {(inputs: New_Job_Creating_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se creează jobul OCR...`)
};

export const new_job_queued_processing = /** @type {(inputs: New_Job_Queued_ProcessingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În coadă pentru procesare...`)
};

export const new_job_extracting_entities = /** @type {(inputs: New_Job_Extracting_EntitiesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se extrag entitățile...`)
};

export const common_apply = /** @type {(inputs: Common_ApplyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Aplică`)
};

export const common_clear = /** @type {(inputs: Common_ClearInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const common_saving = /** @type {(inputs: Common_SavingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se salvează...`)
};

export const common_loading = /** @type {(inputs: Common_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă...`)
};

export const common_refresh = /** @type {(inputs: Common_RefreshInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîmprospătează`)
};

export const common_connected = /** @type {(inputs: Common_ConnectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Conectat`)
};

export const common_connect = /** @type {(inputs: Common_ConnectInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Conectează`)
};

export const common_download = /** @type {(inputs: Common_DownloadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descarcă`)
};

export const common_today = /** @type {(inputs: Common_TodayInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Astăzi`)
};

export const common_this_week = /** @type {(inputs: Common_This_WeekInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Săptămâna aceasta`)
};

export const common_this_month = /** @type {(inputs: Common_This_MonthInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Luna aceasta`)
};

export const common_any = /** @type {(inputs: Common_AnyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Oricare`)
};

export const billing_unavailable = /** @type {(inputs: Billing_UnavailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Indisponibil`)
};

export const billing_credit_blocks_error = /** @type {(inputs: Billing_Credit_Blocks_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creditele trebuie cumpărate în blocuri de 1000 de credite.`)
};

export const billing_checkout_unavailable = /** @type {(inputs: Billing_Checkout_UnavailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Checkoutul nu poate fi pornit`)
};

export const billing_payment_received_title = /** @type {(inputs: Billing_Payment_Received_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Plată primită`)
};

export const billing_payment_received_body = /** @type {(inputs: Billing_Payment_Received_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Soldul de credite se va actualiza în scurt timp după confirmarea plății.`)
};

export const billing_checkout_canceled_title = /** @type {(inputs: Billing_Checkout_Canceled_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Checkout anulat`)
};

export const billing_checkout_canceled_body = /** @type {(inputs: Billing_Checkout_Canceled_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu au fost cumpărate credite și nu s-au făcut debitări.`)
};

export const billing_available_balance = /** @type {(inputs: Billing_Available_BalanceInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sold disponibil`)
};

export const billing_conversion = /** @type {(inputs: Billing_ConversionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Conversie`)
};

export const billing_conversion_rate = /** @type {(inputs: Billing_Conversion_RateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`1 credit = 1 pagină`)
};

export const billing_balance_checked_upload = /** @type {(inputs: Billing_Balance_Checked_UploadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sold verificat la încărcare`)
};

export const billing_debited_after_success = /** @type {(inputs: Billing_Debited_After_SuccessInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Debitat după succes`)
};

export const billing_secure_stripe_checkout = /** @type {(inputs: Billing_Secure_Stripe_CheckoutInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Checkout Stripe securizat`)
};

export const billing_purchase_credits = /** @type {(inputs: Billing_Purchase_CreditsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cumpără credite`)
};

export const billing_credits_to_purchase = /** @type {(inputs: Billing_Credits_To_PurchaseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite de cumpărat`)
};

export const billing_volume_discount_tiers = /** @type {(inputs: Billing_Volume_Discount_TiersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niveluri de discount pentru volum`)
};

export const billing_total_to_pay = /** @type {(inputs: Billing_Total_To_PayInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Total de plată`)
};

export const billing_base_price = /** @type {(inputs: Billing_Base_PriceInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Preț de bază`)
};

export const billing_volume_discount = /** @type {(inputs: Billing_Volume_DiscountInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Discount de volum`)
};

export const billing_starting_checkout = /** @type {(inputs: Billing_Starting_CheckoutInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se pornește checkoutul...`)
};

export const billing_secure_checkout = /** @type {(inputs: Billing_Secure_CheckoutInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Checkout securizat`)
};

export const billing_buy_credits = /** @type {(inputs: Billing_Buy_CreditsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cumpără credite`)
};

export const billing_orders_page_title = /** @type {(inputs: Billing_Orders_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Comenzi de facturare`)
};

export const billing_orders_order_date_filter = /** @type {(inputs: Billing_Orders_Order_Date_FilterInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Data comenzii`)
};

export const billing_orders_amount_column = /** @type {(inputs: Billing_Orders_Amount_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sumă`)
};

export const billing_orders_credits_column = /** @type {(inputs: Billing_Orders_Credits_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite`)
};

export const billing_orders_status_column = /** @type {(inputs: Billing_Orders_Status_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Stare`)
};

export const billing_orders_payment_datetime_column = /** @type {(inputs: Billing_Orders_Payment_Datetime_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Data și ora plății`)
};

export const billing_orders_invoice_column = /** @type {(inputs: Billing_Orders_Invoice_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Factură`)
};

export const billing_orders_presets = /** @type {(inputs: Billing_Orders_PresetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Presetări`)
};

export const billing_orders_filter_status = /** @type {(inputs: Billing_Orders_Filter_StatusInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Filtrează după stare`)
};

export const billing_orders_all_orders = /** @type {(inputs: Billing_Orders_All_OrdersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate comenzile`)
};

export const billing_orders_clear_filters = /** @type {(inputs: Billing_Orders_Clear_FiltersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge filtrele`)
};

export const billing_orders_clear_filters_action = /** @type {(inputs: Billing_Orders_Clear_Filters_ActionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge filtrele`)
};

export const billing_orders_no_orders_found = /** @type {(inputs: Billing_Orders_No_Orders_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit comenzi de facturare`)
};

export const billing_orders_no_orders_yet = /** @type {(inputs: Billing_Orders_No_Orders_YetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există încă comenzi de facturare`)
};

export const billing_orders_no_orders_match = /** @type {(inputs: Billing_Orders_No_Orders_MatchInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio comandă de facturare nu corespunde filtrelor selectate.`)
};

export const billing_orders_empty_body = /** @type {(inputs: Billing_Orders_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Comenzile de cumpărare credite vor apărea aici după pornirea checkoutului.`)
};

export const billing_orders_showing_one = /** @type {(inputs: Billing_Orders_Showing_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} comandă pe această pagină.`)
};

export const billing_orders_showing_other = /** @type {(inputs: Billing_Orders_Showing_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} comenzi pe această pagină.`)
};

export const billing_orders_none_to_show = /** @type {(inputs: Billing_Orders_None_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există comenzi de facturare de afișat.`)
};

export const billing_orders_sort_order_date_ascending = /** @type {(inputs: Billing_Orders_Sort_Order_Date_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data comenzii ascendent`)
};

export const billing_orders_sort_order_date_descending = /** @type {(inputs: Billing_Orders_Sort_Order_Date_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data comenzii descendent`)
};

export const billing_order_status_pending = /** @type {(inputs: Billing_Order_Status_PendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În așteptare`)
};

export const billing_order_status_paid = /** @type {(inputs: Billing_Order_Status_PaidInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Plătită`)
};

export const billing_order_status_failed = /** @type {(inputs: Billing_Order_Status_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Eșuată`)
};

export const billing_order_status_refunded = /** @type {(inputs: Billing_Order_Status_RefundedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rambursată`)
};

export const billing_order_status_canceled = /** @type {(inputs: Billing_Order_Status_CanceledInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anulată`)
};

export const billing_orders_invoice_pdf_title = /** @type {(inputs: Billing_Orders_Invoice_Pdf_TitleInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Previzualizează PDF-ul ${i?.invoice}`)
};

export const billing_orders_invoice_preview_title = /** @type {(inputs: Billing_Orders_Invoice_Preview_TitleInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Factura ${i?.invoice}`)
};

export const billing_orders_invoice_preview_description = /** @type {(inputs: Billing_Orders_Invoice_Preview_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizare PDF`)
};

export const billing_orders_invoice_iframe_title = /** @type {(inputs: Billing_Orders_Invoice_Iframe_TitleInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Previzualizare PDF factură ${i?.invoice}`)
};

export const billing_orders_download_invoice = /** @type {(inputs: Billing_Orders_Download_InvoiceInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descarcă`)
};

export const credit_usage_page_title = /** @type {(inputs: Credit_Usage_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Istoric utilizare credite`)
};

export const credit_usage_date_range_filter = /** @type {(inputs: Credit_Usage_Date_Range_FilterInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Interval de date`)
};

export const credit_usage_created_column = /** @type {(inputs: Credit_Usage_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creat`)
};

export const credit_usage_type_column = /** @type {(inputs: Credit_Usage_Type_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Tip`)
};

export const credit_usage_credits_column = /** @type {(inputs: Credit_Usage_Credits_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Credite`)
};

export const credit_usage_related_id_column = /** @type {(inputs: Credit_Usage_Related_Id_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID asociat`)
};

export const credit_usage_filter_type = /** @type {(inputs: Credit_Usage_Filter_TypeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Filtrează după tip`)
};

export const credit_usage_all_activity = /** @type {(inputs: Credit_Usage_All_ActivityInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toată activitatea`)
};

export const credit_usage_type_purchase = /** @type {(inputs: Credit_Usage_Type_PurchaseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cumpărare`)
};

export const credit_usage_type_debit = /** @type {(inputs: Credit_Usage_Type_DebitInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Debitare`)
};

export const credit_usage_no_usage_found = /** @type {(inputs: Credit_Usage_No_Usage_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-a găsit utilizare de credite`)
};

export const credit_usage_no_usage_yet = /** @type {(inputs: Credit_Usage_No_Usage_YetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există încă utilizare de credite`)
};

export const credit_usage_no_usage_match = /** @type {(inputs: Credit_Usage_No_Usage_MatchInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niciun istoric de utilizare credite nu corespunde filtrelor selectate.`)
};

export const credit_usage_empty_body = /** @type {(inputs: Credit_Usage_Empty_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cumpărările și debitările vor apărea aici după ce se finalizează.`)
};

export const credit_usage_showing_one = /** @type {(inputs: Credit_Usage_Showing_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} înregistrare pe această pagină.`)
};

export const credit_usage_showing_other = /** @type {(inputs: Credit_Usage_Showing_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} înregistrări pe această pagină.`)
};

export const credit_usage_none_to_show = /** @type {(inputs: Credit_Usage_None_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există istoric de utilizare credite de afișat.`)
};

export const credit_usage_sort_created_ascending = /** @type {(inputs: Credit_Usage_Sort_Created_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării ascendent`)
};

export const credit_usage_sort_created_descending = /** @type {(inputs: Credit_Usage_Sort_Created_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării descendent`)
};

export const account_settings_title = /** @type {(inputs: Account_Settings_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setări cont`)
};

export const account_settings_description = /** @type {(inputs: Account_Settings_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Gestionează detaliile contului și securitatea.`)
};

export const account_settings_nav_label = /** @type {(inputs: Account_Settings_Nav_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setări cont`)
};

export const account_settings_account_fallback = /** @type {(inputs: Account_Settings_Account_FallbackInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cont`)
};

export const account_settings_no_email_address = /** @type {(inputs: Account_Settings_No_Email_AddressInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fără adresă de email`)
};

export const account_settings_general = /** @type {(inputs: Account_Settings_GeneralInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`General`)
};

export const account_settings_security = /** @type {(inputs: Account_Settings_SecurityInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Securitate`)
};

export const account_settings_sessions = /** @type {(inputs: Account_Settings_SessionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sesiuni`)
};

export const account_settings_linked_accounts = /** @type {(inputs: Account_Settings_Linked_AccountsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Conturi conectate`)
};

export const account_settings_update_error = /** @type {(inputs: Account_Settings_Update_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setările contului nu pot fi actualizate.`)
};

export const account_settings_save_error = /** @type {(inputs: Account_Settings_Save_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Modificările nu pot fi salvate.`)
};

export const account_settings_revoke_session_title = /** @type {(inputs: Account_Settings_Revoke_Session_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Revoci sesiunea?`)
};

export const account_settings_revoke_session_description = /** @type {(inputs: Account_Settings_Revoke_Session_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Revoci ${i?.session}. Dispozitivul va trebui să se autentifice din nou.`)
};

export const account_settings_revoke = /** @type {(inputs: Account_Settings_RevokeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Revocă`)
};

export const account_settings_session_revoked = /** @type {(inputs: Account_Settings_Session_RevokedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sesiune revocată.`)
};

export const account_settings_unlink_provider_title = /** @type {(inputs: Account_Settings_Unlink_Provider_TitleInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Deconectezi ${i?.provider}?`)
};

export const account_settings_unlink_provider_description = /** @type {(inputs: Account_Settings_Unlink_Provider_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Elimină autentificarea cu ${i?.provider} din acest cont. O poți reconecta mai târziu.`)
};

export const account_settings_unlink = /** @type {(inputs: Account_Settings_UnlinkInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deconectează`)
};

export const account_settings_linked_account_removed = /** @type {(inputs: Account_Settings_Linked_Account_RemovedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Contul conectat a fost eliminat.`)
};

export const account_settings_avatar_saved = /** @type {(inputs: Account_Settings_Avatar_SavedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Avatar actualizat.`)
};

export const account_settings_name_saved = /** @type {(inputs: Account_Settings_Name_SavedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Numele afișat a fost salvat.`)
};

export const account_settings_email_saved = /** @type {(inputs: Account_Settings_Email_SavedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adresa de email a fost salvată.`)
};

export const account_settings_language_saved = /** @type {(inputs: Account_Settings_Language_SavedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Preferința de limbă a fost salvată.`)
};

export const account_settings_password_updated = /** @type {(inputs: Account_Settings_Password_UpdatedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Parola a fost actualizată.`)
};

export const account_settings_current_session = /** @type {(inputs: Account_Settings_Current_SessionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`sesiunea curentă`)
};

export const account_settings_browser_session = /** @type {(inputs: Account_Settings_Browser_SessionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`sesiune browser`)
};

export const account_settings_session_created_at = /** @type {(inputs: Account_Settings_Session_Created_AtInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Creată ${i?.date}`)
};

export const account_settings_session_ip_created_at = /** @type {(inputs: Account_Settings_Session_Ip_Created_AtInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.ip} - creată ${i?.date}`)
};

export const account_settings_unknown = /** @type {(inputs: Account_Settings_UnknownInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Necunoscut`)
};

export const account_settings_avatar = /** @type {(inputs: Account_Settings_AvatarInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Avatar`)
};

export const account_settings_avatar_description = /** @type {(inputs: Account_Settings_Avatar_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încarcă o imagine de profil afișată în contul tău.`)
};

export const account_settings_avatar_uploading = /** @type {(inputs: Account_Settings_Avatar_UploadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă...`)
};

export const account_settings_avatar_upload = /** @type {(inputs: Account_Settings_Avatar_UploadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Fă clic pentru încărcare și decupare`)
};

export const account_settings_avatar_file_hint = /** @type {(inputs: Account_Settings_Avatar_File_HintInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`PNG, JPG, GIF, AVIF, APNG, SVG, WEBP până la 5 MB.`)
};

export const account_settings_crop_avatar = /** @type {(inputs: Account_Settings_Crop_AvatarInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Decupează avatarul`)
};

export const account_settings_crop_avatar_description = /** @type {(inputs: Account_Settings_Crop_Avatar_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ajustează zona de decupare a avatarului înainte de salvare.`)
};

export const account_settings_display_name = /** @type {(inputs: Account_Settings_Display_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume afișat`)
};

export const account_settings_email_address = /** @type {(inputs: Account_Settings_Email_AddressInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adresă de email`)
};

export const account_settings_language = /** @type {(inputs: Account_Settings_LanguageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Limbă`)
};

export const account_settings_save_name = /** @type {(inputs: Account_Settings_Save_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează numele`)
};

export const account_settings_save_email = /** @type {(inputs: Account_Settings_Save_EmailInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează emailul`)
};

export const account_settings_save_language = /** @type {(inputs: Account_Settings_Save_LanguageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează limba`)
};

export const account_settings_save_password = /** @type {(inputs: Account_Settings_Save_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează parola`)
};

export const account_settings_new_password = /** @type {(inputs: Account_Settings_New_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Parolă nouă`)
};

export const account_settings_confirm_password = /** @type {(inputs: Account_Settings_Confirm_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Confirmă parola`)
};

export const account_settings_security_description = /** @type {(inputs: Account_Settings_Security_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schimbă parola contului.`)
};

export const account_settings_sessions_description = /** @type {(inputs: Account_Settings_Sessions_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Verifică browserele și dispozitivele autentificate în acest cont.`)
};

export const account_settings_loading_sessions = /** @type {(inputs: Account_Settings_Loading_SessionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă sesiunile...`)
};

export const account_settings_no_sessions = /** @type {(inputs: Account_Settings_No_SessionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit sesiuni active.`)
};

export const account_settings_current = /** @type {(inputs: Account_Settings_CurrentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Curentă`)
};

export const account_settings_expires = /** @type {(inputs: Account_Settings_ExpiresInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Expiră ${i?.date}`)
};

export const account_settings_current_session_cannot_revoke = /** @type {(inputs: Account_Settings_Current_Session_Cannot_RevokeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sesiunea curentă nu poate fi revocată`)
};

export const account_settings_revoke_session = /** @type {(inputs: Account_Settings_Revoke_SessionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Revocă sesiunea`)
};

export const account_settings_revoking = /** @type {(inputs: Account_Settings_RevokingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se revocă...`)
};

export const account_settings_linked_accounts_description = /** @type {(inputs: Account_Settings_Linked_Accounts_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Gestionează metodele de autentificare conectate la acest cont.`)
};

export const account_settings_loading_linked_accounts = /** @type {(inputs: Account_Settings_Loading_Linked_AccountsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă conturile conectate...`)
};

export const account_settings_no_sign_in_methods = /** @type {(inputs: Account_Settings_No_Sign_In_MethodsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu au fost returnate metode de autentificare.`)
};

export const account_settings_email_password = /** @type {(inputs: Account_Settings_Email_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Email/parolă`)
};

export const account_settings_password_enabled = /** @type {(inputs: Account_Settings_Password_EnabledInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Autentificarea cu parolă este activă pentru ${i?.email}.`)
};

export const account_settings_add_password = /** @type {(inputs: Account_Settings_Add_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adaugă o parolă pentru autentificarea cu email.`)
};

export const account_settings_set_password = /** @type {(inputs: Account_Settings_Set_PasswordInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setează parola`)
};

export const account_settings_provider_google_description = /** @type {(inputs: Account_Settings_Provider_Google_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Folosește contul Google pentru autentificare.`)
};

export const account_settings_provider_github_description = /** @type {(inputs: Account_Settings_Provider_Github_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Folosește contul GitHub pentru autentificare.`)
};

export const account_settings_linked_at = /** @type {(inputs: Account_Settings_Linked_AtInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Conectat la ${i?.date}`)
};

export const account_settings_unlinking = /** @type {(inputs: Account_Settings_UnlinkingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se deconectează...`)
};

export const account_settings_unavailable_title = /** @type {(inputs: Account_Settings_Unavailable_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Indisponibil`)
};

export const account_settings_unavailable_body = /** @type {(inputs: Account_Settings_Unavailable_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Această secțiune nu este disponibilă încă.`)
};

export const billing_profile_title = /** @type {(inputs: Billing_Profile_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Informații de facturare`)
};

export const billing_profile_description = /** @type {(inputs: Billing_Profile_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Gestionează detaliile de facturare folosite pentru facturi.`)
};

export const billing_profile_load_error = /** @type {(inputs: Billing_Profile_Load_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Informațiile de facturare nu pot fi încărcate.`)
};

export const billing_profile_save_error = /** @type {(inputs: Billing_Profile_Save_ErrorInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Informațiile de facturare nu pot fi salvate.`)
};

export const billing_profile_saved = /** @type {(inputs: Billing_Profile_SavedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Informațiile de facturare au fost salvate.`)
};

export const billing_profile_company_name = /** @type {(inputs: Billing_Profile_Company_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume companie`)
};

export const billing_profile_full_name = /** @type {(inputs: Billing_Profile_Full_NameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume complet`)
};

export const billing_profile_error_title = /** @type {(inputs: Billing_Profile_Error_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`A apărut o eroare`)
};

export const billing_profile_loading = /** @type {(inputs: Billing_Profile_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă informațiile de facturare...`)
};

export const billing_profile_loading_body = /** @type {(inputs: Billing_Profile_Loading_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Așteaptă cât timp preluăm detaliile profilului.`)
};

export const billing_profile_failed_load = /** @type {(inputs: Billing_Profile_Failed_LoadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Informațiile de facturare nu au putut fi încărcate`)
};

export const billing_profile_retry_loading = /** @type {(inputs: Billing_Profile_Retry_LoadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă încărcarea`)
};

export const billing_profile_billing_entity = /** @type {(inputs: Billing_Profile_Billing_EntityInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Entitate de facturare`)
};

export const billing_profile_entity_description = /** @type {(inputs: Billing_Profile_Entity_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Alege între un profil individual sau de firmă.`)
};

export const billing_profile_individual = /** @type {(inputs: Billing_Profile_IndividualInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Persoană fizică`)
};

export const billing_profile_company = /** @type {(inputs: Billing_Profile_CompanyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Companie`)
};

export const billing_profile_general_details = /** @type {(inputs: Billing_Profile_General_DetailsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Detalii generale`)
};

export const billing_profile_billing_email = /** @type {(inputs: Billing_Profile_Billing_EmailInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Email de facturare`)
};

export const billing_profile_billing_address = /** @type {(inputs: Billing_Profile_Billing_AddressInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adresă de facturare`)
};

export const billing_profile_address_line1 = /** @type {(inputs: Billing_Profile_Address_Line1Inputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adresă linia 1`)
};

export const billing_profile_address_line2 = /** @type {(inputs: Billing_Profile_Address_Line2Inputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adresă linia 2`)
};

export const billing_profile_city = /** @type {(inputs: Billing_Profile_CityInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Oraș`)
};

export const billing_profile_region_state = /** @type {(inputs: Billing_Profile_Region_StateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Regiune/județ`)
};

export const billing_profile_country = /** @type {(inputs: Billing_Profile_CountryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Țară`)
};

export const billing_profile_postal_code = /** @type {(inputs: Billing_Profile_Postal_CodeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cod poștal`)
};

export const billing_profile_company_details = /** @type {(inputs: Billing_Profile_Company_DetailsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Detalii companie`)
};

export const billing_profile_fiscal_code = /** @type {(inputs: Billing_Profile_Fiscal_CodeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Cod fiscal`)
};

export const billing_profile_registration_number = /** @type {(inputs: Billing_Profile_Registration_NumberInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Număr de înregistrare`)
};

export const billing_profile_save_button = /** @type {(inputs: Billing_Profile_Save_ButtonInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează informațiile de facturare`)
};

export const datasets_page_title = /** @type {(inputs: Datasets_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Seturi de date`)
};

export const datasets_detail_page_title = /** @type {(inputs: Datasets_Detail_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Set de date`)
};

export const datasets_name_column = /** @type {(inputs: Datasets_Name_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume`)
};

export const datasets_schema_column = /** @type {(inputs: Datasets_Schema_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schemă`)
};

export const datasets_fields_column = /** @type {(inputs: Datasets_Fields_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Câmpuri`)
};

export const datasets_created_column = /** @type {(inputs: Datasets_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creat`)
};

export const datasets_actions_column = /** @type {(inputs: Datasets_Actions_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Acțiuni`)
};

export const datasets_sort_created_ascending = /** @type {(inputs: Datasets_Sort_Created_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării ascendent`)
};

export const datasets_sort_created_descending = /** @type {(inputs: Datasets_Sort_Created_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării descendent`)
};

export const datasets_retry = /** @type {(inputs: Datasets_RetryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă`)
};

export const datasets_open = /** @type {(inputs: Datasets_OpenInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deschide`)
};

export const datasets_no_datasets_found = /** @type {(inputs: Datasets_No_Datasets_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit seturi de date.`)
};

export const datasets_showing_datasets_one = /** @type {(inputs: Datasets_Showing_Datasets_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} set de date pe această pagină.`)
};

export const datasets_showing_datasets_other = /** @type {(inputs: Datasets_Showing_Datasets_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} seturi de date pe această pagină.`)
};

export const datasets_no_datasets_to_show = /** @type {(inputs: Datasets_No_Datasets_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există seturi de date de afișat.`)
};

export const datasets_rows_per_page = /** @type {(inputs: Datasets_Rows_Per_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rânduri pe pagină`)
};

export const datasets_previous_page = /** @type {(inputs: Datasets_Previous_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagina anterioară`)
};

export const datasets_next_page = /** @type {(inputs: Datasets_Next_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagina următoare`)
};

export const datasets_field_count_one = /** @type {(inputs: Datasets_Field_Count_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmp`)
};

export const datasets_field_count_other = /** @type {(inputs: Datasets_Field_Count_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmpuri`)
};

export const datasets_date_range = /** @type {(inputs: Datasets_Date_RangeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Interval de date`)
};

export const datasets_any_date = /** @type {(inputs: Datasets_Any_DateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Oricare`)
};

export const datasets_date_range_value = /** @type {(inputs: Datasets_Date_Range_ValueInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.start} - ${i?.end}`)
};

export const datasets_presets = /** @type {(inputs: Datasets_PresetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Presetări`)
};

export const datasets_today = /** @type {(inputs: Datasets_TodayInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Astăzi`)
};

export const datasets_this_week = /** @type {(inputs: Datasets_This_WeekInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Săptămâna aceasta`)
};

export const datasets_this_month = /** @type {(inputs: Datasets_This_MonthInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Luna aceasta`)
};

export const datasets_clear = /** @type {(inputs: Datasets_ClearInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const datasets_apply = /** @type {(inputs: Datasets_ApplyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Aplică`)
};

export const datasets_document_id_column = /** @type {(inputs: Datasets_Document_Id_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID document`)
};

export const datasets_filename_column = /** @type {(inputs: Datasets_Filename_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume fișier`)
};

export const datasets_not_found_title = /** @type {(inputs: Datasets_Not_Found_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Setul de date nu a fost găsit`)
};

export const datasets_not_found_body = /** @type {(inputs: Datasets_Not_Found_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Acest set de date nu există.`)
};

export const datasets_view_datasets = /** @type {(inputs: Datasets_View_DatasetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Vezi seturile de date`)
};

export const datasets_preview_document = /** @type {(inputs: Datasets_Preview_DocumentInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Previzualizează documentul ${i?.documentId}`)
};

export const datasets_no_documents_extracted = /** @type {(inputs: Datasets_No_Documents_ExtractedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există documente extrase pentru acest set de date`)
};

export const datasets_showing_rows_one = /** @type {(inputs: Datasets_Showing_Rows_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rând pe această pagină.`)
};

export const datasets_showing_rows_other = /** @type {(inputs: Datasets_Showing_Rows_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} rânduri pe această pagină.`)
};

export const datasets_no_rows_to_show = /** @type {(inputs: Datasets_No_Rows_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există rânduri de afișat.`)
};

export const datasets_export_csv = /** @type {(inputs: Datasets_Export_CsvInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Export CSV`)
};

export const datasets_export_xlsx = /** @type {(inputs: Datasets_Export_XlsxInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Export XLSX`)
};

export const datasets_failed_export = /** @type {(inputs: Datasets_Failed_ExportInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Exportul setului de date a eșuat`)
};

export const datasets_invalid_date = /** @type {(inputs: Datasets_Invalid_DateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dată invalidă`)
};

export const datasets_missing_document_id = /** @type {(inputs: Datasets_Missing_Document_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul documentului lipsește`)
};

export const datasets_add_dataset = /** @type {(inputs: Datasets_Add_DatasetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adaugă set de date`)
};

export const datasets_all_datasets = /** @type {(inputs: Datasets_All_DatasetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate seturile de date`)
};

export const datasets_retry_datasets = /** @type {(inputs: Datasets_Retry_DatasetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă seturile de date`)
};

export const datasets_no_datasets = /** @type {(inputs: Datasets_No_DatasetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niciun set de date`)
};

export const datasets_dataset_actions = /** @type {(inputs: Datasets_Dataset_ActionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Acțiuni set de date`)
};

export const datasets_edit = /** @type {(inputs: Datasets_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează`)
};

export const datasets_delete = /** @type {(inputs: Datasets_DeleteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const datasets_delete_failed = /** @type {(inputs: Datasets_Delete_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergerea a eșuat`)
};

export const datasets_delete_confirm_title = /** @type {(inputs: Datasets_Delete_Confirm_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi setul de date?`)
};

export const datasets_delete_confirm_description = /** @type {(inputs: Datasets_Delete_Confirm_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi „${i?.name}”?`)
};

export const datasets_dialog_title_new = /** @type {(inputs: Datasets_Dialog_Title_NewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Set de date nou`)
};

export const datasets_dialog_title_edit = /** @type {(inputs: Datasets_Dialog_Title_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează setul de date`)
};

export const datasets_save_changes = /** @type {(inputs: Datasets_Save_ChangesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează modificările`)
};

export const datasets_create_dataset = /** @type {(inputs: Datasets_Create_DatasetInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează setul de date`)
};

export const datasets_selected_schema = /** @type {(inputs: Datasets_Selected_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Schema selectată`)
};

export const datasets_loading_schemas = /** @type {(inputs: Datasets_Loading_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă schemele`)
};

export const datasets_select_schema = /** @type {(inputs: Datasets_Select_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează schema`)
};

export const datasets_no_fields_selected = /** @type {(inputs: Datasets_No_Fields_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niciun câmp selectat`)
};

export const datasets_one_field_selected = /** @type {(inputs: Datasets_One_Field_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`1 câmp selectat`)
};

export const datasets_fields_selected = /** @type {(inputs: Datasets_Fields_SelectedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} câmpuri selectate`)
};

export const datasets_collapse_field = /** @type {(inputs: Datasets_Collapse_FieldInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Restrânge ${i?.label}`)
};

export const datasets_expand_field = /** @type {(inputs: Datasets_Expand_FieldInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Extinde ${i?.label}`)
};

export const datasets_select_field = /** @type {(inputs: Datasets_Select_FieldInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Selectează ${i?.label}`)
};

export const datasets_name_placeholder = /** @type {(inputs: Datasets_Name_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Numele setului de date`)
};

export const datasets_search_schemas = /** @type {(inputs: Datasets_Search_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută scheme`)
};

export const datasets_no_schemas_found = /** @type {(inputs: Datasets_No_Schemas_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit scheme.`)
};

export const datasets_no_fields = /** @type {(inputs: Datasets_No_FieldsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niciun câmp`)
};

export const datasets_cancel = /** @type {(inputs: Datasets_CancelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anulează`)
};

export const datasets_json_badge = /** @type {(inputs: Datasets_Json_BadgeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`JSON`)
};

export const documents_page_title = /** @type {(inputs: Documents_Page_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente`)
};

export const documents_new_ocr_job = /** @type {(inputs: Documents_New_Ocr_JobInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesare OCR nouă`)
};

export const documents_search_filename_placeholder = /** @type {(inputs: Documents_Search_Filename_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută după numele fișierului...`)
};

export const documents_search_filename = /** @type {(inputs: Documents_Search_FilenameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută numele fișierului`)
};

export const documents_date_range = /** @type {(inputs: Documents_Date_RangeInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Interval de date`)
};

export const documents_any_date = /** @type {(inputs: Documents_Any_DateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Oricare`)
};

export const documents_date_range_value = /** @type {(inputs: Documents_Date_Range_ValueInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.start} - ${i?.end}`)
};

export const documents_presets = /** @type {(inputs: Documents_PresetsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Presetări`)
};

export const documents_today = /** @type {(inputs: Documents_TodayInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Astăzi`)
};

export const documents_this_week = /** @type {(inputs: Documents_This_WeekInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Săptămâna aceasta`)
};

export const documents_this_month = /** @type {(inputs: Documents_This_MonthInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Luna aceasta`)
};

export const documents_clear = /** @type {(inputs: Documents_ClearInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const documents_apply = /** @type {(inputs: Documents_ApplyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Aplică`)
};

export const documents_filter_by_collection = /** @type {(inputs: Documents_Filter_By_CollectionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Filtrează după colecție`)
};

export const documents_filter_by_schema = /** @type {(inputs: Documents_Filter_By_SchemaInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Filtrează după schemă`)
};

export const documents_unknown_collection = /** @type {(inputs: Documents_Unknown_CollectionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecție necunoscută`)
};

export const documents_all_collections = /** @type {(inputs: Documents_All_CollectionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate colecțiile`)
};

export const documents_all_schemas = /** @type {(inputs: Documents_All_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate schemele`)
};

export const documents_missing_document_id = /** @type {(inputs: Documents_Missing_Document_IdInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`ID-ul documentului lipsește`)
};

export const documents_failed_load_documents = /** @type {(inputs: Documents_Failed_Load_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încărcarea documentelor a eșuat`)
};

export const documents_failed_load_document = /** @type {(inputs: Documents_Failed_Load_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Încărcarea documentului a eșuat`)
};

export const documents_failed_delete_document = /** @type {(inputs: Documents_Failed_Delete_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergerea documentului a eșuat`)
};

export const documents_failed_update_document = /** @type {(inputs: Documents_Failed_Update_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizarea documentului a eșuat`)
};

export const documents_failed_delete_documents = /** @type {(inputs: Documents_Failed_Delete_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergerea documentelor a eșuat`)
};

export const documents_failed_move_documents = /** @type {(inputs: Documents_Failed_Move_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mutarea documentelor a eșuat`)
};

export const documents_failed_download = /** @type {(inputs: Documents_Failed_DownloadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descărcarea documentelor a eșuat`)
};

export const documents_invalid_date = /** @type {(inputs: Documents_Invalid_DateInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dată invalidă`)
};

export const documents_select_all_on_page = /** @type {(inputs: Documents_Select_All_On_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Selectează toate documentele de pe această pagină`)
};

export const documents_select_document = /** @type {(inputs: Documents_Select_DocumentInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Selectează ${i?.name}`)
};

export const documents_filename_column = /** @type {(inputs: Documents_Filename_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume fișier`)
};

export const documents_collections_column = /** @type {(inputs: Documents_Collections_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecții`)
};

export const documents_pages_column = /** @type {(inputs: Documents_Pages_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Pagini`)
};

export const documents_created_column = /** @type {(inputs: Documents_Created_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creat`)
};

export const documents_file_size_column = /** @type {(inputs: Documents_File_Size_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Dimensiune fișier`)
};

export const documents_sort_created_ascending = /** @type {(inputs: Documents_Sort_Created_AscendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării ascendent`)
};

export const documents_sort_created_descending = /** @type {(inputs: Documents_Sort_Created_DescendingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Sortează după data creării descendent`)
};

export const documents_collection_not_found_title = /** @type {(inputs: Documents_Collection_Not_Found_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecția nu a fost găsită`)
};

export const documents_collection_not_found_body = /** @type {(inputs: Documents_Collection_Not_Found_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Această colecție nu există.`)
};

export const documents_view_all_documents = /** @type {(inputs: Documents_View_All_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Vezi toate documentele`)
};

export const documents_retry = /** @type {(inputs: Documents_RetryInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă`)
};

export const documents_no_documents_found = /** @type {(inputs: Documents_No_Documents_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit documente.`)
};

export const documents_empty_filtered_body = /** @type {(inputs: Documents_Empty_Filtered_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Niciun document nu corespunde criteriilor de filtrare. Încearcă să le resetezi pentru a vedea fișierele.`)
};

export const documents_empty_unfiltered_body = /** @type {(inputs: Documents_Empty_Unfiltered_BodyInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu ai procesat încă niciun document. Încarcă un document pentru a începe.`)
};

export const documents_clear_filters = /** @type {(inputs: Documents_Clear_FiltersInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge filtrele`)
};

export const documents_process_first_document = /** @type {(inputs: Documents_Process_First_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Procesează primul document`)
};

export const documents_showing_documents_one = /** @type {(inputs: Documents_Showing_Documents_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} document pe această pagină.`)
};

export const documents_showing_documents_other = /** @type {(inputs: Documents_Showing_Documents_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Se afișează ${i?.count} documente pe această pagină.`)
};

export const documents_no_documents_to_show = /** @type {(inputs: Documents_No_Documents_To_ShowInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există documente de afișat.`)
};

export const documents_rows_per_page = /** @type {(inputs: Documents_Rows_Per_PageInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Rânduri pe pagină`)
};

export const documents_previous = /** @type {(inputs: Documents_PreviousInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anterior`)
};

export const documents_next = /** @type {(inputs: Documents_NextInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Următor`)
};

export const documents_delete = /** @type {(inputs: Documents_DeleteInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Șterge`)
};

export const documents_delete_single_title = /** @type {(inputs: Documents_Delete_Single_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi documentul?`)
};

export const documents_delete_single_description = /** @type {(inputs: Documents_Delete_Single_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi „${i?.name}”? Această acțiune nu poate fi anulată.`)
};

export const documents_delete_bulk_title_one = /** @type {(inputs: Documents_Delete_Bulk_Title_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} document?`)
};

export const documents_delete_bulk_title_other = /** @type {(inputs: Documents_Delete_Bulk_Title_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} documente?`)
};

export const documents_delete_bulk_description_one = /** @type {(inputs: Documents_Delete_Bulk_Description_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} document selectat? Această acțiune nu poate fi anulată.`)
};

export const documents_delete_bulk_description_other = /** @type {(inputs: Documents_Delete_Bulk_Description_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi ${i?.count} documente selectate? Această acțiune nu poate fi anulată.`)
};

export const documents_selected_count_one = /** @type {(inputs: Documents_Selected_Count_OneInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectat`)
};

export const documents_selected_count_other = /** @type {(inputs: Documents_Selected_Count_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} selectate`)
};

export const documents_download_selected = /** @type {(inputs: Documents_Download_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descarcă documentele selectate`)
};

export const documents_download = /** @type {(inputs: Documents_DownloadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descarcă`)
};

export const documents_downloading = /** @type {(inputs: Documents_DownloadingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se descarcă...`)
};

export const documents_move = /** @type {(inputs: Documents_MoveInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mută`)
};

export const documents_moving = /** @type {(inputs: Documents_MovingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se mută...`)
};

export const documents_deleting = /** @type {(inputs: Documents_DeletingInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se șterge...`)
};

export const documents_open_actions_for = /** @type {(inputs: Documents_Open_Actions_ForInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Deschide acțiunile pentru ${i?.name}`)
};

export const documents_preview = /** @type {(inputs: Documents_PreviewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizează`)
};

export const documents_rename = /** @type {(inputs: Documents_RenameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Redenumește`)
};

export const documents_failed_rename = /** @type {(inputs: Documents_Failed_RenameInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Redenumirea documentului a eșuat`)
};

export const documents_rename_file = /** @type {(inputs: Documents_Rename_FileInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Redenumește ${i?.name}`)
};

export const documents_preview_file = /** @type {(inputs: Documents_Preview_FileInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Previzualizează ${i?.name}`)
};

export const documents_download_dialog_title_one = /** @type {(inputs: Documents_Download_Dialog_Title_OneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Descarcă documentul`)
};

export const documents_download_dialog_title_other = /** @type {(inputs: Documents_Download_Dialog_Title_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Descarcă ${i?.count} documente`)
};

export const documents_selected_documents = /** @type {(inputs: Documents_Selected_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Documente selectate`)
};

export const documents_format_markdown = /** @type {(inputs: Documents_Format_MarkdownInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Markdown`)
};

export const documents_format_html = /** @type {(inputs: Documents_Format_HtmlInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`HTML`)
};

export const documents_format_json = /** @type {(inputs: Documents_Format_JsonInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`JSON`)
};

export const documents_preparing_download = /** @type {(inputs: Documents_Preparing_DownloadInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se pregătește descărcarea...`)
};

export const documents_no_collections_selected = /** @type {(inputs: Documents_No_Collections_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio colecție selectată`)
};

export const documents_one_collection_selected = /** @type {(inputs: Documents_One_Collection_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`1 colecție selectată`)
};

export const documents_collections_selected = /** @type {(inputs: Documents_Collections_SelectedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} colecții selectate`)
};

export const documents_remove_from_all = /** @type {(inputs: Documents_Remove_From_AllInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Elimină din toate`)
};

export const documents_move_documents = /** @type {(inputs: Documents_Move_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mută documentele`)
};

export const documents_move_description_one = /** @type {(inputs: Documents_Move_Description_OneInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Înlocuiește colecțiile pentru 1 document selectat.`)
};

export const documents_move_description_other = /** @type {(inputs: Documents_Move_Description_OtherInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Înlocuiește colecțiile pentru ${i?.count} documente selectate.`)
};

export const documents_collections_label = /** @type {(inputs: Documents_Collections_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecții`)
};

export const documents_search_collections = /** @type {(inputs: Documents_Search_CollectionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută colecții`)
};

export const documents_loading_collections = /** @type {(inputs: Documents_Loading_CollectionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă colecțiile`)
};

export const documents_no_collections_found = /** @type {(inputs: Documents_No_Collections_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit colecții.`)
};

export const documents_cancel = /** @type {(inputs: Documents_CancelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Anulează`)
};

export const documents_collections_nav_label = /** @type {(inputs: Documents_Collections_Nav_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecții`)
};

export const documents_add_collection = /** @type {(inputs: Documents_Add_CollectionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Adaugă colecție`)
};

export const documents_all_documents = /** @type {(inputs: Documents_All_DocumentsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Toate documentele`)
};

export const documents_retry_collections = /** @type {(inputs: Documents_Retry_CollectionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Reîncearcă colecțiile`)
};

export const documents_no_collections = /** @type {(inputs: Documents_No_CollectionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio colecție`)
};

export const documents_collection_actions = /** @type {(inputs: Documents_Collection_ActionsInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Acțiuni colecție`)
};

export const documents_edit = /** @type {(inputs: Documents_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează`)
};

export const documents_delete_failed = /** @type {(inputs: Documents_Delete_FailedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergerea a eșuat`)
};

export const documents_delete_collection_title = /** @type {(inputs: Documents_Delete_Collection_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Ștergi colecția?`)
};

export const documents_delete_collection_description = /** @type {(inputs: Documents_Delete_Collection_DescriptionInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`Ștergi „${i?.name}”? Documentele rămân disponibile în Toate documentele.`)
};

export const documents_collection_dialog_title_new = /** @type {(inputs: Documents_Collection_Dialog_Title_NewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Colecție nouă`)
};

export const documents_collection_dialog_title_edit = /** @type {(inputs: Documents_Collection_Dialog_Title_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Editează colecția`)
};

export const documents_collection_dialog_description_new = /** @type {(inputs: Documents_Collection_Dialog_Description_NewInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Grupează documentele după nume și filtre opționale de schemă.`)
};

export const documents_collection_dialog_description_edit = /** @type {(inputs: Documents_Collection_Dialog_Description_EditInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Actualizează numele colecției și filtrele de schemă.`)
};

export const documents_save_changes = /** @type {(inputs: Documents_Save_ChangesInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Salvează modificările`)
};

export const documents_create_collection = /** @type {(inputs: Documents_Create_CollectionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Creează colecția`)
};

export const documents_name_column = /** @type {(inputs: Documents_Name_ColumnInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nume`)
};

export const documents_collection_name_placeholder = /** @type {(inputs: Documents_Collection_Name_PlaceholderInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Facturi, rapoarte, chitanțe`)
};

export const documents_schemas_label = /** @type {(inputs: Documents_Schemas_LabelInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Scheme`)
};

export const documents_no_schemas_selected = /** @type {(inputs: Documents_No_Schemas_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nicio schemă selectată`)
};

export const documents_one_schema_selected = /** @type {(inputs: Documents_One_Schema_SelectedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`1 schemă selectată`)
};

export const documents_schemas_selected = /** @type {(inputs: Documents_Schemas_SelectedInputs) => LocalizedString} */ (i) => {
	return /** @type {LocalizedString} */ (`${i?.count} scheme selectate`)
};

export const documents_search_schemas = /** @type {(inputs: Documents_Search_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Caută scheme`)
};

export const documents_loading_schemas = /** @type {(inputs: Documents_Loading_SchemasInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă schemele`)
};

export const documents_no_schemas_found = /** @type {(inputs: Documents_No_Schemas_FoundInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu s-au găsit scheme.`)
};

export const documents_collection_schema_hint = /** @type {(inputs: Documents_Collection_Schema_HintInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`În această colecție sunt afișate doar documentele cu scheme potrivite.`)
};

export const documents_preview_fallback_title = /** @type {(inputs: Documents_Preview_Fallback_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Previzualizare document`)
};

export const documents_preview_description = /** @type {(inputs: Documents_Preview_DescriptionInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Revizuiește Markdown-ul extras și JSON-ul pentru acest document.`)
};

export const documents_rename_document_title = /** @type {(inputs: Documents_Rename_Document_TitleInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Redenumește documentul`)
};

export const documents_loading_document = /** @type {(inputs: Documents_Loading_DocumentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Se încarcă documentul...`)
};

export const documents_copy_markdown = /** @type {(inputs: Documents_Copy_MarkdownInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Copiază Markdown`)
};

export const documents_copy_html = /** @type {(inputs: Documents_Copy_HtmlInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Copiază HTML`)
};

export const documents_copy_json = /** @type {(inputs: Documents_Copy_JsonInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Copiază JSON`)
};

export const documents_copied = /** @type {(inputs: Documents_CopiedInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Copiat`)
};

export const documents_no_json_annotation = /** @type {(inputs: Documents_No_Json_AnnotationInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu este disponibilă nicio adnotare JSON.`)
};

export const documents_no_markdown_content = /** @type {(inputs: Documents_No_Markdown_ContentInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu există conținut Markdown disponibil.`)
};

export const documents_no_preview_available = /** @type {(inputs: Documents_No_Preview_AvailableInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Nu este disponibilă nicio previzualizare a documentului.`)
};

export const documents_close = /** @type {(inputs: Documents_CloseInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Închide`)
};

export const documents_more = /** @type {(inputs: Documents_MoreInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Mai multe`)
};

export const documents_open = /** @type {(inputs: Documents_OpenInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Deschide`)
};

export const documents_share = /** @type {(inputs: Documents_ShareInputs) => LocalizedString} */ () => {
	return /** @type {LocalizedString} */ (`Distribuie`)
};
