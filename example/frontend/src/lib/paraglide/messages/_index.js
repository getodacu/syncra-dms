/* eslint-disable */
import { getLocale, experimentalStaticLocale } from "../runtime.js"

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
import * as __en from "./en.js"
import * as __ro from "./ro.js"
/**
* | output |
* | --- |
* | "Language" |
*
* @param {Common_LanguageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_language = /** @type {((inputs?: Common_LanguageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_LanguageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_language(inputs)
	return __ro.common_language(inputs)
});
/**
* | output |
* | --- |
* | "English" |
*
* @param {Common_EnglishInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_english = /** @type {((inputs?: Common_EnglishInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_EnglishInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_english(inputs)
	return __ro.common_english(inputs)
});
/**
* | output |
* | --- |
* | "Romanian" |
*
* @param {Common_RomanianInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_romanian = /** @type {((inputs?: Common_RomanianInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_RomanianInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_romanian(inputs)
	return __ro.common_romanian(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Common_CancelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_cancel = /** @type {((inputs?: Common_CancelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_CancelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_cancel(inputs)
	return __ro.common_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Common_DeleteInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_delete = /** @type {((inputs?: Common_DeleteInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_DeleteInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_delete(inputs)
	return __ro.common_delete(inputs)
});
/**
* | output |
* | --- |
* | "Retry" |
*
* @param {Common_RetryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_retry = /** @type {((inputs?: Common_RetryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_RetryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_retry(inputs)
	return __ro.common_retry(inputs)
});
/**
* | output |
* | --- |
* | "Previous" |
*
* @param {Common_PreviousInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_previous = /** @type {((inputs?: Common_PreviousInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_PreviousInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_previous(inputs)
	return __ro.common_previous(inputs)
});
/**
* | output |
* | --- |
* | "Next" |
*
* @param {Common_NextInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_next = /** @type {((inputs?: Common_NextInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_NextInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_next(inputs)
	return __ro.common_next(inputs)
});
/**
* | output |
* | --- |
* | "Rows per page" |
*
* @param {Common_Rows_Per_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_rows_per_page = /** @type {((inputs?: Common_Rows_Per_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_Rows_Per_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_rows_per_page(inputs)
	return __ro.common_rows_per_page(inputs)
});
/**
* | output |
* | --- |
* | "Strict" |
*
* @param {Common_StrictInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_strict = /** @type {((inputs?: Common_StrictInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_StrictInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_strict(inputs)
	return __ro.common_strict(inputs)
});
/**
* | output |
* | --- |
* | "Flexible" |
*
* @param {Common_FlexibleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_flexible = /** @type {((inputs?: Common_FlexibleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_FlexibleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_flexible(inputs)
	return __ro.common_flexible(inputs)
});
/**
* | output |
* | --- |
* | "Required" |
*
* @param {Common_RequiredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_required = /** @type {((inputs?: Common_RequiredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_RequiredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_required(inputs)
	return __ro.common_required(inputs)
});
/**
* | output |
* | --- |
* | "Unknown" |
*
* @param {Common_UnknownInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_unknown = /** @type {((inputs?: Common_UnknownInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_UnknownInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_unknown(inputs)
	return __ro.common_unknown(inputs)
});
/**
* | output |
* | --- |
* | "Actions" |
*
* @param {Common_ActionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_actions = /** @type {((inputs?: Common_ActionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ActionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_actions(inputs)
	return __ro.common_actions(inputs)
});
/**
* | output |
* | --- |
* | "Toggle theme" |
*
* @param {Common_Toggle_ThemeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_toggle_theme = /** @type {((inputs?: Common_Toggle_ThemeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_Toggle_ThemeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_toggle_theme(inputs)
	return __ro.common_toggle_theme(inputs)
});
/**
* | output |
* | --- |
* | "Credits unavailable" |
*
* @param {Header_Credits_UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const header_credits_unavailable = /** @type {((inputs?: Header_Credits_UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Header_Credits_UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.header_credits_unavailable(inputs)
	return __ro.header_credits_unavailable(inputs)
});
/**
* | output |
* | --- |
* | "{count} credits" |
*
* @param {Header_CreditsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const header_credits = /** @type {((inputs: Header_CreditsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Header_CreditsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.header_credits(inputs)
	return __ro.header_credits(inputs)
});
/**
* | output |
* | --- |
* | "Credit balance unavailable: {message}" |
*
* @param {Header_Credit_Balance_UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const header_credit_balance_unavailable = /** @type {((inputs: Header_Credit_Balance_UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Header_Credit_Balance_UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.header_credit_balance_unavailable(inputs)
	return __ro.header_credit_balance_unavailable(inputs)
});
/**
* | output |
* | --- |
* | "Account" |
*
* @param {Nav_AccountInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account = /** @type {((inputs?: Nav_AccountInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_AccountInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account(inputs)
	return __ro.nav_account(inputs)
});
/**
* | output |
* | --- |
* | "No email address" |
*
* @param {Nav_No_Email_AddressInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_no_email_address = /** @type {((inputs?: Nav_No_Email_AddressInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_No_Email_AddressInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_no_email_address(inputs)
	return __ro.nav_no_email_address(inputs)
});
/**
* | output |
* | --- |
* | "Notifications" |
*
* @param {Nav_NotificationsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_notifications = /** @type {((inputs?: Nav_NotificationsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_NotificationsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_notifications(inputs)
	return __ro.nav_notifications(inputs)
});
/**
* | output |
* | --- |
* | "Log out" |
*
* @param {Nav_Log_OutInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_log_out = /** @type {((inputs?: Nav_Log_OutInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Log_OutInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_log_out(inputs)
	return __ro.nav_log_out(inputs)
});
/**
* | output |
* | --- |
* | "Log out?" |
*
* @param {Nav_Logout_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_logout_title = /** @type {((inputs?: Nav_Logout_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Logout_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_logout_title(inputs)
	return __ro.nav_logout_title(inputs)
});
/**
* | output |
* | --- |
* | "Are you sure you want to log out of Syncra?" |
*
* @param {Nav_Logout_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_logout_description = /** @type {((inputs?: Nav_Logout_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Logout_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_logout_description(inputs)
	return __ro.nav_logout_description(inputs)
});
/**
* | output |
* | --- |
* | "Unable to log out. Please try again." |
*
* @param {Nav_Logout_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_logout_failed = /** @type {((inputs?: Nav_Logout_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Logout_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_logout_failed(inputs)
	return __ro.nav_logout_failed(inputs)
});
/**
* | output |
* | --- |
* | "{provider} account linked." |
*
* @param {Nav_Account_LinkedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_linked = /** @type {((inputs: Nav_Account_LinkedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_LinkedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_linked(inputs)
	return __ro.nav_account_linked(inputs)
});
/**
* | output |
* | --- |
* | "{provider} is already linked to another account." |
*
* @param {Nav_Account_Link_ConflictInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_link_conflict = /** @type {((inputs: Nav_Account_Link_ConflictInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_Link_ConflictInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_link_conflict(inputs)
	return __ro.nav_account_link_conflict(inputs)
});
/**
* | output |
* | --- |
* | "{provider} linking was cancelled." |
*
* @param {Nav_Account_Link_DeniedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_link_denied = /** @type {((inputs: Nav_Account_Link_DeniedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_Link_DeniedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_link_denied(inputs)
	return __ro.nav_account_link_denied(inputs)
});
/**
* | output |
* | --- |
* | "{provider} linking is not configured." |
*
* @param {Nav_Account_Link_Not_ConfiguredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_link_not_configured = /** @type {((inputs: Nav_Account_Link_Not_ConfiguredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_Link_Not_ConfiguredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_link_not_configured(inputs)
	return __ro.nav_account_link_not_configured(inputs)
});
/**
* | output |
* | --- |
* | "Please sign in again before linking accounts." |
*
* @param {Nav_Account_Link_Sign_In_AgainInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_link_sign_in_again = /** @type {((inputs?: Nav_Account_Link_Sign_In_AgainInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_Link_Sign_In_AgainInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_link_sign_in_again(inputs)
	return __ro.nav_account_link_sign_in_again(inputs)
});
/**
* | output |
* | --- |
* | "Unable to link {provider}." |
*
* @param {Nav_Account_Link_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_account_link_failed = /** @type {((inputs: Nav_Account_Link_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Account_Link_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_account_link_failed(inputs)
	return __ro.nav_account_link_failed(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard" |
*
* @param {Nav_DashboardInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_dashboard = /** @type {((inputs?: Nav_DashboardInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_DashboardInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_dashboard(inputs)
	return __ro.nav_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Schemas" |
*
* @param {Nav_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_schemas = /** @type {((inputs?: Nav_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_schemas(inputs)
	return __ro.nav_schemas(inputs)
});
/**
* | output |
* | --- |
* | "New schema" |
*
* @param {Nav_New_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_new_schema = /** @type {((inputs?: Nav_New_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_New_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_new_schema(inputs)
	return __ro.nav_new_schema(inputs)
});
/**
* | output |
* | --- |
* | "Edit schema" |
*
* @param {Nav_Edit_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_edit_schema = /** @type {((inputs?: Nav_Edit_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Edit_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_edit_schema(inputs)
	return __ro.nav_edit_schema(inputs)
});
/**
* | output |
* | --- |
* | "Jobs" |
*
* @param {Nav_JobsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_jobs = /** @type {((inputs?: Nav_JobsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_JobsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_jobs(inputs)
	return __ro.nav_jobs(inputs)
});
/**
* | output |
* | --- |
* | "New Job" |
*
* @param {Nav_New_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_new_job = /** @type {((inputs?: Nav_New_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_New_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_new_job(inputs)
	return __ro.nav_new_job(inputs)
});
/**
* | output |
* | --- |
* | "Billing" |
*
* @param {Nav_BillingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_billing = /** @type {((inputs?: Nav_BillingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_BillingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_billing(inputs)
	return __ro.nav_billing(inputs)
});
/**
* | output |
* | --- |
* | "Billing Orders" |
*
* @param {Nav_Billing_OrdersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_billing_orders = /** @type {((inputs?: Nav_Billing_OrdersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Billing_OrdersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_billing_orders(inputs)
	return __ro.nav_billing_orders(inputs)
});
/**
* | output |
* | --- |
* | "Credit Usage History" |
*
* @param {Nav_Credit_Usage_HistoryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_credit_usage_history = /** @type {((inputs?: Nav_Credit_Usage_HistoryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Credit_Usage_HistoryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_credit_usage_history(inputs)
	return __ro.nav_credit_usage_history(inputs)
});
/**
* | output |
* | --- |
* | "Developer Settings" |
*
* @param {Nav_Developer_SettingsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_developer_settings = /** @type {((inputs?: Nav_Developer_SettingsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Developer_SettingsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_developer_settings(inputs)
	return __ro.nav_developer_settings(inputs)
});
/**
* | output |
* | --- |
* | "Get Help" |
*
* @param {Nav_Get_HelpInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_get_help = /** @type {((inputs?: Nav_Get_HelpInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Get_HelpInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_get_help(inputs)
	return __ro.nav_get_help(inputs)
});
/**
* | output |
* | --- |
* | "Quick OCR" |
*
* @param {Nav_Quick_OcrInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_quick_ocr = /** @type {((inputs?: Nav_Quick_OcrInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Quick_OcrInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_quick_ocr(inputs)
	return __ro.nav_quick_ocr(inputs)
});
/**
* | output |
* | --- |
* | "Create Quick OCR job" |
*
* @param {Nav_Create_Quick_Ocr_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_create_quick_ocr_job = /** @type {((inputs?: Nav_Create_Quick_Ocr_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Create_Quick_Ocr_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_create_quick_ocr_job(inputs)
	return __ro.nav_create_quick_ocr_job(inputs)
});
/**
* | output |
* | --- |
* | "Create schema" |
*
* @param {Nav_Create_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_create_schema = /** @type {((inputs?: Nav_Create_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Create_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_create_schema(inputs)
	return __ro.nav_create_schema(inputs)
});
/**
* | output |
* | --- |
* | "Create job" |
*
* @param {Nav_Create_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const nav_create_job = /** @type {((inputs?: Nav_Create_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Create_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_create_job(inputs)
	return __ro.nav_create_job(inputs)
});
/**
* | output |
* | --- |
* | "Documents processed" |
*
* @param {Dashboard_Metric_Documents_ProcessedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_metric_documents_processed = /** @type {((inputs?: Dashboard_Metric_Documents_ProcessedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Metric_Documents_ProcessedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_metric_documents_processed(inputs)
	return __ro.dashboard_metric_documents_processed(inputs)
});
/**
* | output |
* | --- |
* | "Throughput, recent work, datasets, and credits in one place." |
*
* @param {Dashboard_Page_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_page_description = /** @type {((inputs?: Dashboard_Page_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Page_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_page_description(inputs)
	return __ro.dashboard_page_description(inputs)
});
/**
* | output |
* | --- |
* | "Refreshing" |
*
* @param {Dashboard_RefreshingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_refreshing = /** @type {((inputs?: Dashboard_RefreshingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_RefreshingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_refreshing(inputs)
	return __ro.dashboard_refreshing(inputs)
});
/**
* | output |
* | --- |
* | "Loading dashboard" |
*
* @param {Dashboard_Loading_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_loading_title = /** @type {((inputs?: Dashboard_Loading_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Loading_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_loading_title(inputs)
	return __ro.dashboard_loading_title(inputs)
});
/**
* | output |
* | --- |
* | "Preparing your workspace overview." |
*
* @param {Dashboard_Loading_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_loading_description = /** @type {((inputs?: Dashboard_Loading_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Loading_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_loading_description(inputs)
	return __ro.dashboard_loading_description(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard partially loaded" |
*
* @param {Dashboard_Warning_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_warning_title = /** @type {((inputs?: Dashboard_Warning_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Warning_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_warning_title(inputs)
	return __ro.dashboard_warning_title(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard unavailable" |
*
* @param {Dashboard_Unavailable_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_unavailable_title = /** @type {((inputs?: Dashboard_Unavailable_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Unavailable_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_unavailable_title(inputs)
	return __ro.dashboard_unavailable_title(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard data could not be loaded." |
*
* @param {Dashboard_Unavailable_DefaultInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_unavailable_default = /** @type {((inputs?: Dashboard_Unavailable_DefaultInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Unavailable_DefaultInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_unavailable_default(inputs)
	return __ro.dashboard_unavailable_default(inputs)
});
/**
* | output |
* | --- |
* | "Pages processed" |
*
* @param {Dashboard_Metric_Pages_ProcessedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_metric_pages_processed = /** @type {((inputs?: Dashboard_Metric_Pages_ProcessedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Metric_Pages_ProcessedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_metric_pages_processed(inputs)
	return __ro.dashboard_metric_pages_processed(inputs)
});
/**
* | output |
* | --- |
* | "Completion rate" |
*
* @param {Dashboard_Metric_Completion_RateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_metric_completion_rate = /** @type {((inputs?: Dashboard_Metric_Completion_RateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Metric_Completion_RateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_metric_completion_rate(inputs)
	return __ro.dashboard_metric_completion_rate(inputs)
});
/**
* | output |
* | --- |
* | "Credits spent" |
*
* @param {Dashboard_Metric_Credits_SpentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_metric_credits_spent = /** @type {((inputs?: Dashboard_Metric_Credits_SpentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Metric_Credits_SpentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_metric_credits_spent(inputs)
	return __ro.dashboard_metric_credits_spent(inputs)
});
/**
* | output |
* | --- |
* | "{count} job in progress" |
*
* @param {Dashboard_Jobs_In_Progress_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_jobs_in_progress_one = /** @type {((inputs: Dashboard_Jobs_In_Progress_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Jobs_In_Progress_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_jobs_in_progress_one(inputs)
	return __ro.dashboard_jobs_in_progress_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} jobs in progress" |
*
* @param {Dashboard_Jobs_In_Progress_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_jobs_in_progress_other = /** @type {((inputs: Dashboard_Jobs_In_Progress_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Jobs_In_Progress_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_jobs_in_progress_other(inputs)
	return __ro.dashboard_jobs_in_progress_other(inputs)
});
/**
* | output |
* | --- |
* | "OCR pages completed" |
*
* @param {Dashboard_Pages_CompletedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_pages_completed = /** @type {((inputs?: Dashboard_Pages_CompletedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pages_CompletedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pages_completed(inputs)
	return __ro.dashboard_pages_completed(inputs)
});
/**
* | output |
* | --- |
* | "{completed} completed, {failed} failed" |
*
* @param {Dashboard_Completion_SummaryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_completion_summary = /** @type {((inputs: Dashboard_Completion_SummaryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Completion_SummaryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_completion_summary(inputs)
	return __ro.dashboard_completion_summary(inputs)
});
/**
* | output |
* | --- |
* | "{count} available" |
*
* @param {Dashboard_Credits_Available_ShortInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_available_short = /** @type {((inputs: Dashboard_Credits_Available_ShortInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_Available_ShortInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_available_short(inputs)
	return __ro.dashboard_credits_available_short(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard metrics" |
*
* @param {Dashboard_Metrics_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_metrics_aria = /** @type {((inputs?: Dashboard_Metrics_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Metrics_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_metrics_aria(inputs)
	return __ro.dashboard_metrics_aria(inputs)
});
/**
* | output |
* | --- |
* | "Documents processed" |
*
* @param {Dashboard_Documents_Processed_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_documents_processed_title = /** @type {((inputs?: Dashboard_Documents_Processed_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Documents_Processed_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_documents_processed_title(inputs)
	return __ro.dashboard_documents_processed_title(inputs)
});
/**
* | output |
* | --- |
* | "Documents" |
*
* @param {Dashboard_Chart_Documents_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_chart_documents_label = /** @type {((inputs?: Dashboard_Chart_Documents_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Chart_Documents_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_chart_documents_label(inputs)
	return __ro.dashboard_chart_documents_label(inputs)
});
/**
* | output |
* | --- |
* | "Select dashboard range" |
*
* @param {Dashboard_Select_RangeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_range = /** @type {((inputs?: Dashboard_Select_RangeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_RangeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_range(inputs)
	return __ro.dashboard_select_range(inputs)
});
/**
* | output |
* | --- |
* | "Last 7 days" |
*
* @param {Dashboard_Range_7dInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_range_7d = /** @type {((inputs?: Dashboard_Range_7dInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Range_7dInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_range_7d(inputs)
	return __ro.dashboard_range_7d(inputs)
});
/**
* | output |
* | --- |
* | "Last 30 days" |
*
* @param {Dashboard_Range_30dInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_range_30d = /** @type {((inputs?: Dashboard_Range_30dInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Range_30dInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_range_30d(inputs)
	return __ro.dashboard_range_30d(inputs)
});
/**
* | output |
* | --- |
* | "Last 90 days" |
*
* @param {Dashboard_Range_90dInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_range_90d = /** @type {((inputs?: Dashboard_Range_90dInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Range_90dInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_range_90d(inputs)
	return __ro.dashboard_range_90d(inputs)
});
/**
* | output |
* | --- |
* | "Recent documents" |
*
* @param {Dashboard_Recent_Documents_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_recent_documents_title = /** @type {((inputs?: Dashboard_Recent_Documents_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Recent_Documents_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_recent_documents_title(inputs)
	return __ro.dashboard_recent_documents_title(inputs)
});
/**
* | output |
* | --- |
* | "Latest completed OCR documents" |
*
* @param {Dashboard_Recent_Documents_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_recent_documents_description = /** @type {((inputs?: Dashboard_Recent_Documents_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Recent_Documents_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_recent_documents_description(inputs)
	return __ro.dashboard_recent_documents_description(inputs)
});
/**
* | output |
* | --- |
* | "View" |
*
* @param {Dashboard_ViewInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_view = /** @type {((inputs?: Dashboard_ViewInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ViewInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_view(inputs)
	return __ro.dashboard_view(inputs)
});
/**
* | output |
* | --- |
* | "No saved schema" |
*
* @param {Dashboard_No_Saved_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_saved_schema = /** @type {((inputs?: Dashboard_No_Saved_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Saved_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_saved_schema(inputs)
	return __ro.dashboard_no_saved_schema(inputs)
});
/**
* | output |
* | --- |
* | "{count} page" |
*
* @param {Dashboard_Pages_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_pages_one = /** @type {((inputs: Dashboard_Pages_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pages_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pages_one(inputs)
	return __ro.dashboard_pages_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} pages" |
*
* @param {Dashboard_Pages_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_pages_other = /** @type {((inputs: Dashboard_Pages_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pages_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pages_other(inputs)
	return __ro.dashboard_pages_other(inputs)
});
/**
* | output |
* | --- |
* | "No completed documents yet" |
*
* @param {Dashboard_No_Completed_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_completed_documents = /** @type {((inputs?: Dashboard_No_Completed_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Completed_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_completed_documents(inputs)
	return __ro.dashboard_no_completed_documents(inputs)
});
/**
* | output |
* | --- |
* | "Schema throughput" |
*
* @param {Dashboard_Schema_Throughput_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_schema_throughput_title = /** @type {((inputs?: Dashboard_Schema_Throughput_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Schema_Throughput_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_schema_throughput_title(inputs)
	return __ro.dashboard_schema_throughput_title(inputs)
});
/**
* | output |
* | --- |
* | "Completed documents by schema" |
*
* @param {Dashboard_Schema_Throughput_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_schema_throughput_description = /** @type {((inputs?: Dashboard_Schema_Throughput_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Schema_Throughput_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_schema_throughput_description(inputs)
	return __ro.dashboard_schema_throughput_description(inputs)
});
/**
* | output |
* | --- |
* | "{count} document processed" |
*
* @param {Dashboard_Documents_Processed_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_documents_processed_one = /** @type {((inputs: Dashboard_Documents_Processed_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Documents_Processed_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_documents_processed_one(inputs)
	return __ro.dashboard_documents_processed_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} documents processed" |
*
* @param {Dashboard_Documents_Processed_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_documents_processed_other = /** @type {((inputs: Dashboard_Documents_Processed_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Documents_Processed_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_documents_processed_other(inputs)
	return __ro.dashboard_documents_processed_other(inputs)
});
/**
* | output |
* | --- |
* | "No schema throughput in this range" |
*
* @param {Dashboard_No_Schema_ThroughputInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_schema_throughput = /** @type {((inputs?: Dashboard_No_Schema_ThroughputInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Schema_ThroughputInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_schema_throughput(inputs)
	return __ro.dashboard_no_schema_throughput(inputs)
});
/**
* | output |
* | --- |
* | "Datasets" |
*
* @param {Dashboard_Datasets_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_datasets_title = /** @type {((inputs?: Dashboard_Datasets_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Datasets_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_datasets_title(inputs)
	return __ro.dashboard_datasets_title(inputs)
});
/**
* | output |
* | --- |
* | "{count} total dataset" |
*
* @param {Dashboard_Total_Datasets_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_total_datasets_one = /** @type {((inputs: Dashboard_Total_Datasets_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Total_Datasets_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_total_datasets_one(inputs)
	return __ro.dashboard_total_datasets_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} total datasets" |
*
* @param {Dashboard_Total_Datasets_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_total_datasets_other = /** @type {((inputs: Dashboard_Total_Datasets_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Total_Datasets_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_total_datasets_other(inputs)
	return __ro.dashboard_total_datasets_other(inputs)
});
/**
* | output |
* | --- |
* | "{count} field" |
*
* @param {Dashboard_Fields_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_fields_one = /** @type {((inputs: Dashboard_Fields_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Fields_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_fields_one(inputs)
	return __ro.dashboard_fields_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} fields" |
*
* @param {Dashboard_Fields_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_fields_other = /** @type {((inputs: Dashboard_Fields_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Fields_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_fields_other(inputs)
	return __ro.dashboard_fields_other(inputs)
});
/**
* | output |
* | --- |
* | "No datasets yet" |
*
* @param {Dashboard_No_DatasetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_datasets = /** @type {((inputs?: Dashboard_No_DatasetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_DatasetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_datasets(inputs)
	return __ro.dashboard_no_datasets(inputs)
});
/**
* | output |
* | --- |
* | "Credits" |
*
* @param {Dashboard_Credits_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_title = /** @type {((inputs?: Dashboard_Credits_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_title(inputs)
	return __ro.dashboard_credits_title(inputs)
});
/**
* | output |
* | --- |
* | "Balance and usage in range" |
*
* @param {Dashboard_Credits_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_description = /** @type {((inputs?: Dashboard_Credits_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_description(inputs)
	return __ro.dashboard_credits_description(inputs)
});
/**
* | output |
* | --- |
* | "Low credit" |
*
* @param {Dashboard_Low_CreditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_low_credit = /** @type {((inputs?: Dashboard_Low_CreditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Low_CreditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_low_credit(inputs)
	return __ro.dashboard_low_credit(inputs)
});
/**
* | output |
* | --- |
* | "Available credits" |
*
* @param {Dashboard_Available_CreditsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_available_credits = /** @type {((inputs?: Dashboard_Available_CreditsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Available_CreditsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_available_credits(inputs)
	return __ro.dashboard_available_credits(inputs)
});
/**
* | output |
* | --- |
* | "Credits spent in selected range" |
*
* @param {Dashboard_Credits_Spent_In_RangeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_spent_in_range = /** @type {((inputs?: Dashboard_Credits_Spent_In_RangeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_Spent_In_RangeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_spent_in_range(inputs)
	return __ro.dashboard_credits_spent_in_range(inputs)
});
/**
* | output |
* | --- |
* | "Billing" |
*
* @param {Dashboard_BillingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_billing = /** @type {((inputs?: Dashboard_BillingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_BillingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_billing(inputs)
	return __ro.dashboard_billing(inputs)
});
/**
* | output |
* | --- |
* | "Start processing documents" |
*
* @param {Dashboard_Onboarding_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_onboarding_title = /** @type {((inputs?: Dashboard_Onboarding_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Onboarding_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_onboarding_title(inputs)
	return __ro.dashboard_onboarding_title(inputs)
});
/**
* | output |
* | --- |
* | "Create a schema, run OCR, then turn results into datasets." |
*
* @param {Dashboard_Onboarding_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_onboarding_description = /** @type {((inputs?: Dashboard_Onboarding_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Onboarding_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_onboarding_description(inputs)
	return __ro.dashboard_onboarding_description(inputs)
});
/**
* | output |
* | --- |
* | "New OCR job" |
*
* @param {Dashboard_New_Ocr_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_new_ocr_job = /** @type {((inputs?: Dashboard_New_Ocr_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_New_Ocr_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_new_ocr_job(inputs)
	return __ro.dashboard_new_ocr_job(inputs)
});
/**
* | output |
* | --- |
* | "{count} credit" |
*
* @param {Dashboard_Credits_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_one = /** @type {((inputs: Dashboard_Credits_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_one(inputs)
	return __ro.dashboard_credits_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} credits" |
*
* @param {Dashboard_Credits_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_credits_other = /** @type {((inputs: Dashboard_Credits_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Credits_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_credits_other(inputs)
	return __ro.dashboard_credits_other(inputs)
});
/**
* | output |
* | --- |
* | "Schema" |
*
* @param {Dashboard_Step_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_schema = /** @type {((inputs?: Dashboard_Step_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_schema(inputs)
	return __ro.dashboard_step_schema(inputs)
});
/**
* | output |
* | --- |
* | "OCR job" |
*
* @param {Dashboard_Step_Ocr_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_ocr_job = /** @type {((inputs?: Dashboard_Step_Ocr_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_Ocr_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_ocr_job(inputs)
	return __ro.dashboard_step_ocr_job(inputs)
});
/**
* | output |
* | --- |
* | "Dataset" |
*
* @param {Dashboard_Step_DatasetInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_dataset = /** @type {((inputs?: Dashboard_Step_DatasetInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_DatasetInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_dataset(inputs)
	return __ro.dashboard_step_dataset(inputs)
});
/**
* | output |
* | --- |
* | "API key" |
*
* @param {Dashboard_Step_Api_KeyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_api_key = /** @type {((inputs?: Dashboard_Step_Api_KeyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_Api_KeyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_api_key(inputs)
	return __ro.dashboard_step_api_key(inputs)
});
/**
* | output |
* | --- |
* | "Webhook" |
*
* @param {Dashboard_Step_WebhookInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_webhook = /** @type {((inputs?: Dashboard_Step_WebhookInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_WebhookInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_webhook(inputs)
	return __ro.dashboard_step_webhook(inputs)
});
/**
* | output |
* | --- |
* | "Ready" |
*
* @param {Dashboard_Step_ReadyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_ready = /** @type {((inputs?: Dashboard_Step_ReadyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_ReadyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_ready(inputs)
	return __ro.dashboard_step_ready(inputs)
});
/**
* | output |
* | --- |
* | "Open" |
*
* @param {Dashboard_Step_OpenInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const dashboard_step_open = /** @type {((inputs?: Dashboard_Step_OpenInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Step_OpenInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_step_open(inputs)
	return __ro.dashboard_step_open(inputs)
});
/**
* | output |
* | --- |
* | "Users" |
*
* @param {Admin_Nav_UsersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_users = /** @type {((inputs?: Admin_Nav_UsersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_UsersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_users(inputs)
	return __ro.admin_nav_users(inputs)
});
/**
* | output |
* | --- |
* | "User" |
*
* @param {Admin_Nav_UserInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_user = /** @type {((inputs?: Admin_Nav_UserInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_UserInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_user(inputs)
	return __ro.admin_nav_user(inputs)
});
/**
* | output |
* | --- |
* | "Invoices" |
*
* @param {Admin_Nav_InvoicesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_invoices = /** @type {((inputs?: Admin_Nav_InvoicesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_InvoicesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_invoices(inputs)
	return __ro.admin_nav_invoices(inputs)
});
/**
* | output |
* | --- |
* | "Orders" |
*
* @param {Admin_Nav_OrdersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_orders = /** @type {((inputs?: Admin_Nav_OrdersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_OrdersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_orders(inputs)
	return __ro.admin_nav_orders(inputs)
});
/**
* | output |
* | --- |
* | "Recipes" |
*
* @param {Admin_Nav_Json_RecipesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_json_recipes = /** @type {((inputs?: Admin_Nav_Json_RecipesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_Json_RecipesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_json_recipes(inputs)
	return __ro.admin_nav_json_recipes(inputs)
});
/**
* | output |
* | --- |
* | "Admin" |
*
* @param {Admin_Nav_AdminInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_nav_admin = /** @type {((inputs?: Admin_Nav_AdminInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_Nav_AdminInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_nav_admin(inputs)
	return __ro.admin_nav_admin(inputs)
});
/**
* | output |
* | --- |
* | "Admin" |
*
* @param {Admin_User_FallbackInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const admin_user_fallback = /** @type {((inputs?: Admin_User_FallbackInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Admin_User_FallbackInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.admin_user_fallback(inputs)
	return __ro.admin_user_fallback(inputs)
});
/**
* | output |
* | --- |
* | "Syncra" |
*
* @param {Sidebar_SyncraInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const sidebar_syncra = /** @type {((inputs?: Sidebar_SyncraInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_SyncraInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.sidebar_syncra(inputs)
	return __ro.sidebar_syncra(inputs)
});
/**
* | output |
* | --- |
* | "Syncra Admin" |
*
* @param {Sidebar_Syncra_AdminInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const sidebar_syncra_admin = /** @type {((inputs?: Sidebar_Syncra_AdminInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_Syncra_AdminInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.sidebar_syncra_admin(inputs)
	return __ro.sidebar_syncra_admin(inputs)
});
/**
* | output |
* | --- |
* | "User Space" |
*
* @param {Sidebar_User_SpaceInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const sidebar_user_space = /** @type {((inputs?: Sidebar_User_SpaceInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_User_SpaceInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.sidebar_user_space(inputs)
	return __ro.sidebar_user_space(inputs)
});
/**
* | output |
* | --- |
* | "Admin Portal" |
*
* @param {Sidebar_Admin_PortalInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const sidebar_admin_portal = /** @type {((inputs?: Sidebar_Admin_PortalInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_Admin_PortalInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.sidebar_admin_portal(inputs)
	return __ro.sidebar_admin_portal(inputs)
});
/**
* | output |
* | --- |
* | "Switch space" |
*
* @param {Sidebar_Switch_SpaceInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const sidebar_switch_space = /** @type {((inputs?: Sidebar_Switch_SpaceInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Sidebar_Switch_SpaceInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.sidebar_switch_space(inputs)
	return __ro.sidebar_switch_space(inputs)
});
/**
* | output |
* | --- |
* | "New schema" |
*
* @param {Schemas_New_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_new_title = /** @type {((inputs?: Schemas_New_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_New_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_new_title(inputs)
	return __ro.schemas_new_title(inputs)
});
/**
* | output |
* | --- |
* | "Library" |
*
* @param {Schemas_LibraryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_library = /** @type {((inputs?: Schemas_LibraryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_LibraryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_library(inputs)
	return __ro.schemas_library(inputs)
});
/**
* | output |
* | --- |
* | "Define schema metadata and structure." |
*
* @param {Schemas_New_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_new_description = /** @type {((inputs?: Schemas_New_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_New_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_new_description(inputs)
	return __ro.schemas_new_description(inputs)
});
/**
* | output |
* | --- |
* | "Edit schema" |
*
* @param {Schemas_Edit_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_edit_title = /** @type {((inputs?: Schemas_Edit_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Edit_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_edit_title(inputs)
	return __ro.schemas_edit_title(inputs)
});
/**
* | output |
* | --- |
* | "Update schema metadata and structure." |
*
* @param {Schemas_Edit_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_edit_description = /** @type {((inputs?: Schemas_Edit_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Edit_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_edit_description(inputs)
	return __ro.schemas_edit_description(inputs)
});
/**
* | output |
* | --- |
* | "Save schema" |
*
* @param {Schemas_Save_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_save_schema = /** @type {((inputs?: Schemas_Save_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Save_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_save_schema(inputs)
	return __ro.schemas_save_schema(inputs)
});
/**
* | output |
* | --- |
* | "Save changes" |
*
* @param {Schemas_Save_ChangesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_save_changes = /** @type {((inputs?: Schemas_Save_ChangesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Save_ChangesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_save_changes(inputs)
	return __ro.schemas_save_changes(inputs)
});
/**
* | output |
* | --- |
* | "Schema {name} saved successfully." |
*
* @param {Schemas_Saved_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_saved_success = /** @type {((inputs: Schemas_Saved_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Saved_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_saved_success(inputs)
	return __ro.schemas_saved_success(inputs)
});
/**
* | output |
* | --- |
* | "Schema {name} ({id}) saved successfully." |
*
* @param {Schemas_Saved_Success_With_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_saved_success_with_id = /** @type {((inputs: Schemas_Saved_Success_With_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Saved_Success_With_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_saved_success_with_id(inputs)
	return __ro.schemas_saved_success_with_id(inputs)
});
/**
* | output |
* | --- |
* | "Saved {name} ({id})" |
*
* @param {Schemas_Saved_FeedbackInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_saved_feedback = /** @type {((inputs: Schemas_Saved_FeedbackInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Saved_FeedbackInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_saved_feedback(inputs)
	return __ro.schemas_saved_feedback(inputs)
});
/**
* | output |
* | --- |
* | "Schema must include at least one field." |
*
* @param {Schemas_Empty_Schema_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_empty_schema_error = /** @type {((inputs?: Schemas_Empty_Schema_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Empty_Schema_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_empty_schema_error(inputs)
	return __ro.schemas_empty_schema_error(inputs)
});
/**
* | output |
* | --- |
* | "Delete schema?" |
*
* @param {Schemas_Delete_Single_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_single_title = /** @type {((inputs?: Schemas_Delete_Single_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Single_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_single_title(inputs)
	return __ro.schemas_delete_single_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete \"{name}\"? This action cannot be undone." |
*
* @param {Schemas_Delete_Single_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_single_description = /** @type {((inputs: Schemas_Delete_Single_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Single_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_single_description(inputs)
	return __ro.schemas_delete_single_description(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} schema?" |
*
* @param {Schemas_Delete_Bulk_Title_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_bulk_title_one = /** @type {((inputs: Schemas_Delete_Bulk_Title_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Bulk_Title_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_bulk_title_one(inputs)
	return __ro.schemas_delete_bulk_title_one(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} schemas?" |
*
* @param {Schemas_Delete_Bulk_Title_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_bulk_title_other = /** @type {((inputs: Schemas_Delete_Bulk_Title_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Bulk_Title_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_bulk_title_other(inputs)
	return __ro.schemas_delete_bulk_title_other(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} selected schema? This action cannot be undone." |
*
* @param {Schemas_Delete_Bulk_Description_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_bulk_description_one = /** @type {((inputs: Schemas_Delete_Bulk_Description_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Bulk_Description_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_bulk_description_one(inputs)
	return __ro.schemas_delete_bulk_description_one(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} selected schemas? This action cannot be undone." |
*
* @param {Schemas_Delete_Bulk_Description_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_bulk_description_other = /** @type {((inputs: Schemas_Delete_Bulk_Description_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_Bulk_Description_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_bulk_description_other(inputs)
	return __ro.schemas_delete_bulk_description_other(inputs)
});
/**
* | output |
* | --- |
* | "Select all schemas on this page" |
*
* @param {Schemas_Select_All_On_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_select_all_on_page = /** @type {((inputs?: Schemas_Select_All_On_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Select_All_On_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_select_all_on_page(inputs)
	return __ro.schemas_select_all_on_page(inputs)
});
/**
* | output |
* | --- |
* | "Select {name}" |
*
* @param {Schemas_Select_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_select_schema = /** @type {((inputs: Schemas_Select_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Select_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_select_schema(inputs)
	return __ro.schemas_select_schema(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Schemas_Name_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_name_column = /** @type {((inputs?: Schemas_Name_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Name_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_name_column(inputs)
	return __ro.schemas_name_column(inputs)
});
/**
* | output |
* | --- |
* | "ID" |
*
* @param {Schemas_Id_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_id_column = /** @type {((inputs?: Schemas_Id_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Id_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_id_column(inputs)
	return __ro.schemas_id_column(inputs)
});
/**
* | output |
* | --- |
* | "Schema ID" |
*
* @param {Schemas_Id_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_id_label = /** @type {((inputs?: Schemas_Id_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Id_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_id_label(inputs)
	return __ro.schemas_id_label(inputs)
});
/**
* | output |
* | --- |
* | "Copy ID" |
*
* @param {Schemas_Copy_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_copy_id = /** @type {((inputs?: Schemas_Copy_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Copy_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_copy_id(inputs)
	return __ro.schemas_copy_id(inputs)
});
/**
* | output |
* | --- |
* | "Copy schema ID {id}" |
*
* @param {Schemas_Copy_Id_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_copy_id_aria = /** @type {((inputs: Schemas_Copy_Id_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Copy_Id_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_copy_id_aria(inputs)
	return __ro.schemas_copy_id_aria(inputs)
});
/**
* | output |
* | --- |
* | "Schema ID copied." |
*
* @param {Schemas_Copy_Id_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_copy_id_success = /** @type {((inputs?: Schemas_Copy_Id_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Copy_Id_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_copy_id_success(inputs)
	return __ro.schemas_copy_id_success(inputs)
});
/**
* | output |
* | --- |
* | "Unable to copy schema ID." |
*
* @param {Schemas_Copy_Id_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_copy_id_error = /** @type {((inputs?: Schemas_Copy_Id_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Copy_Id_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_copy_id_error(inputs)
	return __ro.schemas_copy_id_error(inputs)
});
/**
* | output |
* | --- |
* | "Strict mode" |
*
* @param {Schemas_Strict_Mode_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_strict_mode_column = /** @type {((inputs?: Schemas_Strict_Mode_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Strict_Mode_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_strict_mode_column(inputs)
	return __ro.schemas_strict_mode_column(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Schemas_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_created_column = /** @type {((inputs?: Schemas_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_created_column(inputs)
	return __ro.schemas_created_column(inputs)
});
/**
* | output |
* | --- |
* | "Updated" |
*
* @param {Schemas_Updated_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_updated_column = /** @type {((inputs?: Schemas_Updated_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Updated_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_updated_column(inputs)
	return __ro.schemas_updated_column(inputs)
});
/**
* | output |
* | --- |
* | "New schema" |
*
* @param {Schemas_New_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_new_schema = /** @type {((inputs?: Schemas_New_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_New_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_new_schema(inputs)
	return __ro.schemas_new_schema(inputs)
});
/**
* | output |
* | --- |
* | "No schemas found" |
*
* @param {Schemas_No_Schemas_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_no_schemas_found = /** @type {((inputs?: Schemas_No_Schemas_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_No_Schemas_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_no_schemas_found(inputs)
	return __ro.schemas_no_schemas_found(inputs)
});
/**
* | output |
* | --- |
* | "Create a schema to define structured fields for document extraction." |
*
* @param {Schemas_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_empty_body = /** @type {((inputs?: Schemas_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_empty_body(inputs)
	return __ro.schemas_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Create schema" |
*
* @param {Schemas_Create_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_create_schema = /** @type {((inputs?: Schemas_Create_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Create_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_create_schema(inputs)
	return __ro.schemas_create_schema(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} schema on this page." |
*
* @param {Schemas_Showing_Schemas_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_showing_schemas_one = /** @type {((inputs: Schemas_Showing_Schemas_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Showing_Schemas_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_showing_schemas_one(inputs)
	return __ro.schemas_showing_schemas_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} schemas on this page." |
*
* @param {Schemas_Showing_Schemas_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_showing_schemas_other = /** @type {((inputs: Schemas_Showing_Schemas_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Showing_Schemas_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_showing_schemas_other(inputs)
	return __ro.schemas_showing_schemas_other(inputs)
});
/**
* | output |
* | --- |
* | "No schemas to show." |
*
* @param {Schemas_No_Schemas_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_no_schemas_to_show = /** @type {((inputs?: Schemas_No_Schemas_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_No_Schemas_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_no_schemas_to_show(inputs)
	return __ro.schemas_no_schemas_to_show(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Schemas_Selected_Count_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_selected_count_one = /** @type {((inputs: Schemas_Selected_Count_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Selected_Count_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_selected_count_one(inputs)
	return __ro.schemas_selected_count_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Schemas_Selected_Count_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_selected_count_other = /** @type {((inputs: Schemas_Selected_Count_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Selected_Count_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_selected_count_other(inputs)
	return __ro.schemas_selected_count_other(inputs)
});
/**
* | output |
* | --- |
* | "Deleting..." |
*
* @param {Schemas_DeletingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_deleting = /** @type {((inputs?: Schemas_DeletingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_DeletingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_deleting(inputs)
	return __ro.schemas_deleting(inputs)
});
/**
* | output |
* | --- |
* | "No description" |
*
* @param {Schemas_No_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_no_description = /** @type {((inputs?: Schemas_No_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_No_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_no_description(inputs)
	return __ro.schemas_no_description(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date ascending" |
*
* @param {Schemas_Sort_Created_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_sort_created_ascending = /** @type {((inputs?: Schemas_Sort_Created_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Sort_Created_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_sort_created_ascending(inputs)
	return __ro.schemas_sort_created_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date descending" |
*
* @param {Schemas_Sort_Created_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_sort_created_descending = /** @type {((inputs?: Schemas_Sort_Created_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Sort_Created_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_sort_created_descending(inputs)
	return __ro.schemas_sort_created_descending(inputs)
});
/**
* | output |
* | --- |
* | "Edit {name}" |
*
* @param {Schemas_Edit_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_edit_aria = /** @type {((inputs: Schemas_Edit_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Edit_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_edit_aria(inputs)
	return __ro.schemas_edit_aria(inputs)
});
/**
* | output |
* | --- |
* | "Create job with {name}" |
*
* @param {Schemas_Create_Job_WithInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_create_job_with = /** @type {((inputs: Schemas_Create_Job_WithInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Create_Job_WithInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_create_job_with(inputs)
	return __ro.schemas_create_job_with(inputs)
});
/**
* | output |
* | --- |
* | "Clone {name}" |
*
* @param {Schemas_Clone_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_clone_aria = /** @type {((inputs: Schemas_Clone_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Clone_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_clone_aria(inputs)
	return __ro.schemas_clone_aria(inputs)
});
/**
* | output |
* | --- |
* | "Delete {name}" |
*
* @param {Schemas_Delete_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_delete_aria = /** @type {((inputs: Schemas_Delete_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Delete_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_delete_aria(inputs)
	return __ro.schemas_delete_aria(inputs)
});
/**
* | output |
* | --- |
* | "Loading schema..." |
*
* @param {Schemas_Loading_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_loading_schema = /** @type {((inputs?: Schemas_Loading_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Loading_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_loading_schema(inputs)
	return __ro.schemas_loading_schema(inputs)
});
/**
* | output |
* | --- |
* | "Schema not found" |
*
* @param {Schemas_Not_Found_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_not_found_title = /** @type {((inputs?: Schemas_Not_Found_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Not_Found_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_not_found_title(inputs)
	return __ro.schemas_not_found_title(inputs)
});
/**
* | output |
* | --- |
* | "This schema does not exist." |
*
* @param {Schemas_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_not_found_body = /** @type {((inputs?: Schemas_Not_Found_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Not_Found_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_not_found_body(inputs)
	return __ro.schemas_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "View schemas" |
*
* @param {Schemas_View_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_view_schemas = /** @type {((inputs?: Schemas_View_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_View_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_view_schemas(inputs)
	return __ro.schemas_view_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Schema could not be loaded" |
*
* @param {Schemas_Could_Not_LoadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_could_not_load = /** @type {((inputs?: Schemas_Could_Not_LoadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Could_Not_LoadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_could_not_load(inputs)
	return __ro.schemas_could_not_load(inputs)
});
/**
* | output |
* | --- |
* | "Schema Editor" |
*
* @param {Schemas_Editor_BadgeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_editor_badge = /** @type {((inputs?: Schemas_Editor_BadgeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Editor_BadgeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_editor_badge(inputs)
	return __ro.schemas_editor_badge(inputs)
});
/**
* | output |
* | --- |
* | "General Settings" |
*
* @param {Schemas_General_SettingsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_general_settings = /** @type {((inputs?: Schemas_General_SettingsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_General_SettingsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_general_settings(inputs)
	return __ro.schemas_general_settings(inputs)
});
/**
* | output |
* | --- |
* | "Schema Name" |
*
* @param {Schemas_Schema_Name_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_schema_name_label = /** @type {((inputs?: Schemas_Schema_Name_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Schema_Name_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_schema_name_label(inputs)
	return __ro.schemas_schema_name_label(inputs)
});
/**
* | output |
* | --- |
* | "Schema name" |
*
* @param {Schemas_Schema_Name_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_schema_name_placeholder = /** @type {((inputs?: Schemas_Schema_Name_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Schema_Name_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_schema_name_placeholder(inputs)
	return __ro.schemas_schema_name_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Description" |
*
* @param {Schemas_Description_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_description_label = /** @type {((inputs?: Schemas_Description_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Description_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_description_label(inputs)
	return __ro.schemas_description_label(inputs)
});
/**
* | output |
* | --- |
* | "Provide optional context or instructions for this schema..." |
*
* @param {Schemas_Description_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_description_placeholder = /** @type {((inputs?: Schemas_Description_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Description_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_description_placeholder(inputs)
	return __ro.schemas_description_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Strict Mode" |
*
* @param {Schemas_Strict_ModeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_strict_mode = /** @type {((inputs?: Schemas_Strict_ModeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Strict_ModeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_strict_mode(inputs)
	return __ro.schemas_strict_mode(inputs)
});
/**
* | output |
* | --- |
* | "Flexible Mode" |
*
* @param {Schemas_Flexible_ModeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_flexible_mode = /** @type {((inputs?: Schemas_Flexible_ModeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Flexible_ModeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_flexible_mode(inputs)
	return __ro.schemas_flexible_mode(inputs)
});
/**
* | output |
* | --- |
* | "Reject fields not explicitly declared in this schema. Highly recommended for structured entity extraction." |
*
* @param {Schemas_Strict_Mode_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_strict_mode_description = /** @type {((inputs?: Schemas_Strict_Mode_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Strict_Mode_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_strict_mode_description(inputs)
	return __ro.schemas_strict_mode_description(inputs)
});
/**
* | output |
* | --- |
* | "Structure Designer" |
*
* @param {Schemas_Structure_DesignerInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_structure_designer = /** @type {((inputs?: Schemas_Structure_DesignerInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Structure_DesignerInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_structure_designer(inputs)
	return __ro.schemas_structure_designer(inputs)
});
/**
* | output |
* | --- |
* | "Visual Node Designer" |
*
* @param {Schemas_Visual_Node_DesignerInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_visual_node_designer = /** @type {((inputs?: Schemas_Visual_Node_DesignerInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Visual_Node_DesignerInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_visual_node_designer(inputs)
	return __ro.schemas_visual_node_designer(inputs)
});
/**
* | output |
* | --- |
* | "Name is required." |
*
* @param {Schemas_Validation_Name_RequiredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_validation_name_required = /** @type {((inputs?: Schemas_Validation_Name_RequiredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Validation_Name_RequiredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_validation_name_required(inputs)
	return __ro.schemas_validation_name_required(inputs)
});
/**
* | output |
* | --- |
* | "Name must be at most 160 characters." |
*
* @param {Schemas_Validation_Name_Too_LongInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_validation_name_too_long = /** @type {((inputs?: Schemas_Validation_Name_Too_LongInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Validation_Name_Too_LongInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_validation_name_too_long(inputs)
	return __ro.schemas_validation_name_too_long(inputs)
});
/**
* | output |
* | --- |
* | "Schema must be a JSON object." |
*
* @param {Schemas_Validation_Schema_ObjectInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_validation_schema_object = /** @type {((inputs?: Schemas_Validation_Schema_ObjectInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_Validation_Schema_ObjectInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_validation_schema_object(inputs)
	return __ro.schemas_validation_schema_object(inputs)
});
/**
* | output |
* | --- |
* | "Clone" |
*
* @param {Schemas_CloneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_clone = /** @type {((inputs?: Schemas_CloneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_CloneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_clone(inputs)
	return __ro.schemas_clone(inputs)
});
/**
* | output |
* | --- |
* | "Cloning..." |
*
* @param {Schemas_CloningInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_cloning = /** @type {((inputs?: Schemas_CloningInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_CloningInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_cloning(inputs)
	return __ro.schemas_cloning(inputs)
});
/**
* | output |
* | --- |
* | "Saving..." |
*
* @param {Schemas_SavingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const schemas_saving = /** @type {((inputs?: Schemas_SavingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Schemas_SavingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.schemas_saving(inputs)
	return __ro.schemas_saving(inputs)
});
/**
* | output |
* | --- |
* | "JSON Recipes" |
*
* @param {Json_Recipes_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_title = /** @type {((inputs?: Json_Recipes_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_title(inputs)
	return __ro.json_recipes_title(inputs)
});
/**
* | output |
* | --- |
* | "Admin-managed templates for deploying extraction schemas." |
*
* @param {Json_Recipes_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_description = /** @type {((inputs?: Json_Recipes_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_description(inputs)
	return __ro.json_recipes_description(inputs)
});
/**
* | output |
* | --- |
* | "New recipe" |
*
* @param {Json_Recipes_New_RecipeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_new_recipe = /** @type {((inputs?: Json_Recipes_New_RecipeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_New_RecipeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_new_recipe(inputs)
	return __ro.json_recipes_new_recipe(inputs)
});
/**
* | output |
* | --- |
* | "No recipes found" |
*
* @param {Json_Recipes_No_Recipes_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_no_recipes_found = /** @type {((inputs?: Json_Recipes_No_Recipes_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_No_Recipes_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_no_recipes_found(inputs)
	return __ro.json_recipes_no_recipes_found(inputs)
});
/**
* | output |
* | --- |
* | "Create a recipe to make a reusable extraction schema available for deployment." |
*
* @param {Json_Recipes_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_empty_body = /** @type {((inputs?: Json_Recipes_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_empty_body(inputs)
	return __ro.json_recipes_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Loading recipes..." |
*
* @param {Json_Recipes_LoadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_loading = /** @type {((inputs?: Json_Recipes_LoadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_LoadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_loading(inputs)
	return __ro.json_recipes_loading(inputs)
});
/**
* | output |
* | --- |
* | "Loading recipe..." |
*
* @param {Json_Recipes_Loading_RecipeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_loading_recipe = /** @type {((inputs?: Json_Recipes_Loading_RecipeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Loading_RecipeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_loading_recipe(inputs)
	return __ro.json_recipes_loading_recipe(inputs)
});
/**
* | output |
* | --- |
* | "Deploys" |
*
* @param {Json_Recipes_Counter_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_counter_column = /** @type {((inputs?: Json_Recipes_Counter_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Counter_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_counter_column(inputs)
	return __ro.json_recipes_counter_column(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Json_Recipes_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_created_column = /** @type {((inputs?: Json_Recipes_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_created_column(inputs)
	return __ro.json_recipes_created_column(inputs)
});
/**
* | output |
* | --- |
* | "Updated" |
*
* @param {Json_Recipes_Updated_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_updated_column = /** @type {((inputs?: Json_Recipes_Updated_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Updated_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_updated_column(inputs)
	return __ro.json_recipes_updated_column(inputs)
});
/**
* | output |
* | --- |
* | "Fields" |
*
* @param {Json_Recipes_Json_Fields_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_json_fields_column = /** @type {((inputs?: Json_Recipes_Json_Fields_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Json_Fields_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_json_fields_column(inputs)
	return __ro.json_recipes_json_fields_column(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date ascending" |
*
* @param {Json_Recipes_Sort_Created_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_sort_created_ascending = /** @type {((inputs?: Json_Recipes_Sort_Created_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Sort_Created_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_sort_created_ascending(inputs)
	return __ro.json_recipes_sort_created_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date descending" |
*
* @param {Json_Recipes_Sort_Created_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_sort_created_descending = /** @type {((inputs?: Json_Recipes_Sort_Created_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Sort_Created_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_sort_created_descending(inputs)
	return __ro.json_recipes_sort_created_descending(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} recipe on this page." |
*
* @param {Json_Recipes_Showing_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_showing_one = /** @type {((inputs: Json_Recipes_Showing_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Showing_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_showing_one(inputs)
	return __ro.json_recipes_showing_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} recipes on this page." |
*
* @param {Json_Recipes_Showing_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_showing_other = /** @type {((inputs: Json_Recipes_Showing_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Showing_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_showing_other(inputs)
	return __ro.json_recipes_showing_other(inputs)
});
/**
* | output |
* | --- |
* | "No recipes to show." |
*
* @param {Json_Recipes_No_Recipes_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_no_recipes_to_show = /** @type {((inputs?: Json_Recipes_No_Recipes_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_No_Recipes_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_no_recipes_to_show(inputs)
	return __ro.json_recipes_no_recipes_to_show(inputs)
});
/**
* | output |
* | --- |
* | "Edit {name}" |
*
* @param {Json_Recipes_Edit_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_edit_aria = /** @type {((inputs: Json_Recipes_Edit_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Edit_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_edit_aria(inputs)
	return __ro.json_recipes_edit_aria(inputs)
});
/**
* | output |
* | --- |
* | "Delete {name}" |
*
* @param {Json_Recipes_Delete_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_delete_aria = /** @type {((inputs: Json_Recipes_Delete_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Delete_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_delete_aria(inputs)
	return __ro.json_recipes_delete_aria(inputs)
});
/**
* | output |
* | --- |
* | "New JSON recipe" |
*
* @param {Json_Recipes_New_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_new_title = /** @type {((inputs?: Json_Recipes_New_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_New_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_new_title(inputs)
	return __ro.json_recipes_new_title(inputs)
});
/**
* | output |
* | --- |
* | "Define recipe metadata and JSON Schema structure." |
*
* @param {Json_Recipes_New_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_new_description = /** @type {((inputs?: Json_Recipes_New_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_New_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_new_description(inputs)
	return __ro.json_recipes_new_description(inputs)
});
/**
* | output |
* | --- |
* | "Edit JSON recipe" |
*
* @param {Json_Recipes_Edit_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_edit_title = /** @type {((inputs?: Json_Recipes_Edit_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Edit_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_edit_title(inputs)
	return __ro.json_recipes_edit_title(inputs)
});
/**
* | output |
* | --- |
* | "Update recipe metadata and JSON Schema structure." |
*
* @param {Json_Recipes_Edit_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_edit_description = /** @type {((inputs?: Json_Recipes_Edit_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Edit_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_edit_description(inputs)
	return __ro.json_recipes_edit_description(inputs)
});
/**
* | output |
* | --- |
* | "Save recipe" |
*
* @param {Json_Recipes_Save_RecipeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_save_recipe = /** @type {((inputs?: Json_Recipes_Save_RecipeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Save_RecipeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_save_recipe(inputs)
	return __ro.json_recipes_save_recipe(inputs)
});
/**
* | output |
* | --- |
* | "Save changes" |
*
* @param {Json_Recipes_Save_ChangesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_save_changes = /** @type {((inputs?: Json_Recipes_Save_ChangesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Save_ChangesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_save_changes(inputs)
	return __ro.json_recipes_save_changes(inputs)
});
/**
* | output |
* | --- |
* | "Recipe {name} created." |
*
* @param {Json_Recipes_Created_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_created_success = /** @type {((inputs: Json_Recipes_Created_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Created_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_created_success(inputs)
	return __ro.json_recipes_created_success(inputs)
});
/**
* | output |
* | --- |
* | "Recipe {name} saved." |
*
* @param {Json_Recipes_Saved_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_saved_success = /** @type {((inputs: Json_Recipes_Saved_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Saved_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_saved_success(inputs)
	return __ro.json_recipes_saved_success(inputs)
});
/**
* | output |
* | --- |
* | "Recipe {name} deleted." |
*
* @param {Json_Recipes_Deleted_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_deleted_success = /** @type {((inputs: Json_Recipes_Deleted_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Deleted_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_deleted_success(inputs)
	return __ro.json_recipes_deleted_success(inputs)
});
/**
* | output |
* | --- |
* | "Delete this recipe? Deployed schemas will remain unchanged." |
*
* @param {Json_Recipes_Delete_ConfirmInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_delete_confirm = /** @type {((inputs?: Json_Recipes_Delete_ConfirmInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Delete_ConfirmInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_delete_confirm(inputs)
	return __ro.json_recipes_delete_confirm(inputs)
});
/**
* | output |
* | --- |
* | "Recipe not found" |
*
* @param {Json_Recipes_Not_Found_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_not_found_title = /** @type {((inputs?: Json_Recipes_Not_Found_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Not_Found_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_not_found_title(inputs)
	return __ro.json_recipes_not_found_title(inputs)
});
/**
* | output |
* | --- |
* | "This JSON recipe does not exist." |
*
* @param {Json_Recipes_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_not_found_body = /** @type {((inputs?: Json_Recipes_Not_Found_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Not_Found_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_not_found_body(inputs)
	return __ro.json_recipes_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "View recipes" |
*
* @param {Json_Recipes_View_RecipesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_view_recipes = /** @type {((inputs?: Json_Recipes_View_RecipesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_View_RecipesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_view_recipes(inputs)
	return __ro.json_recipes_view_recipes(inputs)
});
/**
* | output |
* | --- |
* | "Recipe could not be loaded" |
*
* @param {Json_Recipes_Could_Not_LoadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_could_not_load = /** @type {((inputs?: Json_Recipes_Could_Not_LoadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Could_Not_LoadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_could_not_load(inputs)
	return __ro.json_recipes_could_not_load(inputs)
});
/**
* | output |
* | --- |
* | "Recipe Editor" |
*
* @param {Json_Recipes_Editor_BadgeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_editor_badge = /** @type {((inputs?: Json_Recipes_Editor_BadgeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Editor_BadgeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_editor_badge(inputs)
	return __ro.json_recipes_editor_badge(inputs)
});
/**
* | output |
* | --- |
* | "Recipe Settings" |
*
* @param {Json_Recipes_General_SettingsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_general_settings = /** @type {((inputs?: Json_Recipes_General_SettingsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_General_SettingsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_general_settings(inputs)
	return __ro.json_recipes_general_settings(inputs)
});
/**
* | output |
* | --- |
* | "Title" |
*
* @param {Json_Recipes_Title_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_title_label = /** @type {((inputs?: Json_Recipes_Title_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Title_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_title_label(inputs)
	return __ro.json_recipes_title_label(inputs)
});
/**
* | output |
* | --- |
* | "Recipe title" |
*
* @param {Json_Recipes_Title_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_title_placeholder = /** @type {((inputs?: Json_Recipes_Title_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Title_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_title_placeholder(inputs)
	return __ro.json_recipes_title_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Description" |
*
* @param {Json_Recipes_Description_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_description_label = /** @type {((inputs?: Json_Recipes_Description_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Description_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_description_label(inputs)
	return __ro.json_recipes_description_label(inputs)
});
/**
* | output |
* | --- |
* | "Describe when this recipe should be used..." |
*
* @param {Json_Recipes_Description_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_description_placeholder = /** @type {((inputs?: Json_Recipes_Description_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Description_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_description_placeholder(inputs)
	return __ro.json_recipes_description_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Structure Designer" |
*
* @param {Json_Recipes_Structure_DesignerInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_structure_designer = /** @type {((inputs?: Json_Recipes_Structure_DesignerInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Structure_DesignerInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_structure_designer(inputs)
	return __ro.json_recipes_structure_designer(inputs)
});
/**
* | output |
* | --- |
* | "Visual Node Designer" |
*
* @param {Json_Recipes_Visual_Node_DesignerInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_visual_node_designer = /** @type {((inputs?: Json_Recipes_Visual_Node_DesignerInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Visual_Node_DesignerInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_visual_node_designer(inputs)
	return __ro.json_recipes_visual_node_designer(inputs)
});
/**
* | output |
* | --- |
* | "Category" |
*
* @param {Json_Recipes_Category_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_category_label = /** @type {((inputs?: Json_Recipes_Category_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Category_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_category_label(inputs)
	return __ro.json_recipes_category_label(inputs)
});
/**
* | output |
* | --- |
* | "Others" |
*
* @param {Json_Recipes_OthersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_others = /** @type {((inputs?: Json_Recipes_OthersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_OthersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_others(inputs)
	return __ro.json_recipes_others(inputs)
});
/**
* | output |
* | --- |
* | "Manage categories" |
*
* @param {Json_Recipes_Manage_CategoriesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_manage_categories = /** @type {((inputs?: Json_Recipes_Manage_CategoriesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Manage_CategoriesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_manage_categories(inputs)
	return __ro.json_recipes_manage_categories(inputs)
});
/**
* | output |
* | --- |
* | "Title is required." |
*
* @param {Json_Recipes_Validation_Title_RequiredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_validation_title_required = /** @type {((inputs?: Json_Recipes_Validation_Title_RequiredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Validation_Title_RequiredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_validation_title_required(inputs)
	return __ro.json_recipes_validation_title_required(inputs)
});
/**
* | output |
* | --- |
* | "Title must be at most 160 characters." |
*
* @param {Json_Recipes_Validation_Title_Too_LongInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_validation_title_too_long = /** @type {((inputs?: Json_Recipes_Validation_Title_Too_LongInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Validation_Title_Too_LongInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_validation_title_too_long(inputs)
	return __ro.json_recipes_validation_title_too_long(inputs)
});
/**
* | output |
* | --- |
* | "Recipe JSON must be a JSON object." |
*
* @param {Json_Recipes_Validation_Json_ObjectInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_validation_json_object = /** @type {((inputs?: Json_Recipes_Validation_Json_ObjectInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_Validation_Json_ObjectInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_validation_json_object(inputs)
	return __ro.json_recipes_validation_json_object(inputs)
});
/**
* | output |
* | --- |
* | "Saving..." |
*
* @param {Json_Recipes_SavingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_saving = /** @type {((inputs?: Json_Recipes_SavingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_SavingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_saving(inputs)
	return __ro.json_recipes_saving(inputs)
});
/**
* | output |
* | --- |
* | "Deleting..." |
*
* @param {Json_Recipes_DeletingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipes_deleting = /** @type {((inputs?: Json_Recipes_DeletingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipes_DeletingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipes_deleting(inputs)
	return __ro.json_recipes_deleting(inputs)
});
/**
* | output |
* | --- |
* | "JSON Recipe Categories" |
*
* @param {Json_Recipe_Categories_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_title = /** @type {((inputs?: Json_Recipe_Categories_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_title(inputs)
	return __ro.json_recipe_categories_title(inputs)
});
/**
* | output |
* | --- |
* | "Define localized category labels for grouping JSON recipes." |
*
* @param {Json_Recipe_Categories_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_description = /** @type {((inputs?: Json_Recipe_Categories_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_description(inputs)
	return __ro.json_recipe_categories_description(inputs)
});
/**
* | output |
* | --- |
* | "English title" |
*
* @param {Json_Recipe_Categories_Title_En_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_title_en_label = /** @type {((inputs?: Json_Recipe_Categories_Title_En_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Title_En_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_title_en_label(inputs)
	return __ro.json_recipe_categories_title_en_label(inputs)
});
/**
* | output |
* | --- |
* | "Romanian title" |
*
* @param {Json_Recipe_Categories_Title_Ro_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_title_ro_label = /** @type {((inputs?: Json_Recipe_Categories_Title_Ro_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Title_Ro_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_title_ro_label(inputs)
	return __ro.json_recipe_categories_title_ro_label(inputs)
});
/**
* | output |
* | --- |
* | "Create category" |
*
* @param {Json_Recipe_Categories_Create_CategoryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_create_category = /** @type {((inputs?: Json_Recipe_Categories_Create_CategoryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Create_CategoryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_create_category(inputs)
	return __ro.json_recipe_categories_create_category(inputs)
});
/**
* | output |
* | --- |
* | "Save category" |
*
* @param {Json_Recipe_Categories_Save_CategoryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_save_category = /** @type {((inputs?: Json_Recipe_Categories_Save_CategoryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Save_CategoryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_save_category(inputs)
	return __ro.json_recipe_categories_save_category(inputs)
});
/**
* | output |
* | --- |
* | "Edit category" |
*
* @param {Json_Recipe_Categories_Edit_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_edit_title = /** @type {((inputs?: Json_Recipe_Categories_Edit_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Edit_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_edit_title(inputs)
	return __ro.json_recipe_categories_edit_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete this category? Recipes must be moved out before deletion." |
*
* @param {Json_Recipe_Categories_Delete_ConfirmInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_delete_confirm = /** @type {((inputs?: Json_Recipe_Categories_Delete_ConfirmInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Delete_ConfirmInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_delete_confirm(inputs)
	return __ro.json_recipe_categories_delete_confirm(inputs)
});
/**
* | output |
* | --- |
* | "Loading categories..." |
*
* @param {Json_Recipe_Categories_LoadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_loading = /** @type {((inputs?: Json_Recipe_Categories_LoadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_LoadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_loading(inputs)
	return __ro.json_recipe_categories_loading(inputs)
});
/**
* | output |
* | --- |
* | "Categories could not be loaded" |
*
* @param {Json_Recipe_Categories_Could_Not_LoadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_could_not_load = /** @type {((inputs?: Json_Recipe_Categories_Could_Not_LoadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Could_Not_LoadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_could_not_load(inputs)
	return __ro.json_recipe_categories_could_not_load(inputs)
});
/**
* | output |
* | --- |
* | "No categories yet" |
*
* @param {Json_Recipe_Categories_Empty_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_empty_title = /** @type {((inputs?: Json_Recipe_Categories_Empty_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Empty_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_empty_title(inputs)
	return __ro.json_recipe_categories_empty_title(inputs)
});
/**
* | output |
* | --- |
* | "Recipes without a category will appear under Others." |
*
* @param {Json_Recipe_Categories_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_empty_body = /** @type {((inputs?: Json_Recipe_Categories_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_empty_body(inputs)
	return __ro.json_recipe_categories_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Category {name} created." |
*
* @param {Json_Recipe_Categories_Created_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_created_success = /** @type {((inputs: Json_Recipe_Categories_Created_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Created_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_created_success(inputs)
	return __ro.json_recipe_categories_created_success(inputs)
});
/**
* | output |
* | --- |
* | "Category {name} saved." |
*
* @param {Json_Recipe_Categories_Saved_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_saved_success = /** @type {((inputs: Json_Recipe_Categories_Saved_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Saved_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_saved_success(inputs)
	return __ro.json_recipe_categories_saved_success(inputs)
});
/**
* | output |
* | --- |
* | "Category {name} deleted." |
*
* @param {Json_Recipe_Categories_Deleted_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_deleted_success = /** @type {((inputs: Json_Recipe_Categories_Deleted_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Deleted_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_deleted_success(inputs)
	return __ro.json_recipe_categories_deleted_success(inputs)
});
/**
* | output |
* | --- |
* | "Both English and Romanian titles are required." |
*
* @param {Json_Recipe_Categories_Validation_Titles_RequiredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_validation_titles_required = /** @type {((inputs?: Json_Recipe_Categories_Validation_Titles_RequiredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Validation_Titles_RequiredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_validation_titles_required(inputs)
	return __ro.json_recipe_categories_validation_titles_required(inputs)
});
/**
* | output |
* | --- |
* | "Titles must be at most 160 characters." |
*
* @param {Json_Recipe_Categories_Validation_Titles_Too_LongInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_validation_titles_too_long = /** @type {((inputs?: Json_Recipe_Categories_Validation_Titles_Too_LongInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Validation_Titles_Too_LongInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_validation_titles_too_long(inputs)
	return __ro.json_recipe_categories_validation_titles_too_long(inputs)
});
/**
* | output |
* | --- |
* | "Edit {name}" |
*
* @param {Json_Recipe_Categories_Edit_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_edit_aria = /** @type {((inputs: Json_Recipe_Categories_Edit_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Edit_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_edit_aria(inputs)
	return __ro.json_recipe_categories_edit_aria(inputs)
});
/**
* | output |
* | --- |
* | "Delete {name}" |
*
* @param {Json_Recipe_Categories_Delete_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const json_recipe_categories_delete_aria = /** @type {((inputs: Json_Recipe_Categories_Delete_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Json_Recipe_Categories_Delete_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.json_recipe_categories_delete_aria(inputs)
	return __ro.json_recipe_categories_delete_aria(inputs)
});
/**
* | output |
* | --- |
* | "OCR Recipes" |
*
* @param {Ocr_Recipes_NavInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_nav = /** @type {((inputs?: Ocr_Recipes_NavInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_NavInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_nav(inputs)
	return __ro.ocr_recipes_nav(inputs)
});
/**
* | output |
* | --- |
* | "OCR Recipes" |
*
* @param {Ocr_Recipes_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_title = /** @type {((inputs?: Ocr_Recipes_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_title(inputs)
	return __ro.ocr_recipes_title(inputs)
});
/**
* | output |
* | --- |
* | "Browse system OCR JSON recipes and clone them into your Syncra schemas." |
*
* @param {Ocr_Recipes_Meta_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_meta_description = /** @type {((inputs?: Ocr_Recipes_Meta_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Meta_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_meta_description(inputs)
	return __ro.ocr_recipes_meta_description(inputs)
});
/**
* | output |
* | --- |
* | "System extraction templates" |
*
* @param {Ocr_Recipes_EyebrowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_eyebrow = /** @type {((inputs?: Ocr_Recipes_EyebrowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_EyebrowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_eyebrow(inputs)
	return __ro.ocr_recipes_eyebrow(inputs)
});
/**
* | output |
* | --- |
* | "Start from a proven OCR recipe" |
*
* @param {Ocr_Recipes_Hero_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_hero_title = /** @type {((inputs?: Ocr_Recipes_Hero_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Hero_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_hero_title(inputs)
	return __ro.ocr_recipes_hero_title(inputs)
});
/**
* | output |
* | --- |
* | "Browse reusable JSON Schema recipes for common Romanian documents. Clone a recipe into your workspace, then tune it in the schema editor before running OCR j..." |
*
* @param {Ocr_Recipes_Hero_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_hero_description = /** @type {((inputs?: Ocr_Recipes_Hero_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Hero_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_hero_description(inputs)
	return __ro.ocr_recipes_hero_description(inputs)
});
/**
* | output |
* | --- |
* | "Search recipes" |
*
* @param {Ocr_Recipes_Search_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_search_label = /** @type {((inputs?: Ocr_Recipes_Search_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Search_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_search_label(inputs)
	return __ro.ocr_recipes_search_label(inputs)
});
/**
* | output |
* | --- |
* | "Search by recipe, field, type, or description" |
*
* @param {Ocr_Recipes_Search_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_search_placeholder = /** @type {((inputs?: Ocr_Recipes_Search_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Search_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_search_placeholder(inputs)
	return __ro.ocr_recipes_search_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Filter by category" |
*
* @param {Ocr_Recipes_Category_FilterInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_category_filter = /** @type {((inputs?: Ocr_Recipes_Category_FilterInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Category_FilterInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_category_filter(inputs)
	return __ro.ocr_recipes_category_filter(inputs)
});
/**
* | output |
* | --- |
* | "All categories" |
*
* @param {Ocr_Recipes_All_CategoriesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_all_categories = /** @type {((inputs?: Ocr_Recipes_All_CategoriesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_All_CategoriesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_all_categories(inputs)
	return __ro.ocr_recipes_all_categories(inputs)
});
/**
* | output |
* | --- |
* | "Sort recipes" |
*
* @param {Ocr_Recipes_Sort_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_sort_label = /** @type {((inputs?: Ocr_Recipes_Sort_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Sort_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_sort_label(inputs)
	return __ro.ocr_recipes_sort_label(inputs)
});
/**
* | output |
* | --- |
* | "Most cloned" |
*
* @param {Ocr_Recipes_Sort_PopularInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_sort_popular = /** @type {((inputs?: Ocr_Recipes_Sort_PopularInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Sort_PopularInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_sort_popular(inputs)
	return __ro.ocr_recipes_sort_popular(inputs)
});
/**
* | output |
* | --- |
* | "Newest" |
*
* @param {Ocr_Recipes_Sort_NewestInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_sort_newest = /** @type {((inputs?: Ocr_Recipes_Sort_NewestInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Sort_NewestInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_sort_newest(inputs)
	return __ro.ocr_recipes_sort_newest(inputs)
});
/**
* | output |
* | --- |
* | "A-Z" |
*
* @param {Ocr_Recipes_Sort_AzInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_sort_az = /** @type {((inputs?: Ocr_Recipes_Sort_AzInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Sort_AzInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_sort_az(inputs)
	return __ro.ocr_recipes_sort_az(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} recipe" |
*
* @param {Ocr_Recipes_Showing_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_showing_one = /** @type {((inputs: Ocr_Recipes_Showing_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Showing_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_showing_one(inputs)
	return __ro.ocr_recipes_showing_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} recipes" |
*
* @param {Ocr_Recipes_Showing_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_showing_other = /** @type {((inputs: Ocr_Recipes_Showing_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Showing_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_showing_other(inputs)
	return __ro.ocr_recipes_showing_other(inputs)
});
/**
* | output |
* | --- |
* | "No recipes match your search" |
*
* @param {Ocr_Recipes_No_Matches_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_no_matches_title = /** @type {((inputs?: Ocr_Recipes_No_Matches_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_No_Matches_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_no_matches_title(inputs)
	return __ro.ocr_recipes_no_matches_title(inputs)
});
/**
* | output |
* | --- |
* | "Clear the search field to browse every system recipe." |
*
* @param {Ocr_Recipes_No_Matches_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_no_matches_body = /** @type {((inputs?: Ocr_Recipes_No_Matches_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_No_Matches_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_no_matches_body(inputs)
	return __ro.ocr_recipes_no_matches_body(inputs)
});
/**
* | output |
* | --- |
* | "Others" |
*
* @param {Ocr_Recipes_OthersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_others = /** @type {((inputs?: Ocr_Recipes_OthersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_OthersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_others(inputs)
	return __ro.ocr_recipes_others(inputs)
});
/**
* | output |
* | --- |
* | "{count} field" |
*
* @param {Ocr_Recipes_Fields_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_fields_one = /** @type {((inputs: Ocr_Recipes_Fields_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Fields_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_fields_one(inputs)
	return __ro.ocr_recipes_fields_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} fields" |
*
* @param {Ocr_Recipes_Fields_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_fields_other = /** @type {((inputs: Ocr_Recipes_Fields_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Fields_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_fields_other(inputs)
	return __ro.ocr_recipes_fields_other(inputs)
});
/**
* | output |
* | --- |
* | "{count} required" |
*
* @param {Ocr_Recipes_Required_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_required_one = /** @type {((inputs: Ocr_Recipes_Required_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Required_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_required_one(inputs)
	return __ro.ocr_recipes_required_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} required" |
*
* @param {Ocr_Recipes_Required_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_required_other = /** @type {((inputs: Ocr_Recipes_Required_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Required_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_required_other(inputs)
	return __ro.ocr_recipes_required_other(inputs)
});
/**
* | output |
* | --- |
* | "{count} clone" |
*
* @param {Ocr_Recipes_Deploys_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_deploys_one = /** @type {((inputs: Ocr_Recipes_Deploys_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Deploys_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_deploys_one(inputs)
	return __ro.ocr_recipes_deploys_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} clones" |
*
* @param {Ocr_Recipes_Deploys_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_deploys_other = /** @type {((inputs: Ocr_Recipes_Deploys_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Deploys_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_deploys_other(inputs)
	return __ro.ocr_recipes_deploys_other(inputs)
});
/**
* | output |
* | --- |
* | "JSON fields" |
*
* @param {Ocr_Recipes_Json_FieldsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_json_fields = /** @type {((inputs?: Ocr_Recipes_Json_FieldsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Json_FieldsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_json_fields(inputs)
	return __ro.ocr_recipes_json_fields(inputs)
});
/**
* | output |
* | --- |
* | "System recipe" |
*
* @param {Ocr_Recipes_System_RecipeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_system_recipe = /** @type {((inputs?: Ocr_Recipes_System_RecipeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_System_RecipeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_system_recipe(inputs)
	return __ro.ocr_recipes_system_recipe(inputs)
});
/**
* | output |
* | --- |
* | "Strict JSON Schema" |
*
* @param {Ocr_Recipes_Strict_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_strict_schema = /** @type {((inputs?: Ocr_Recipes_Strict_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Strict_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_strict_schema(inputs)
	return __ro.ocr_recipes_strict_schema(inputs)
});
/**
* | output |
* | --- |
* | "Required" |
*
* @param {Ocr_Recipes_RequiredInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_required = /** @type {((inputs?: Ocr_Recipes_RequiredInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_RequiredInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_required(inputs)
	return __ro.ocr_recipes_required(inputs)
});
/**
* | output |
* | --- |
* | "Preview JSON" |
*
* @param {Ocr_Recipes_Preview_JsonInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_preview_json = /** @type {((inputs?: Ocr_Recipes_Preview_JsonInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Preview_JsonInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_preview_json(inputs)
	return __ro.ocr_recipes_preview_json(inputs)
});
/**
* | output |
* | --- |
* | "No JSON fields defined." |
*
* @param {Ocr_Recipes_No_FieldsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_no_fields = /** @type {((inputs?: Ocr_Recipes_No_FieldsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_No_FieldsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_no_fields(inputs)
	return __ro.ocr_recipes_no_fields(inputs)
});
/**
* | output |
* | --- |
* | "Clone recipe" |
*
* @param {Ocr_Recipes_Clone_RecipeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_clone_recipe = /** @type {((inputs?: Ocr_Recipes_Clone_RecipeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Clone_RecipeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_clone_recipe(inputs)
	return __ro.ocr_recipes_clone_recipe(inputs)
});
/**
* | output |
* | --- |
* | "Clone {name}" |
*
* @param {Ocr_Recipes_Clone_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_clone_aria = /** @type {((inputs: Ocr_Recipes_Clone_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Clone_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_clone_aria(inputs)
	return __ro.ocr_recipes_clone_aria(inputs)
});
/**
* | output |
* | --- |
* | "Log in to clone" |
*
* @param {Ocr_Recipes_Log_In_To_CloneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_log_in_to_clone = /** @type {((inputs?: Ocr_Recipes_Log_In_To_CloneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Log_In_To_CloneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_log_in_to_clone(inputs)
	return __ro.ocr_recipes_log_in_to_clone(inputs)
});
/**
* | output |
* | --- |
* | "Unable to clone recipe." |
*
* @param {Ocr_Recipes_Clone_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_clone_failed = /** @type {((inputs?: Ocr_Recipes_Clone_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Clone_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_clone_failed(inputs)
	return __ro.ocr_recipes_clone_failed(inputs)
});
/**
* | output |
* | --- |
* | "Unable to load OCR recipes." |
*
* @param {Ocr_Recipes_Load_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const ocr_recipes_load_failed = /** @type {((inputs?: Ocr_Recipes_Load_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Ocr_Recipes_Load_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.ocr_recipes_load_failed(inputs)
	return __ro.ocr_recipes_load_failed(inputs)
});
/**
* | output |
* | --- |
* | "Jobs" |
*
* @param {Jobs_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_page_title = /** @type {((inputs?: Jobs_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_page_title(inputs)
	return __ro.jobs_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Missing schema id" |
*
* @param {Jobs_Missing_Schema_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_missing_schema_id = /** @type {((inputs?: Jobs_Missing_Schema_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Missing_Schema_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_missing_schema_id(inputs)
	return __ro.jobs_missing_schema_id(inputs)
});
/**
* | output |
* | --- |
* | "Missing job id" |
*
* @param {Jobs_Missing_Job_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_missing_job_id = /** @type {((inputs?: Jobs_Missing_Job_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Missing_Job_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_missing_job_id(inputs)
	return __ro.jobs_missing_job_id(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} job?" |
*
* @param {Jobs_Delete_Bulk_Title_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_bulk_title_one = /** @type {((inputs: Jobs_Delete_Bulk_Title_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Bulk_Title_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_bulk_title_one(inputs)
	return __ro.jobs_delete_bulk_title_one(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} jobs?" |
*
* @param {Jobs_Delete_Bulk_Title_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_bulk_title_other = /** @type {((inputs: Jobs_Delete_Bulk_Title_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Bulk_Title_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_bulk_title_other(inputs)
	return __ro.jobs_delete_bulk_title_other(inputs)
});
/**
* | output |
* | --- |
* | "Remove {count} selected job from the jobs list. Generated documents remain available." |
*
* @param {Jobs_Delete_Bulk_Description_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_bulk_description_one = /** @type {((inputs: Jobs_Delete_Bulk_Description_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Bulk_Description_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_bulk_description_one(inputs)
	return __ro.jobs_delete_bulk_description_one(inputs)
});
/**
* | output |
* | --- |
* | "Remove {count} selected jobs from the jobs list. Generated documents remain available." |
*
* @param {Jobs_Delete_Bulk_Description_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_bulk_description_other = /** @type {((inputs: Jobs_Delete_Bulk_Description_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Bulk_Description_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_bulk_description_other(inputs)
	return __ro.jobs_delete_bulk_description_other(inputs)
});
/**
* | output |
* | --- |
* | "Delete job?" |
*
* @param {Jobs_Delete_Single_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_single_title = /** @type {((inputs?: Jobs_Delete_Single_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Single_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_single_title(inputs)
	return __ro.jobs_delete_single_title(inputs)
});
/**
* | output |
* | --- |
* | "Remove \"{name}\" from the jobs list. Generated documents remain available." |
*
* @param {Jobs_Delete_Single_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_single_description = /** @type {((inputs: Jobs_Delete_Single_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_Single_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_single_description(inputs)
	return __ro.jobs_delete_single_description(inputs)
});
/**
* | output |
* | --- |
* | "Queued" |
*
* @param {Jobs_Status_QueuedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_queued = /** @type {((inputs?: Jobs_Status_QueuedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_QueuedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_queued(inputs)
	return __ro.jobs_status_queued(inputs)
});
/**
* | output |
* | --- |
* | "Pending" |
*
* @param {Jobs_Status_PendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_pending = /** @type {((inputs?: Jobs_Status_PendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_PendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_pending(inputs)
	return __ro.jobs_status_pending(inputs)
});
/**
* | output |
* | --- |
* | "Processing" |
*
* @param {Jobs_Status_ProcessingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_processing = /** @type {((inputs?: Jobs_Status_ProcessingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_ProcessingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_processing(inputs)
	return __ro.jobs_status_processing(inputs)
});
/**
* | output |
* | --- |
* | "Completed" |
*
* @param {Jobs_Status_CompletedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_completed = /** @type {((inputs?: Jobs_Status_CompletedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_CompletedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_completed(inputs)
	return __ro.jobs_status_completed(inputs)
});
/**
* | output |
* | --- |
* | "Failed" |
*
* @param {Jobs_Status_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_failed = /** @type {((inputs?: Jobs_Status_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_failed(inputs)
	return __ro.jobs_status_failed(inputs)
});
/**
* | output |
* | --- |
* | "Inline schema" |
*
* @param {Jobs_Inline_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_inline_schema = /** @type {((inputs?: Jobs_Inline_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Inline_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_inline_schema(inputs)
	return __ro.jobs_inline_schema(inputs)
});
/**
* | output |
* | --- |
* | "No schema" |
*
* @param {Jobs_No_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_no_schema = /** @type {((inputs?: Jobs_No_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_No_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_no_schema(inputs)
	return __ro.jobs_no_schema(inputs)
});
/**
* | output |
* | --- |
* | "Schema" |
*
* @param {Jobs_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_schema = /** @type {((inputs?: Jobs_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_schema(inputs)
	return __ro.jobs_schema(inputs)
});
/**
* | output |
* | --- |
* | "Select all jobs on this page" |
*
* @param {Jobs_Select_All_On_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_select_all_on_page = /** @type {((inputs?: Jobs_Select_All_On_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Select_All_On_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_select_all_on_page(inputs)
	return __ro.jobs_select_all_on_page(inputs)
});
/**
* | output |
* | --- |
* | "Select {name}" |
*
* @param {Jobs_Select_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_select_job = /** @type {((inputs: Jobs_Select_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Select_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_select_job(inputs)
	return __ro.jobs_select_job(inputs)
});
/**
* | output |
* | --- |
* | "Filename" |
*
* @param {Jobs_Filename_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_filename_column = /** @type {((inputs?: Jobs_Filename_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Filename_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_filename_column(inputs)
	return __ro.jobs_filename_column(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Jobs_Status_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_status_column = /** @type {((inputs?: Jobs_Status_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Status_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_status_column(inputs)
	return __ro.jobs_status_column(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Jobs_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_created_column = /** @type {((inputs?: Jobs_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_created_column(inputs)
	return __ro.jobs_created_column(inputs)
});
/**
* | output |
* | --- |
* | "File size" |
*
* @param {Jobs_File_Size_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_file_size_column = /** @type {((inputs?: Jobs_File_Size_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_File_Size_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_file_size_column(inputs)
	return __ro.jobs_file_size_column(inputs)
});
/**
* | output |
* | --- |
* | "Pages" |
*
* @param {Jobs_Pages_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_pages_column = /** @type {((inputs?: Jobs_Pages_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Pages_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_pages_column(inputs)
	return __ro.jobs_pages_column(inputs)
});
/**
* | output |
* | --- |
* | "New job" |
*
* @param {Jobs_New_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_new_job = /** @type {((inputs?: Jobs_New_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_New_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_new_job(inputs)
	return __ro.jobs_new_job(inputs)
});
/**
* | output |
* | --- |
* | "No jobs found" |
*
* @param {Jobs_No_Jobs_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_no_jobs_found = /** @type {((inputs?: Jobs_No_Jobs_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_No_Jobs_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_no_jobs_found(inputs)
	return __ro.jobs_no_jobs_found(inputs)
});
/**
* | output |
* | --- |
* | "Start a batch job to process documents and track progress here." |
*
* @param {Jobs_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_empty_body = /** @type {((inputs?: Jobs_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_empty_body(inputs)
	return __ro.jobs_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} job on this page." |
*
* @param {Jobs_Showing_Jobs_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_showing_jobs_one = /** @type {((inputs: Jobs_Showing_Jobs_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Showing_Jobs_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_showing_jobs_one(inputs)
	return __ro.jobs_showing_jobs_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} jobs on this page." |
*
* @param {Jobs_Showing_Jobs_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_showing_jobs_other = /** @type {((inputs: Jobs_Showing_Jobs_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Showing_Jobs_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_showing_jobs_other(inputs)
	return __ro.jobs_showing_jobs_other(inputs)
});
/**
* | output |
* | --- |
* | "No jobs to show." |
*
* @param {Jobs_No_Jobs_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_no_jobs_to_show = /** @type {((inputs?: Jobs_No_Jobs_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_No_Jobs_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_no_jobs_to_show(inputs)
	return __ro.jobs_no_jobs_to_show(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Jobs_Selected_Count_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_selected_count_one = /** @type {((inputs: Jobs_Selected_Count_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Selected_Count_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_selected_count_one(inputs)
	return __ro.jobs_selected_count_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Jobs_Selected_Count_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_selected_count_other = /** @type {((inputs: Jobs_Selected_Count_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Selected_Count_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_selected_count_other(inputs)
	return __ro.jobs_selected_count_other(inputs)
});
/**
* | output |
* | --- |
* | "Deleting..." |
*
* @param {Jobs_DeletingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_deleting = /** @type {((inputs?: Jobs_DeletingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_DeletingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_deleting(inputs)
	return __ro.jobs_deleting(inputs)
});
/**
* | output |
* | --- |
* | "Delete {name}" |
*
* @param {Jobs_Delete_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_delete_job = /** @type {((inputs: Jobs_Delete_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Delete_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_delete_job(inputs)
	return __ro.jobs_delete_job(inputs)
});
/**
* | output |
* | --- |
* | "Saved extraction schema" |
*
* @param {Jobs_Saved_Extraction_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_saved_extraction_schema = /** @type {((inputs?: Jobs_Saved_Extraction_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Saved_Extraction_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_saved_extraction_schema(inputs)
	return __ro.jobs_saved_extraction_schema(inputs)
});
/**
* | output |
* | --- |
* | "Schema submitted directly with this job." |
*
* @param {Jobs_Inline_Schema_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_inline_schema_description = /** @type {((inputs?: Jobs_Inline_Schema_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Inline_Schema_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_inline_schema_description(inputs)
	return __ro.jobs_inline_schema_description(inputs)
});
/**
* | output |
* | --- |
* | "Extraction schema details." |
*
* @param {Jobs_Extraction_Schema_DetailsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const jobs_extraction_schema_details = /** @type {((inputs?: Jobs_Extraction_Schema_DetailsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Jobs_Extraction_Schema_DetailsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.jobs_extraction_schema_details(inputs)
	return __ro.jobs_extraction_schema_details(inputs)
});
/**
* | output |
* | --- |
* | "Missing document id" |
*
* @param {New_Job_Missing_Document_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_missing_document_id = /** @type {((inputs?: New_Job_Missing_Document_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Missing_Document_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_missing_document_id(inputs)
	return __ro.new_job_missing_document_id(inputs)
});
/**
* | output |
* | --- |
* | "Failed to create OCR job" |
*
* @param {New_Job_Failed_CreateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed_create = /** @type {((inputs?: New_Job_Failed_CreateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Failed_CreateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed_create(inputs)
	return __ro.new_job_failed_create(inputs)
});
/**
* | output |
* | --- |
* | "Insufficient credits. Buy credits to process this document." |
*
* @param {New_Job_Insufficient_Credits_BuyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_insufficient_credits_buy = /** @type {((inputs?: New_Job_Insufficient_Credits_BuyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Insufficient_Credits_BuyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_insufficient_credits_buy(inputs)
	return __ro.new_job_insufficient_credits_buy(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load document" |
*
* @param {New_Job_Failed_Load_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed_load_document = /** @type {((inputs?: New_Job_Failed_Load_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Failed_Load_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed_load_document(inputs)
	return __ro.new_job_failed_load_document(inputs)
});
/**
* | output |
* | --- |
* | "Invalid document response" |
*
* @param {New_Job_Invalid_Document_ResponseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_invalid_document_response = /** @type {((inputs?: New_Job_Invalid_Document_ResponseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Invalid_Document_ResponseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_invalid_document_response(inputs)
	return __ro.new_job_invalid_document_response(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load schemas" |
*
* @param {New_Job_Failed_Load_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed_load_schemas = /** @type {((inputs?: New_Job_Failed_Load_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Failed_Load_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed_load_schemas(inputs)
	return __ro.new_job_failed_load_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Invalid schema response" |
*
* @param {New_Job_Invalid_Schema_ResponseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_invalid_schema_response = /** @type {((inputs?: New_Job_Invalid_Schema_ResponseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Invalid_Schema_ResponseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_invalid_schema_response(inputs)
	return __ro.new_job_invalid_schema_response(inputs)
});
/**
* | output |
* | --- |
* | "Invalid OCR job response" |
*
* @param {New_Job_Invalid_Job_ResponseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_invalid_job_response = /** @type {((inputs?: New_Job_Invalid_Job_ResponseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Invalid_Job_ResponseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_invalid_job_response(inputs)
	return __ro.new_job_invalid_job_response(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load OCR job" |
*
* @param {New_Job_Failed_Load_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed_load_job = /** @type {((inputs?: New_Job_Failed_Load_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Failed_Load_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed_load_job(inputs)
	return __ro.new_job_failed_load_job(inputs)
});
/**
* | output |
* | --- |
* | "Failed to poll OCR job" |
*
* @param {New_Job_Failed_Poll_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed_poll_job = /** @type {((inputs?: New_Job_Failed_Poll_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Failed_Poll_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed_poll_job(inputs)
	return __ro.new_job_failed_poll_job(inputs)
});
/**
* | output |
* | --- |
* | "Select Schema" |
*
* @param {New_Job_Select_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_select_schema = /** @type {((inputs?: New_Job_Select_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Select_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_select_schema(inputs)
	return __ro.new_job_select_schema(inputs)
});
/**
* | output |
* | --- |
* | "Select schema" |
*
* @param {New_Job_Select_Schema_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_select_schema_placeholder = /** @type {((inputs?: New_Job_Select_Schema_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Select_Schema_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_select_schema_placeholder(inputs)
	return __ro.new_job_select_schema_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Configure payload format" |
*
* @param {New_Job_Configure_Payload_FormatInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_configure_payload_format = /** @type {((inputs?: New_Job_Configure_Payload_FormatInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Configure_Payload_FormatInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_configure_payload_format(inputs)
	return __ro.new_job_configure_payload_format(inputs)
});
/**
* | output |
* | --- |
* | "Upload Documents" |
*
* @param {New_Job_Upload_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_upload_documents = /** @type {((inputs?: New_Job_Upload_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Upload_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_upload_documents(inputs)
	return __ro.new_job_upload_documents(inputs)
});
/**
* | output |
* | --- |
* | "{count} file selected" |
*
* @param {New_Job_Files_Selected_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_files_selected_one = /** @type {((inputs: New_Job_Files_Selected_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Files_Selected_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_files_selected_one(inputs)
	return __ro.new_job_files_selected_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} files selected" |
*
* @param {New_Job_Files_Selected_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_files_selected_other = /** @type {((inputs: New_Job_Files_Selected_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Files_Selected_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_files_selected_other(inputs)
	return __ro.new_job_files_selected_other(inputs)
});
/**
* | output |
* | --- |
* | "Drag or browse files" |
*
* @param {New_Job_Drag_Or_Browse_FilesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_drag_or_browse_files = /** @type {((inputs?: New_Job_Drag_Or_Browse_FilesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Drag_Or_Browse_FilesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_drag_or_browse_files(inputs)
	return __ro.new_job_drag_or_browse_files(inputs)
});
/**
* | output |
* | --- |
* | "Run & Monitor" |
*
* @param {New_Job_Run_MonitorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_run_monitor = /** @type {((inputs?: New_Job_Run_MonitorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Run_MonitorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_run_monitor(inputs)
	return __ro.new_job_run_monitor(inputs)
});
/**
* | output |
* | --- |
* | "Processing batch..." |
*
* @param {New_Job_Processing_BatchInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_processing_batch = /** @type {((inputs?: New_Job_Processing_BatchInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Processing_BatchInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_processing_batch(inputs)
	return __ro.new_job_processing_batch(inputs)
});
/**
* | output |
* | --- |
* | "Start extraction pipeline" |
*
* @param {New_Job_Start_Extraction_PipelineInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_start_extraction_pipeline = /** @type {((inputs?: New_Job_Start_Extraction_PipelineInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Start_Extraction_PipelineInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_start_extraction_pipeline(inputs)
	return __ro.new_job_start_extraction_pipeline(inputs)
});
/**
* | output |
* | --- |
* | "Select Extraction Schema" |
*
* @param {New_Job_Select_Extraction_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_select_extraction_schema = /** @type {((inputs?: New_Job_Select_Extraction_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Select_Extraction_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_select_extraction_schema(inputs)
	return __ro.new_job_select_extraction_schema(inputs)
});
/**
* | output |
* | --- |
* | "Choose a schema to define structured fields for AI extraction, or proceed in raw OCR-only text mode." |
*
* @param {New_Job_Select_Schema_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_select_schema_description = /** @type {((inputs?: New_Job_Select_Schema_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Select_Schema_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_select_schema_description(inputs)
	return __ro.new_job_select_schema_description(inputs)
});
/**
* | output |
* | --- |
* | "Select extraction schema" |
*
* @param {New_Job_Select_Extraction_Schema_AriaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_select_extraction_schema_aria = /** @type {((inputs?: New_Job_Select_Extraction_Schema_AriaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Select_Extraction_Schema_AriaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_select_extraction_schema_aria(inputs)
	return __ro.new_job_select_extraction_schema_aria(inputs)
});
/**
* | output |
* | --- |
* | "Search schemas..." |
*
* @param {New_Job_Search_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_search_schemas = /** @type {((inputs?: New_Job_Search_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Search_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_search_schemas(inputs)
	return __ro.new_job_search_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Loading schemas" |
*
* @param {New_Job_Loading_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_loading_schemas = /** @type {((inputs?: New_Job_Loading_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Loading_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_loading_schemas(inputs)
	return __ro.new_job_loading_schemas(inputs)
});
/**
* | output |
* | --- |
* | "No schemas found." |
*
* @param {New_Job_No_Schemas_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_schemas_found = /** @type {((inputs?: New_Job_No_Schemas_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Schemas_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_schemas_found(inputs)
	return __ro.new_job_no_schemas_found(inputs)
});
/**
* | output |
* | --- |
* | "No schema (OCR only)" |
*
* @param {New_Job_No_Schema_Ocr_OnlyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_schema_ocr_only = /** @type {((inputs?: New_Job_No_Schema_Ocr_OnlyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Schema_Ocr_OnlyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_schema_ocr_only(inputs)
	return __ro.new_job_no_schema_ocr_only(inputs)
});
/**
* | output |
* | --- |
* | "Process documents without structured extraction." |
*
* @param {New_Job_No_Schema_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_schema_description = /** @type {((inputs?: New_Job_No_Schema_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Schema_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_schema_description(inputs)
	return __ro.new_job_no_schema_description(inputs)
});
/**
* | output |
* | --- |
* | "No personal schemas are available." |
*
* @param {New_Job_No_Personal_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_personal_schemas = /** @type {((inputs?: New_Job_No_Personal_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Personal_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_personal_schemas(inputs)
	return __ro.new_job_no_personal_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Create one" |
*
* @param {New_Job_Create_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_create_one = /** @type {((inputs?: New_Job_Create_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Create_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_create_one(inputs)
	return __ro.new_job_create_one(inputs)
});
/**
* | output |
* | --- |
* | "Selected schema defines the structured fields that will be extracted from your files." |
*
* @param {New_Job_Selected_Schema_HelpInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_selected_schema_help = /** @type {((inputs?: New_Job_Selected_Schema_HelpInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Selected_Schema_HelpInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_selected_schema_help(inputs)
	return __ro.new_job_selected_schema_help(inputs)
});
/**
* | output |
* | --- |
* | "No schema selected. Files will be OCRized without structured extraction." |
*
* @param {New_Job_No_Schema_Selected_HelpInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_schema_selected_help = /** @type {((inputs?: New_Job_No_Schema_Selected_HelpInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Schema_Selected_HelpInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_schema_selected_help(inputs)
	return __ro.new_job_no_schema_selected_help(inputs)
});
/**
* | output |
* | --- |
* | "Target Mapped Fields ({count})" |
*
* @param {New_Job_Target_Mapped_FieldsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_target_mapped_fields = /** @type {((inputs: New_Job_Target_Mapped_FieldsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Target_Mapped_FieldsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_target_mapped_fields(inputs)
	return __ro.new_job_target_mapped_fields(inputs)
});
/**
* | output |
* | --- |
* | "No fields defined in this schema." |
*
* @param {New_Job_No_Fields_DefinedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_fields_defined = /** @type {((inputs?: New_Job_No_Fields_DefinedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Fields_DefinedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_fields_defined(inputs)
	return __ro.new_job_no_fields_defined(inputs)
});
/**
* | output |
* | --- |
* | "OCR-Only Mode Active" |
*
* @param {New_Job_Ocr_Only_Mode_ActiveInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_ocr_only_mode_active = /** @type {((inputs?: New_Job_Ocr_Only_Mode_ActiveInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Ocr_Only_Mode_ActiveInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_ocr_only_mode_active(inputs)
	return __ro.new_job_ocr_only_mode_active(inputs)
});
/**
* | output |
* | --- |
* | "Documents will be processed for high-fidelity OCR text extraction without converting fields into structured schema payloads. Select a schema above to extract..." |
*
* @param {New_Job_Ocr_Only_Mode_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_ocr_only_mode_body = /** @type {((inputs?: New_Job_Ocr_Only_Mode_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Ocr_Only_Mode_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_ocr_only_mode_body(inputs)
	return __ro.new_job_ocr_only_mode_body(inputs)
});
/**
* | output |
* | --- |
* | "Select PDF or image files to extract content. You can batch upload up to {count} files simultaneously." |
*
* @param {New_Job_Upload_Documents_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_upload_documents_description = /** @type {((inputs: New_Job_Upload_Documents_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Upload_Documents_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_upload_documents_description(inputs)
	return __ro.new_job_upload_documents_description(inputs)
});
/**
* | output |
* | --- |
* | "Drag & drop files or click to upload" |
*
* @param {New_Job_Dropzone_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_dropzone_title = /** @type {((inputs?: New_Job_Dropzone_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Dropzone_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_dropzone_title(inputs)
	return __ro.new_job_dropzone_title(inputs)
});
/**
* | output |
* | --- |
* | "Supports PDF, PNG, and JPG up to {size} per file" |
*
* @param {New_Job_Dropzone_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_dropzone_description = /** @type {((inputs: New_Job_Dropzone_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Dropzone_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_dropzone_description(inputs)
	return __ro.new_job_dropzone_description(inputs)
});
/**
* | output |
* | --- |
* | "Browse Files" |
*
* @param {New_Job_Browse_FilesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_browse_files = /** @type {((inputs?: New_Job_Browse_FilesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Browse_FilesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_browse_files(inputs)
	return __ro.new_job_browse_files(inputs)
});
/**
* | output |
* | --- |
* | "Pending Upload Queue ({count})" |
*
* @param {New_Job_Pending_Upload_QueueInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_pending_upload_queue = /** @type {((inputs: New_Job_Pending_Upload_QueueInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Pending_Upload_QueueInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_pending_upload_queue(inputs)
	return __ro.new_job_pending_upload_queue(inputs)
});
/**
* | output |
* | --- |
* | "Clear All" |
*
* @param {New_Job_Clear_AllInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_clear_all = /** @type {((inputs?: New_Job_Clear_AllInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Clear_AllInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_clear_all(inputs)
	return __ro.new_job_clear_all(inputs)
});
/**
* | output |
* | --- |
* | "Remove file" |
*
* @param {New_Job_Remove_FileInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_remove_file = /** @type {((inputs?: New_Job_Remove_FileInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Remove_FileInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_remove_file(inputs)
	return __ro.new_job_remove_file(inputs)
});
/**
* | output |
* | --- |
* | "Extraction Queue & Results" |
*
* @param {New_Job_Extraction_Queue_ResultsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_extraction_queue_results = /** @type {((inputs?: New_Job_Extraction_Queue_ResultsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Extraction_Queue_ResultsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_extraction_queue_results(inputs)
	return __ro.new_job_extraction_queue_results(inputs)
});
/**
* | output |
* | --- |
* | "{count} File" |
*
* @param {New_Job_File_Count_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_file_count_one = /** @type {((inputs: New_Job_File_Count_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_File_Count_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_file_count_one(inputs)
	return __ro.new_job_file_count_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} Files" |
*
* @param {New_Job_File_Count_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_file_count_other = /** @type {((inputs: New_Job_File_Count_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_File_Count_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_file_count_other(inputs)
	return __ro.new_job_file_count_other(inputs)
});
/**
* | output |
* | --- |
* | "{label} total" |
*
* @param {New_Job_TotalInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_total = /** @type {((inputs: New_Job_TotalInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_TotalInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_total(inputs)
	return __ro.new_job_total(inputs)
});
/**
* | output |
* | --- |
* | "Active Batch Status" |
*
* @param {New_Job_Active_Batch_StatusInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_active_batch_status = /** @type {((inputs?: New_Job_Active_Batch_StatusInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Active_Batch_StatusInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_active_batch_status(inputs)
	return __ro.new_job_active_batch_status(inputs)
});
/**
* | output |
* | --- |
* | "Monitors your batch documents in real-time execution." |
*
* @param {New_Job_Active_Batch_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_active_batch_description = /** @type {((inputs?: New_Job_Active_Batch_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Active_Batch_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_active_batch_description(inputs)
	return __ro.new_job_active_batch_description(inputs)
});
/**
* | output |
* | --- |
* | "Progress: {progress}%" |
*
* @param {New_Job_ProgressInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_progress = /** @type {((inputs: New_Job_ProgressInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_ProgressInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_progress(inputs)
	return __ro.new_job_progress(inputs)
});
/**
* | output |
* | --- |
* | "Total Files" |
*
* @param {New_Job_Total_FilesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_total_files = /** @type {((inputs?: New_Job_Total_FilesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Total_FilesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_total_files(inputs)
	return __ro.new_job_total_files(inputs)
});
/**
* | output |
* | --- |
* | "Completed" |
*
* @param {New_Job_CompletedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_completed = /** @type {((inputs?: New_Job_CompletedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_CompletedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_completed(inputs)
	return __ro.new_job_completed(inputs)
});
/**
* | output |
* | --- |
* | "Processing" |
*
* @param {New_Job_ProcessingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_processing = /** @type {((inputs?: New_Job_ProcessingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_ProcessingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_processing(inputs)
	return __ro.new_job_processing(inputs)
});
/**
* | output |
* | --- |
* | "Failed" |
*
* @param {New_Job_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_failed = /** @type {((inputs?: New_Job_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_failed(inputs)
	return __ro.new_job_failed(inputs)
});
/**
* | output |
* | --- |
* | "No active extraction jobs" |
*
* @param {New_Job_No_Active_Extraction_JobsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_active_extraction_jobs = /** @type {((inputs?: New_Job_No_Active_Extraction_JobsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Active_Extraction_JobsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_active_extraction_jobs(inputs)
	return __ro.new_job_no_active_extraction_jobs(inputs)
});
/**
* | output |
* | --- |
* | "Upload documents above and select a schema to launch the automated OCR and extraction process." |
*
* @param {New_Job_No_Active_Extraction_Jobs_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_no_active_extraction_jobs_body = /** @type {((inputs?: New_Job_No_Active_Extraction_Jobs_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_No_Active_Extraction_Jobs_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_no_active_extraction_jobs_body(inputs)
	return __ro.new_job_no_active_extraction_jobs_body(inputs)
});
/**
* | output |
* | --- |
* | "Preview document" |
*
* @param {New_Job_Preview_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_preview_document = /** @type {((inputs?: New_Job_Preview_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Preview_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_preview_document(inputs)
	return __ro.new_job_preview_document(inputs)
});
/**
* | output |
* | --- |
* | "Document preview is not available yet" |
*
* @param {New_Job_Preview_UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_preview_unavailable = /** @type {((inputs?: New_Job_Preview_UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Preview_UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_preview_unavailable(inputs)
	return __ro.new_job_preview_unavailable(inputs)
});
/**
* | output |
* | --- |
* | "Remove failed job" |
*
* @param {New_Job_Remove_Failed_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_remove_failed_job = /** @type {((inputs?: New_Job_Remove_Failed_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Remove_Failed_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_remove_failed_job(inputs)
	return __ro.new_job_remove_failed_job(inputs)
});
/**
* | output |
* | --- |
* | "Queueing Documents..." |
*
* @param {New_Job_Queueing_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_queueing_documents = /** @type {((inputs?: New_Job_Queueing_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Queueing_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_queueing_documents(inputs)
	return __ro.new_job_queueing_documents(inputs)
});
/**
* | output |
* | --- |
* | "Extracting Content..." |
*
* @param {New_Job_Extracting_ContentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_extracting_content = /** @type {((inputs?: New_Job_Extracting_ContentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Extracting_ContentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_extracting_content(inputs)
	return __ro.new_job_extracting_content(inputs)
});
/**
* | output |
* | --- |
* | "Run Extraction ({count} File)" |
*
* @param {New_Job_Run_Extraction_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_run_extraction_one = /** @type {((inputs: New_Job_Run_Extraction_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Run_Extraction_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_run_extraction_one(inputs)
	return __ro.new_job_run_extraction_one(inputs)
});
/**
* | output |
* | --- |
* | "Run Extraction ({count} Files)" |
*
* @param {New_Job_Run_Extraction_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_run_extraction_other = /** @type {((inputs: New_Job_Run_Extraction_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Run_Extraction_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_run_extraction_other(inputs)
	return __ro.new_job_run_extraction_other(inputs)
});
/**
* | output |
* | --- |
* | "Insufficient credits for this document." |
*
* @param {New_Job_Insufficient_Credits_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_insufficient_credits_document = /** @type {((inputs?: New_Job_Insufficient_Credits_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Insufficient_Credits_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_insufficient_credits_document(inputs)
	return __ro.new_job_insufficient_credits_document(inputs)
});
/**
* | output |
* | --- |
* | "Processing failed" |
*
* @param {New_Job_Processing_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_processing_failed = /** @type {((inputs?: New_Job_Processing_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Processing_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_processing_failed(inputs)
	return __ro.new_job_processing_failed(inputs)
});
/**
* | output |
* | --- |
* | "Processed" |
*
* @param {New_Job_ProcessedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_processed = /** @type {((inputs?: New_Job_ProcessedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_ProcessedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_processed(inputs)
	return __ro.new_job_processed(inputs)
});
/**
* | output |
* | --- |
* | "Document {id}" |
*
* @param {New_Job_Document_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_document_id = /** @type {((inputs: New_Job_Document_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Document_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_document_id(inputs)
	return __ro.new_job_document_id(inputs)
});
/**
* | output |
* | --- |
* | "Creating OCR job..." |
*
* @param {New_Job_Creating_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_creating_job = /** @type {((inputs?: New_Job_Creating_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Creating_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_creating_job(inputs)
	return __ro.new_job_creating_job(inputs)
});
/**
* | output |
* | --- |
* | "Queued for processing..." |
*
* @param {New_Job_Queued_ProcessingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_queued_processing = /** @type {((inputs?: New_Job_Queued_ProcessingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Queued_ProcessingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_queued_processing(inputs)
	return __ro.new_job_queued_processing(inputs)
});
/**
* | output |
* | --- |
* | "Extracting entities..." |
*
* @param {New_Job_Extracting_EntitiesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const new_job_extracting_entities = /** @type {((inputs?: New_Job_Extracting_EntitiesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<New_Job_Extracting_EntitiesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.new_job_extracting_entities(inputs)
	return __ro.new_job_extracting_entities(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Common_ApplyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_apply = /** @type {((inputs?: Common_ApplyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ApplyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_apply(inputs)
	return __ro.common_apply(inputs)
});
/**
* | output |
* | --- |
* | "Clear" |
*
* @param {Common_ClearInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_clear = /** @type {((inputs?: Common_ClearInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ClearInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_clear(inputs)
	return __ro.common_clear(inputs)
});
/**
* | output |
* | --- |
* | "Saving..." |
*
* @param {Common_SavingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_saving = /** @type {((inputs?: Common_SavingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_SavingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_saving(inputs)
	return __ro.common_saving(inputs)
});
/**
* | output |
* | --- |
* | "Loading..." |
*
* @param {Common_LoadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_loading = /** @type {((inputs?: Common_LoadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_LoadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_loading(inputs)
	return __ro.common_loading(inputs)
});
/**
* | output |
* | --- |
* | "Refresh" |
*
* @param {Common_RefreshInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_refresh = /** @type {((inputs?: Common_RefreshInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_RefreshInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_refresh(inputs)
	return __ro.common_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Connected" |
*
* @param {Common_ConnectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_connected = /** @type {((inputs?: Common_ConnectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ConnectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_connected(inputs)
	return __ro.common_connected(inputs)
});
/**
* | output |
* | --- |
* | "Connect" |
*
* @param {Common_ConnectInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_connect = /** @type {((inputs?: Common_ConnectInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ConnectInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_connect(inputs)
	return __ro.common_connect(inputs)
});
/**
* | output |
* | --- |
* | "Download" |
*
* @param {Common_DownloadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_download = /** @type {((inputs?: Common_DownloadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_DownloadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_download(inputs)
	return __ro.common_download(inputs)
});
/**
* | output |
* | --- |
* | "Today" |
*
* @param {Common_TodayInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_today = /** @type {((inputs?: Common_TodayInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_TodayInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_today(inputs)
	return __ro.common_today(inputs)
});
/**
* | output |
* | --- |
* | "This week" |
*
* @param {Common_This_WeekInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_this_week = /** @type {((inputs?: Common_This_WeekInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_This_WeekInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_this_week(inputs)
	return __ro.common_this_week(inputs)
});
/**
* | output |
* | --- |
* | "This month" |
*
* @param {Common_This_MonthInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_this_month = /** @type {((inputs?: Common_This_MonthInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_This_MonthInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_this_month(inputs)
	return __ro.common_this_month(inputs)
});
/**
* | output |
* | --- |
* | "Any" |
*
* @param {Common_AnyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const common_any = /** @type {((inputs?: Common_AnyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_AnyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_any(inputs)
	return __ro.common_any(inputs)
});
/**
* | output |
* | --- |
* | "Unavailable" |
*
* @param {Billing_UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_unavailable = /** @type {((inputs?: Billing_UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_unavailable(inputs)
	return __ro.billing_unavailable(inputs)
});
/**
* | output |
* | --- |
* | "Credits must be purchased in 1000-credit blocks." |
*
* @param {Billing_Credit_Blocks_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_credit_blocks_error = /** @type {((inputs?: Billing_Credit_Blocks_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Credit_Blocks_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_credit_blocks_error(inputs)
	return __ro.billing_credit_blocks_error(inputs)
});
/**
* | output |
* | --- |
* | "Unable to start checkout" |
*
* @param {Billing_Checkout_UnavailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_checkout_unavailable = /** @type {((inputs?: Billing_Checkout_UnavailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Checkout_UnavailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_checkout_unavailable(inputs)
	return __ro.billing_checkout_unavailable(inputs)
});
/**
* | output |
* | --- |
* | "Payment Received" |
*
* @param {Billing_Payment_Received_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_payment_received_title = /** @type {((inputs?: Billing_Payment_Received_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Payment_Received_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_payment_received_title(inputs)
	return __ro.billing_payment_received_title(inputs)
});
/**
* | output |
* | --- |
* | "Your credit balance will update momentarily once we receive the payment confirmation." |
*
* @param {Billing_Payment_Received_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_payment_received_body = /** @type {((inputs?: Billing_Payment_Received_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Payment_Received_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_payment_received_body(inputs)
	return __ro.billing_payment_received_body(inputs)
});
/**
* | output |
* | --- |
* | "Checkout Canceled" |
*
* @param {Billing_Checkout_Canceled_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_checkout_canceled_title = /** @type {((inputs?: Billing_Checkout_Canceled_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Checkout_Canceled_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_checkout_canceled_title(inputs)
	return __ro.billing_checkout_canceled_title(inputs)
});
/**
* | output |
* | --- |
* | "No credits were purchased and no charges were made." |
*
* @param {Billing_Checkout_Canceled_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_checkout_canceled_body = /** @type {((inputs?: Billing_Checkout_Canceled_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Checkout_Canceled_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_checkout_canceled_body(inputs)
	return __ro.billing_checkout_canceled_body(inputs)
});
/**
* | output |
* | --- |
* | "Available Balance" |
*
* @param {Billing_Available_BalanceInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_available_balance = /** @type {((inputs?: Billing_Available_BalanceInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Available_BalanceInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_available_balance(inputs)
	return __ro.billing_available_balance(inputs)
});
/**
* | output |
* | --- |
* | "Conversion" |
*
* @param {Billing_ConversionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_conversion = /** @type {((inputs?: Billing_ConversionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_ConversionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_conversion(inputs)
	return __ro.billing_conversion(inputs)
});
/**
* | output |
* | --- |
* | "1 credit = 1 page" |
*
* @param {Billing_Conversion_RateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_conversion_rate = /** @type {((inputs?: Billing_Conversion_RateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Conversion_RateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_conversion_rate(inputs)
	return __ro.billing_conversion_rate(inputs)
});
/**
* | output |
* | --- |
* | "Balance Checked on Upload" |
*
* @param {Billing_Balance_Checked_UploadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_balance_checked_upload = /** @type {((inputs?: Billing_Balance_Checked_UploadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Balance_Checked_UploadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_balance_checked_upload(inputs)
	return __ro.billing_balance_checked_upload(inputs)
});
/**
* | output |
* | --- |
* | "Debited After Success" |
*
* @param {Billing_Debited_After_SuccessInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_debited_after_success = /** @type {((inputs?: Billing_Debited_After_SuccessInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Debited_After_SuccessInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_debited_after_success(inputs)
	return __ro.billing_debited_after_success(inputs)
});
/**
* | output |
* | --- |
* | "Secure Stripe Checkout" |
*
* @param {Billing_Secure_Stripe_CheckoutInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_secure_stripe_checkout = /** @type {((inputs?: Billing_Secure_Stripe_CheckoutInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Secure_Stripe_CheckoutInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_secure_stripe_checkout(inputs)
	return __ro.billing_secure_stripe_checkout(inputs)
});
/**
* | output |
* | --- |
* | "Purchase Credits" |
*
* @param {Billing_Purchase_CreditsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_purchase_credits = /** @type {((inputs?: Billing_Purchase_CreditsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Purchase_CreditsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_purchase_credits(inputs)
	return __ro.billing_purchase_credits(inputs)
});
/**
* | output |
* | --- |
* | "Credits to purchase" |
*
* @param {Billing_Credits_To_PurchaseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_credits_to_purchase = /** @type {((inputs?: Billing_Credits_To_PurchaseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Credits_To_PurchaseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_credits_to_purchase(inputs)
	return __ro.billing_credits_to_purchase(inputs)
});
/**
* | output |
* | --- |
* | "Volume Discount Tiers" |
*
* @param {Billing_Volume_Discount_TiersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_volume_discount_tiers = /** @type {((inputs?: Billing_Volume_Discount_TiersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Volume_Discount_TiersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_volume_discount_tiers(inputs)
	return __ro.billing_volume_discount_tiers(inputs)
});
/**
* | output |
* | --- |
* | "Total to pay" |
*
* @param {Billing_Total_To_PayInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_total_to_pay = /** @type {((inputs?: Billing_Total_To_PayInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Total_To_PayInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_total_to_pay(inputs)
	return __ro.billing_total_to_pay(inputs)
});
/**
* | output |
* | --- |
* | "Base price" |
*
* @param {Billing_Base_PriceInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_base_price = /** @type {((inputs?: Billing_Base_PriceInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Base_PriceInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_base_price(inputs)
	return __ro.billing_base_price(inputs)
});
/**
* | output |
* | --- |
* | "Volume Discount" |
*
* @param {Billing_Volume_DiscountInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_volume_discount = /** @type {((inputs?: Billing_Volume_DiscountInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Volume_DiscountInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_volume_discount(inputs)
	return __ro.billing_volume_discount(inputs)
});
/**
* | output |
* | --- |
* | "Starting checkout..." |
*
* @param {Billing_Starting_CheckoutInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_starting_checkout = /** @type {((inputs?: Billing_Starting_CheckoutInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Starting_CheckoutInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_starting_checkout(inputs)
	return __ro.billing_starting_checkout(inputs)
});
/**
* | output |
* | --- |
* | "Secure Checkout" |
*
* @param {Billing_Secure_CheckoutInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_secure_checkout = /** @type {((inputs?: Billing_Secure_CheckoutInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Secure_CheckoutInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_secure_checkout(inputs)
	return __ro.billing_secure_checkout(inputs)
});
/**
* | output |
* | --- |
* | "Buy Credits" |
*
* @param {Billing_Buy_CreditsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_buy_credits = /** @type {((inputs?: Billing_Buy_CreditsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Buy_CreditsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_buy_credits(inputs)
	return __ro.billing_buy_credits(inputs)
});
/**
* | output |
* | --- |
* | "Billing Orders" |
*
* @param {Billing_Orders_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_page_title = /** @type {((inputs?: Billing_Orders_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_page_title(inputs)
	return __ro.billing_orders_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Order date" |
*
* @param {Billing_Orders_Order_Date_FilterInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_order_date_filter = /** @type {((inputs?: Billing_Orders_Order_Date_FilterInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Order_Date_FilterInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_order_date_filter(inputs)
	return __ro.billing_orders_order_date_filter(inputs)
});
/**
* | output |
* | --- |
* | "Amount" |
*
* @param {Billing_Orders_Amount_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_amount_column = /** @type {((inputs?: Billing_Orders_Amount_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Amount_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_amount_column(inputs)
	return __ro.billing_orders_amount_column(inputs)
});
/**
* | output |
* | --- |
* | "Credits" |
*
* @param {Billing_Orders_Credits_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_credits_column = /** @type {((inputs?: Billing_Orders_Credits_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Credits_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_credits_column(inputs)
	return __ro.billing_orders_credits_column(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Billing_Orders_Status_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_status_column = /** @type {((inputs?: Billing_Orders_Status_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Status_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_status_column(inputs)
	return __ro.billing_orders_status_column(inputs)
});
/**
* | output |
* | --- |
* | "Payment datetime" |
*
* @param {Billing_Orders_Payment_Datetime_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_payment_datetime_column = /** @type {((inputs?: Billing_Orders_Payment_Datetime_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Payment_Datetime_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_payment_datetime_column(inputs)
	return __ro.billing_orders_payment_datetime_column(inputs)
});
/**
* | output |
* | --- |
* | "Invoice" |
*
* @param {Billing_Orders_Invoice_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_invoice_column = /** @type {((inputs?: Billing_Orders_Invoice_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Invoice_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_invoice_column(inputs)
	return __ro.billing_orders_invoice_column(inputs)
});
/**
* | output |
* | --- |
* | "Presets" |
*
* @param {Billing_Orders_PresetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_presets = /** @type {((inputs?: Billing_Orders_PresetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_PresetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_presets(inputs)
	return __ro.billing_orders_presets(inputs)
});
/**
* | output |
* | --- |
* | "Filter by status" |
*
* @param {Billing_Orders_Filter_StatusInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_filter_status = /** @type {((inputs?: Billing_Orders_Filter_StatusInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Filter_StatusInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_filter_status(inputs)
	return __ro.billing_orders_filter_status(inputs)
});
/**
* | output |
* | --- |
* | "All Orders" |
*
* @param {Billing_Orders_All_OrdersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_all_orders = /** @type {((inputs?: Billing_Orders_All_OrdersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_All_OrdersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_all_orders(inputs)
	return __ro.billing_orders_all_orders(inputs)
});
/**
* | output |
* | --- |
* | "Clear Filters" |
*
* @param {Billing_Orders_Clear_FiltersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_clear_filters = /** @type {((inputs?: Billing_Orders_Clear_FiltersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Clear_FiltersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_clear_filters(inputs)
	return __ro.billing_orders_clear_filters(inputs)
});
/**
* | output |
* | --- |
* | "Clear filters" |
*
* @param {Billing_Orders_Clear_Filters_ActionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_clear_filters_action = /** @type {((inputs?: Billing_Orders_Clear_Filters_ActionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Clear_Filters_ActionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_clear_filters_action(inputs)
	return __ro.billing_orders_clear_filters_action(inputs)
});
/**
* | output |
* | --- |
* | "No billing orders found" |
*
* @param {Billing_Orders_No_Orders_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_no_orders_found = /** @type {((inputs?: Billing_Orders_No_Orders_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_No_Orders_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_no_orders_found(inputs)
	return __ro.billing_orders_no_orders_found(inputs)
});
/**
* | output |
* | --- |
* | "No billing orders yet" |
*
* @param {Billing_Orders_No_Orders_YetInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_no_orders_yet = /** @type {((inputs?: Billing_Orders_No_Orders_YetInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_No_Orders_YetInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_no_orders_yet(inputs)
	return __ro.billing_orders_no_orders_yet(inputs)
});
/**
* | output |
* | --- |
* | "No billing orders match the selected filters." |
*
* @param {Billing_Orders_No_Orders_MatchInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_no_orders_match = /** @type {((inputs?: Billing_Orders_No_Orders_MatchInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_No_Orders_MatchInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_no_orders_match(inputs)
	return __ro.billing_orders_no_orders_match(inputs)
});
/**
* | output |
* | --- |
* | "Credit purchase orders will appear here after checkout starts." |
*
* @param {Billing_Orders_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_empty_body = /** @type {((inputs?: Billing_Orders_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_empty_body(inputs)
	return __ro.billing_orders_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} order on this page." |
*
* @param {Billing_Orders_Showing_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_showing_one = /** @type {((inputs: Billing_Orders_Showing_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Showing_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_showing_one(inputs)
	return __ro.billing_orders_showing_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} orders on this page." |
*
* @param {Billing_Orders_Showing_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_showing_other = /** @type {((inputs: Billing_Orders_Showing_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Showing_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_showing_other(inputs)
	return __ro.billing_orders_showing_other(inputs)
});
/**
* | output |
* | --- |
* | "No billing orders to show." |
*
* @param {Billing_Orders_None_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_none_to_show = /** @type {((inputs?: Billing_Orders_None_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_None_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_none_to_show(inputs)
	return __ro.billing_orders_none_to_show(inputs)
});
/**
* | output |
* | --- |
* | "Sort by order date ascending" |
*
* @param {Billing_Orders_Sort_Order_Date_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_sort_order_date_ascending = /** @type {((inputs?: Billing_Orders_Sort_Order_Date_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Sort_Order_Date_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_sort_order_date_ascending(inputs)
	return __ro.billing_orders_sort_order_date_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by order date descending" |
*
* @param {Billing_Orders_Sort_Order_Date_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_sort_order_date_descending = /** @type {((inputs?: Billing_Orders_Sort_Order_Date_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Sort_Order_Date_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_sort_order_date_descending(inputs)
	return __ro.billing_orders_sort_order_date_descending(inputs)
});
/**
* | output |
* | --- |
* | "Pending" |
*
* @param {Billing_Order_Status_PendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_order_status_pending = /** @type {((inputs?: Billing_Order_Status_PendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Order_Status_PendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_order_status_pending(inputs)
	return __ro.billing_order_status_pending(inputs)
});
/**
* | output |
* | --- |
* | "Paid" |
*
* @param {Billing_Order_Status_PaidInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_order_status_paid = /** @type {((inputs?: Billing_Order_Status_PaidInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Order_Status_PaidInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_order_status_paid(inputs)
	return __ro.billing_order_status_paid(inputs)
});
/**
* | output |
* | --- |
* | "Failed" |
*
* @param {Billing_Order_Status_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_order_status_failed = /** @type {((inputs?: Billing_Order_Status_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Order_Status_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_order_status_failed(inputs)
	return __ro.billing_order_status_failed(inputs)
});
/**
* | output |
* | --- |
* | "Refunded" |
*
* @param {Billing_Order_Status_RefundedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_order_status_refunded = /** @type {((inputs?: Billing_Order_Status_RefundedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Order_Status_RefundedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_order_status_refunded(inputs)
	return __ro.billing_order_status_refunded(inputs)
});
/**
* | output |
* | --- |
* | "Canceled" |
*
* @param {Billing_Order_Status_CanceledInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_order_status_canceled = /** @type {((inputs?: Billing_Order_Status_CanceledInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Order_Status_CanceledInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_order_status_canceled(inputs)
	return __ro.billing_order_status_canceled(inputs)
});
/**
* | output |
* | --- |
* | "Preview {invoice} PDF" |
*
* @param {Billing_Orders_Invoice_Pdf_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_invoice_pdf_title = /** @type {((inputs: Billing_Orders_Invoice_Pdf_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Invoice_Pdf_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_invoice_pdf_title(inputs)
	return __ro.billing_orders_invoice_pdf_title(inputs)
});
/**
* | output |
* | --- |
* | "Invoice {invoice}" |
*
* @param {Billing_Orders_Invoice_Preview_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_invoice_preview_title = /** @type {((inputs: Billing_Orders_Invoice_Preview_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Invoice_Preview_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_invoice_preview_title(inputs)
	return __ro.billing_orders_invoice_preview_title(inputs)
});
/**
* | output |
* | --- |
* | "PDF preview" |
*
* @param {Billing_Orders_Invoice_Preview_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_invoice_preview_description = /** @type {((inputs?: Billing_Orders_Invoice_Preview_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Invoice_Preview_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_invoice_preview_description(inputs)
	return __ro.billing_orders_invoice_preview_description(inputs)
});
/**
* | output |
* | --- |
* | "Invoice {invoice} PDF preview" |
*
* @param {Billing_Orders_Invoice_Iframe_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_invoice_iframe_title = /** @type {((inputs: Billing_Orders_Invoice_Iframe_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Invoice_Iframe_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_invoice_iframe_title(inputs)
	return __ro.billing_orders_invoice_iframe_title(inputs)
});
/**
* | output |
* | --- |
* | "Download" |
*
* @param {Billing_Orders_Download_InvoiceInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_orders_download_invoice = /** @type {((inputs?: Billing_Orders_Download_InvoiceInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Orders_Download_InvoiceInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_orders_download_invoice(inputs)
	return __ro.billing_orders_download_invoice(inputs)
});
/**
* | output |
* | --- |
* | "Credit Usage History" |
*
* @param {Credit_Usage_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_page_title = /** @type {((inputs?: Credit_Usage_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_page_title(inputs)
	return __ro.credit_usage_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Date range" |
*
* @param {Credit_Usage_Date_Range_FilterInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_date_range_filter = /** @type {((inputs?: Credit_Usage_Date_Range_FilterInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Date_Range_FilterInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_date_range_filter(inputs)
	return __ro.credit_usage_date_range_filter(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Credit_Usage_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_created_column = /** @type {((inputs?: Credit_Usage_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_created_column(inputs)
	return __ro.credit_usage_created_column(inputs)
});
/**
* | output |
* | --- |
* | "Type" |
*
* @param {Credit_Usage_Type_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_type_column = /** @type {((inputs?: Credit_Usage_Type_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Type_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_type_column(inputs)
	return __ro.credit_usage_type_column(inputs)
});
/**
* | output |
* | --- |
* | "Credits" |
*
* @param {Credit_Usage_Credits_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_credits_column = /** @type {((inputs?: Credit_Usage_Credits_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Credits_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_credits_column(inputs)
	return __ro.credit_usage_credits_column(inputs)
});
/**
* | output |
* | --- |
* | "Related ID" |
*
* @param {Credit_Usage_Related_Id_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_related_id_column = /** @type {((inputs?: Credit_Usage_Related_Id_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Related_Id_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_related_id_column(inputs)
	return __ro.credit_usage_related_id_column(inputs)
});
/**
* | output |
* | --- |
* | "Filter by type" |
*
* @param {Credit_Usage_Filter_TypeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_filter_type = /** @type {((inputs?: Credit_Usage_Filter_TypeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Filter_TypeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_filter_type(inputs)
	return __ro.credit_usage_filter_type(inputs)
});
/**
* | output |
* | --- |
* | "All Activity" |
*
* @param {Credit_Usage_All_ActivityInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_all_activity = /** @type {((inputs?: Credit_Usage_All_ActivityInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_All_ActivityInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_all_activity(inputs)
	return __ro.credit_usage_all_activity(inputs)
});
/**
* | output |
* | --- |
* | "Purchase" |
*
* @param {Credit_Usage_Type_PurchaseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_type_purchase = /** @type {((inputs?: Credit_Usage_Type_PurchaseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Type_PurchaseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_type_purchase(inputs)
	return __ro.credit_usage_type_purchase(inputs)
});
/**
* | output |
* | --- |
* | "Debit" |
*
* @param {Credit_Usage_Type_DebitInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_type_debit = /** @type {((inputs?: Credit_Usage_Type_DebitInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Type_DebitInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_type_debit(inputs)
	return __ro.credit_usage_type_debit(inputs)
});
/**
* | output |
* | --- |
* | "No credit usage found" |
*
* @param {Credit_Usage_No_Usage_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_no_usage_found = /** @type {((inputs?: Credit_Usage_No_Usage_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_No_Usage_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_no_usage_found(inputs)
	return __ro.credit_usage_no_usage_found(inputs)
});
/**
* | output |
* | --- |
* | "No credit usage yet" |
*
* @param {Credit_Usage_No_Usage_YetInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_no_usage_yet = /** @type {((inputs?: Credit_Usage_No_Usage_YetInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_No_Usage_YetInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_no_usage_yet(inputs)
	return __ro.credit_usage_no_usage_yet(inputs)
});
/**
* | output |
* | --- |
* | "No credit usage history matches the selected filters." |
*
* @param {Credit_Usage_No_Usage_MatchInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_no_usage_match = /** @type {((inputs?: Credit_Usage_No_Usage_MatchInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_No_Usage_MatchInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_no_usage_match(inputs)
	return __ro.credit_usage_no_usage_match(inputs)
});
/**
* | output |
* | --- |
* | "Purchases and debits will appear here after they settle." |
*
* @param {Credit_Usage_Empty_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_empty_body = /** @type {((inputs?: Credit_Usage_Empty_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Empty_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_empty_body(inputs)
	return __ro.credit_usage_empty_body(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} record on this page." |
*
* @param {Credit_Usage_Showing_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_showing_one = /** @type {((inputs: Credit_Usage_Showing_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Showing_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_showing_one(inputs)
	return __ro.credit_usage_showing_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} records on this page." |
*
* @param {Credit_Usage_Showing_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_showing_other = /** @type {((inputs: Credit_Usage_Showing_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Showing_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_showing_other(inputs)
	return __ro.credit_usage_showing_other(inputs)
});
/**
* | output |
* | --- |
* | "No credit usage history to show." |
*
* @param {Credit_Usage_None_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_none_to_show = /** @type {((inputs?: Credit_Usage_None_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_None_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_none_to_show(inputs)
	return __ro.credit_usage_none_to_show(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date ascending" |
*
* @param {Credit_Usage_Sort_Created_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_sort_created_ascending = /** @type {((inputs?: Credit_Usage_Sort_Created_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Sort_Created_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_sort_created_ascending(inputs)
	return __ro.credit_usage_sort_created_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date descending" |
*
* @param {Credit_Usage_Sort_Created_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const credit_usage_sort_created_descending = /** @type {((inputs?: Credit_Usage_Sort_Created_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Credit_Usage_Sort_Created_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.credit_usage_sort_created_descending(inputs)
	return __ro.credit_usage_sort_created_descending(inputs)
});
/**
* | output |
* | --- |
* | "Account Settings" |
*
* @param {Account_Settings_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_title = /** @type {((inputs?: Account_Settings_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_title(inputs)
	return __ro.account_settings_title(inputs)
});
/**
* | output |
* | --- |
* | "Manage your account details and security." |
*
* @param {Account_Settings_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_description = /** @type {((inputs?: Account_Settings_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_description(inputs)
	return __ro.account_settings_description(inputs)
});
/**
* | output |
* | --- |
* | "Account settings" |
*
* @param {Account_Settings_Nav_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_nav_label = /** @type {((inputs?: Account_Settings_Nav_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Nav_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_nav_label(inputs)
	return __ro.account_settings_nav_label(inputs)
});
/**
* | output |
* | --- |
* | "Account" |
*
* @param {Account_Settings_Account_FallbackInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_account_fallback = /** @type {((inputs?: Account_Settings_Account_FallbackInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Account_FallbackInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_account_fallback(inputs)
	return __ro.account_settings_account_fallback(inputs)
});
/**
* | output |
* | --- |
* | "No email address" |
*
* @param {Account_Settings_No_Email_AddressInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_no_email_address = /** @type {((inputs?: Account_Settings_No_Email_AddressInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_No_Email_AddressInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_no_email_address(inputs)
	return __ro.account_settings_no_email_address(inputs)
});
/**
* | output |
* | --- |
* | "General" |
*
* @param {Account_Settings_GeneralInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_general = /** @type {((inputs?: Account_Settings_GeneralInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_GeneralInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_general(inputs)
	return __ro.account_settings_general(inputs)
});
/**
* | output |
* | --- |
* | "Security" |
*
* @param {Account_Settings_SecurityInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_security = /** @type {((inputs?: Account_Settings_SecurityInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_SecurityInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_security(inputs)
	return __ro.account_settings_security(inputs)
});
/**
* | output |
* | --- |
* | "Sessions" |
*
* @param {Account_Settings_SessionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_sessions = /** @type {((inputs?: Account_Settings_SessionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_SessionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_sessions(inputs)
	return __ro.account_settings_sessions(inputs)
});
/**
* | output |
* | --- |
* | "Linked Accounts" |
*
* @param {Account_Settings_Linked_AccountsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_linked_accounts = /** @type {((inputs?: Account_Settings_Linked_AccountsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Linked_AccountsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_linked_accounts(inputs)
	return __ro.account_settings_linked_accounts(inputs)
});
/**
* | output |
* | --- |
* | "Unable to update account settings." |
*
* @param {Account_Settings_Update_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_update_error = /** @type {((inputs?: Account_Settings_Update_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Update_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_update_error(inputs)
	return __ro.account_settings_update_error(inputs)
});
/**
* | output |
* | --- |
* | "Unable to save changes." |
*
* @param {Account_Settings_Save_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_save_error = /** @type {((inputs?: Account_Settings_Save_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Save_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_save_error(inputs)
	return __ro.account_settings_save_error(inputs)
});
/**
* | output |
* | --- |
* | "Revoke session?" |
*
* @param {Account_Settings_Revoke_Session_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_revoke_session_title = /** @type {((inputs?: Account_Settings_Revoke_Session_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Revoke_Session_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_revoke_session_title(inputs)
	return __ro.account_settings_revoke_session_title(inputs)
});
/**
* | output |
* | --- |
* | "Revoke {session}. That device will need to sign in again." |
*
* @param {Account_Settings_Revoke_Session_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_revoke_session_description = /** @type {((inputs: Account_Settings_Revoke_Session_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Revoke_Session_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_revoke_session_description(inputs)
	return __ro.account_settings_revoke_session_description(inputs)
});
/**
* | output |
* | --- |
* | "Revoke" |
*
* @param {Account_Settings_RevokeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_revoke = /** @type {((inputs?: Account_Settings_RevokeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_RevokeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_revoke(inputs)
	return __ro.account_settings_revoke(inputs)
});
/**
* | output |
* | --- |
* | "Session revoked." |
*
* @param {Account_Settings_Session_RevokedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_session_revoked = /** @type {((inputs?: Account_Settings_Session_RevokedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Session_RevokedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_session_revoked(inputs)
	return __ro.account_settings_session_revoked(inputs)
});
/**
* | output |
* | --- |
* | "Unlink {provider}?" |
*
* @param {Account_Settings_Unlink_Provider_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unlink_provider_title = /** @type {((inputs: Account_Settings_Unlink_Provider_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Unlink_Provider_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unlink_provider_title(inputs)
	return __ro.account_settings_unlink_provider_title(inputs)
});
/**
* | output |
* | --- |
* | "Remove {provider} sign-in from this account. You can reconnect it later." |
*
* @param {Account_Settings_Unlink_Provider_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unlink_provider_description = /** @type {((inputs: Account_Settings_Unlink_Provider_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Unlink_Provider_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unlink_provider_description(inputs)
	return __ro.account_settings_unlink_provider_description(inputs)
});
/**
* | output |
* | --- |
* | "Unlink" |
*
* @param {Account_Settings_UnlinkInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unlink = /** @type {((inputs?: Account_Settings_UnlinkInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_UnlinkInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unlink(inputs)
	return __ro.account_settings_unlink(inputs)
});
/**
* | output |
* | --- |
* | "Linked account removed." |
*
* @param {Account_Settings_Linked_Account_RemovedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_linked_account_removed = /** @type {((inputs?: Account_Settings_Linked_Account_RemovedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Linked_Account_RemovedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_linked_account_removed(inputs)
	return __ro.account_settings_linked_account_removed(inputs)
});
/**
* | output |
* | --- |
* | "Avatar updated." |
*
* @param {Account_Settings_Avatar_SavedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar_saved = /** @type {((inputs?: Account_Settings_Avatar_SavedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Avatar_SavedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar_saved(inputs)
	return __ro.account_settings_avatar_saved(inputs)
});
/**
* | output |
* | --- |
* | "Display name saved." |
*
* @param {Account_Settings_Name_SavedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_name_saved = /** @type {((inputs?: Account_Settings_Name_SavedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Name_SavedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_name_saved(inputs)
	return __ro.account_settings_name_saved(inputs)
});
/**
* | output |
* | --- |
* | "Email address saved." |
*
* @param {Account_Settings_Email_SavedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_email_saved = /** @type {((inputs?: Account_Settings_Email_SavedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Email_SavedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_email_saved(inputs)
	return __ro.account_settings_email_saved(inputs)
});
/**
* | output |
* | --- |
* | "Language preference saved." |
*
* @param {Account_Settings_Language_SavedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_language_saved = /** @type {((inputs?: Account_Settings_Language_SavedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Language_SavedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_language_saved(inputs)
	return __ro.account_settings_language_saved(inputs)
});
/**
* | output |
* | --- |
* | "Password updated." |
*
* @param {Account_Settings_Password_UpdatedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_password_updated = /** @type {((inputs?: Account_Settings_Password_UpdatedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Password_UpdatedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_password_updated(inputs)
	return __ro.account_settings_password_updated(inputs)
});
/**
* | output |
* | --- |
* | "current session" |
*
* @param {Account_Settings_Current_SessionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_current_session = /** @type {((inputs?: Account_Settings_Current_SessionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Current_SessionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_current_session(inputs)
	return __ro.account_settings_current_session(inputs)
});
/**
* | output |
* | --- |
* | "browser session" |
*
* @param {Account_Settings_Browser_SessionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_browser_session = /** @type {((inputs?: Account_Settings_Browser_SessionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Browser_SessionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_browser_session(inputs)
	return __ro.account_settings_browser_session(inputs)
});
/**
* | output |
* | --- |
* | "Created {date}" |
*
* @param {Account_Settings_Session_Created_AtInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_session_created_at = /** @type {((inputs: Account_Settings_Session_Created_AtInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Session_Created_AtInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_session_created_at(inputs)
	return __ro.account_settings_session_created_at(inputs)
});
/**
* | output |
* | --- |
* | "{ip} - created {date}" |
*
* @param {Account_Settings_Session_Ip_Created_AtInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_session_ip_created_at = /** @type {((inputs: Account_Settings_Session_Ip_Created_AtInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Session_Ip_Created_AtInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_session_ip_created_at(inputs)
	return __ro.account_settings_session_ip_created_at(inputs)
});
/**
* | output |
* | --- |
* | "Unknown" |
*
* @param {Account_Settings_UnknownInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unknown = /** @type {((inputs?: Account_Settings_UnknownInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_UnknownInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unknown(inputs)
	return __ro.account_settings_unknown(inputs)
});
/**
* | output |
* | --- |
* | "Avatar" |
*
* @param {Account_Settings_AvatarInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar = /** @type {((inputs?: Account_Settings_AvatarInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_AvatarInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar(inputs)
	return __ro.account_settings_avatar(inputs)
});
/**
* | output |
* | --- |
* | "Upload a profile image shown across your account." |
*
* @param {Account_Settings_Avatar_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar_description = /** @type {((inputs?: Account_Settings_Avatar_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Avatar_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar_description(inputs)
	return __ro.account_settings_avatar_description(inputs)
});
/**
* | output |
* | --- |
* | "Uploading..." |
*
* @param {Account_Settings_Avatar_UploadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar_uploading = /** @type {((inputs?: Account_Settings_Avatar_UploadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Avatar_UploadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar_uploading(inputs)
	return __ro.account_settings_avatar_uploading(inputs)
});
/**
* | output |
* | --- |
* | "Click to upload and crop" |
*
* @param {Account_Settings_Avatar_UploadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar_upload = /** @type {((inputs?: Account_Settings_Avatar_UploadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Avatar_UploadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar_upload(inputs)
	return __ro.account_settings_avatar_upload(inputs)
});
/**
* | output |
* | --- |
* | "PNG, JPG, GIF, AVIF, APNG, SVG, WEBP up to 5 MB." |
*
* @param {Account_Settings_Avatar_File_HintInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_avatar_file_hint = /** @type {((inputs?: Account_Settings_Avatar_File_HintInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Avatar_File_HintInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_avatar_file_hint(inputs)
	return __ro.account_settings_avatar_file_hint(inputs)
});
/**
* | output |
* | --- |
* | "Crop Avatar" |
*
* @param {Account_Settings_Crop_AvatarInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_crop_avatar = /** @type {((inputs?: Account_Settings_Crop_AvatarInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Crop_AvatarInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_crop_avatar(inputs)
	return __ro.account_settings_crop_avatar(inputs)
});
/**
* | output |
* | --- |
* | "Adjust your avatar's crop area before saving." |
*
* @param {Account_Settings_Crop_Avatar_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_crop_avatar_description = /** @type {((inputs?: Account_Settings_Crop_Avatar_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Crop_Avatar_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_crop_avatar_description(inputs)
	return __ro.account_settings_crop_avatar_description(inputs)
});
/**
* | output |
* | --- |
* | "Display Name" |
*
* @param {Account_Settings_Display_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_display_name = /** @type {((inputs?: Account_Settings_Display_NameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Display_NameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_display_name(inputs)
	return __ro.account_settings_display_name(inputs)
});
/**
* | output |
* | --- |
* | "Email Address" |
*
* @param {Account_Settings_Email_AddressInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_email_address = /** @type {((inputs?: Account_Settings_Email_AddressInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Email_AddressInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_email_address(inputs)
	return __ro.account_settings_email_address(inputs)
});
/**
* | output |
* | --- |
* | "Language" |
*
* @param {Account_Settings_LanguageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_language = /** @type {((inputs?: Account_Settings_LanguageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_LanguageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_language(inputs)
	return __ro.account_settings_language(inputs)
});
/**
* | output |
* | --- |
* | "Save name" |
*
* @param {Account_Settings_Save_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_save_name = /** @type {((inputs?: Account_Settings_Save_NameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Save_NameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_save_name(inputs)
	return __ro.account_settings_save_name(inputs)
});
/**
* | output |
* | --- |
* | "Save email" |
*
* @param {Account_Settings_Save_EmailInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_save_email = /** @type {((inputs?: Account_Settings_Save_EmailInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Save_EmailInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_save_email(inputs)
	return __ro.account_settings_save_email(inputs)
});
/**
* | output |
* | --- |
* | "Save language" |
*
* @param {Account_Settings_Save_LanguageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_save_language = /** @type {((inputs?: Account_Settings_Save_LanguageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Save_LanguageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_save_language(inputs)
	return __ro.account_settings_save_language(inputs)
});
/**
* | output |
* | --- |
* | "Save password" |
*
* @param {Account_Settings_Save_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_save_password = /** @type {((inputs?: Account_Settings_Save_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Save_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_save_password(inputs)
	return __ro.account_settings_save_password(inputs)
});
/**
* | output |
* | --- |
* | "New Password" |
*
* @param {Account_Settings_New_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_new_password = /** @type {((inputs?: Account_Settings_New_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_New_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_new_password(inputs)
	return __ro.account_settings_new_password(inputs)
});
/**
* | output |
* | --- |
* | "Confirm Password" |
*
* @param {Account_Settings_Confirm_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_confirm_password = /** @type {((inputs?: Account_Settings_Confirm_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Confirm_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_confirm_password(inputs)
	return __ro.account_settings_confirm_password(inputs)
});
/**
* | output |
* | --- |
* | "Change your account password." |
*
* @param {Account_Settings_Security_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_security_description = /** @type {((inputs?: Account_Settings_Security_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Security_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_security_description(inputs)
	return __ro.account_settings_security_description(inputs)
});
/**
* | output |
* | --- |
* | "Review browsers and devices signed in to this account." |
*
* @param {Account_Settings_Sessions_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_sessions_description = /** @type {((inputs?: Account_Settings_Sessions_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Sessions_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_sessions_description(inputs)
	return __ro.account_settings_sessions_description(inputs)
});
/**
* | output |
* | --- |
* | "Loading sessions..." |
*
* @param {Account_Settings_Loading_SessionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_loading_sessions = /** @type {((inputs?: Account_Settings_Loading_SessionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Loading_SessionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_loading_sessions(inputs)
	return __ro.account_settings_loading_sessions(inputs)
});
/**
* | output |
* | --- |
* | "No active sessions found." |
*
* @param {Account_Settings_No_SessionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_no_sessions = /** @type {((inputs?: Account_Settings_No_SessionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_No_SessionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_no_sessions(inputs)
	return __ro.account_settings_no_sessions(inputs)
});
/**
* | output |
* | --- |
* | "Current" |
*
* @param {Account_Settings_CurrentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_current = /** @type {((inputs?: Account_Settings_CurrentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_CurrentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_current(inputs)
	return __ro.account_settings_current(inputs)
});
/**
* | output |
* | --- |
* | "Expires {date}" |
*
* @param {Account_Settings_ExpiresInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_expires = /** @type {((inputs: Account_Settings_ExpiresInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_ExpiresInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_expires(inputs)
	return __ro.account_settings_expires(inputs)
});
/**
* | output |
* | --- |
* | "Current session cannot be revoked" |
*
* @param {Account_Settings_Current_Session_Cannot_RevokeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_current_session_cannot_revoke = /** @type {((inputs?: Account_Settings_Current_Session_Cannot_RevokeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Current_Session_Cannot_RevokeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_current_session_cannot_revoke(inputs)
	return __ro.account_settings_current_session_cannot_revoke(inputs)
});
/**
* | output |
* | --- |
* | "Revoke session" |
*
* @param {Account_Settings_Revoke_SessionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_revoke_session = /** @type {((inputs?: Account_Settings_Revoke_SessionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Revoke_SessionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_revoke_session(inputs)
	return __ro.account_settings_revoke_session(inputs)
});
/**
* | output |
* | --- |
* | "Revoking..." |
*
* @param {Account_Settings_RevokingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_revoking = /** @type {((inputs?: Account_Settings_RevokingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_RevokingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_revoking(inputs)
	return __ro.account_settings_revoking(inputs)
});
/**
* | output |
* | --- |
* | "Manage sign-in methods connected to this account." |
*
* @param {Account_Settings_Linked_Accounts_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_linked_accounts_description = /** @type {((inputs?: Account_Settings_Linked_Accounts_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Linked_Accounts_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_linked_accounts_description(inputs)
	return __ro.account_settings_linked_accounts_description(inputs)
});
/**
* | output |
* | --- |
* | "Loading linked accounts..." |
*
* @param {Account_Settings_Loading_Linked_AccountsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_loading_linked_accounts = /** @type {((inputs?: Account_Settings_Loading_Linked_AccountsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Loading_Linked_AccountsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_loading_linked_accounts(inputs)
	return __ro.account_settings_loading_linked_accounts(inputs)
});
/**
* | output |
* | --- |
* | "No sign-in methods were returned." |
*
* @param {Account_Settings_No_Sign_In_MethodsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_no_sign_in_methods = /** @type {((inputs?: Account_Settings_No_Sign_In_MethodsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_No_Sign_In_MethodsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_no_sign_in_methods(inputs)
	return __ro.account_settings_no_sign_in_methods(inputs)
});
/**
* | output |
* | --- |
* | "Email/password" |
*
* @param {Account_Settings_Email_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_email_password = /** @type {((inputs?: Account_Settings_Email_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Email_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_email_password(inputs)
	return __ro.account_settings_email_password(inputs)
});
/**
* | output |
* | --- |
* | "Password sign-in is enabled for {email}." |
*
* @param {Account_Settings_Password_EnabledInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_password_enabled = /** @type {((inputs: Account_Settings_Password_EnabledInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Password_EnabledInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_password_enabled(inputs)
	return __ro.account_settings_password_enabled(inputs)
});
/**
* | output |
* | --- |
* | "Add a password to use email sign-in." |
*
* @param {Account_Settings_Add_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_add_password = /** @type {((inputs?: Account_Settings_Add_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Add_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_add_password(inputs)
	return __ro.account_settings_add_password(inputs)
});
/**
* | output |
* | --- |
* | "Set password" |
*
* @param {Account_Settings_Set_PasswordInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_set_password = /** @type {((inputs?: Account_Settings_Set_PasswordInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Set_PasswordInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_set_password(inputs)
	return __ro.account_settings_set_password(inputs)
});
/**
* | output |
* | --- |
* | "Use your Google account to sign in." |
*
* @param {Account_Settings_Provider_Google_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_provider_google_description = /** @type {((inputs?: Account_Settings_Provider_Google_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Provider_Google_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_provider_google_description(inputs)
	return __ro.account_settings_provider_google_description(inputs)
});
/**
* | output |
* | --- |
* | "Use your GitHub account to sign in." |
*
* @param {Account_Settings_Provider_Github_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_provider_github_description = /** @type {((inputs?: Account_Settings_Provider_Github_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Provider_Github_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_provider_github_description(inputs)
	return __ro.account_settings_provider_github_description(inputs)
});
/**
* | output |
* | --- |
* | "Linked {date}" |
*
* @param {Account_Settings_Linked_AtInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_linked_at = /** @type {((inputs: Account_Settings_Linked_AtInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Linked_AtInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_linked_at(inputs)
	return __ro.account_settings_linked_at(inputs)
});
/**
* | output |
* | --- |
* | "Unlinking..." |
*
* @param {Account_Settings_UnlinkingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unlinking = /** @type {((inputs?: Account_Settings_UnlinkingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_UnlinkingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unlinking(inputs)
	return __ro.account_settings_unlinking(inputs)
});
/**
* | output |
* | --- |
* | "Unavailable" |
*
* @param {Account_Settings_Unavailable_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unavailable_title = /** @type {((inputs?: Account_Settings_Unavailable_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Unavailable_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unavailable_title(inputs)
	return __ro.account_settings_unavailable_title(inputs)
});
/**
* | output |
* | --- |
* | "This section is not available yet." |
*
* @param {Account_Settings_Unavailable_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const account_settings_unavailable_body = /** @type {((inputs?: Account_Settings_Unavailable_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Account_Settings_Unavailable_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.account_settings_unavailable_body(inputs)
	return __ro.account_settings_unavailable_body(inputs)
});
/**
* | output |
* | --- |
* | "Billing Information" |
*
* @param {Billing_Profile_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_title = /** @type {((inputs?: Billing_Profile_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_title(inputs)
	return __ro.billing_profile_title(inputs)
});
/**
* | output |
* | --- |
* | "Manage the billing details used for invoices." |
*
* @param {Billing_Profile_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_description = /** @type {((inputs?: Billing_Profile_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_description(inputs)
	return __ro.billing_profile_description(inputs)
});
/**
* | output |
* | --- |
* | "Unable to load billing information." |
*
* @param {Billing_Profile_Load_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_load_error = /** @type {((inputs?: Billing_Profile_Load_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Load_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_load_error(inputs)
	return __ro.billing_profile_load_error(inputs)
});
/**
* | output |
* | --- |
* | "Unable to save billing information." |
*
* @param {Billing_Profile_Save_ErrorInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_save_error = /** @type {((inputs?: Billing_Profile_Save_ErrorInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Save_ErrorInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_save_error(inputs)
	return __ro.billing_profile_save_error(inputs)
});
/**
* | output |
* | --- |
* | "Billing information saved." |
*
* @param {Billing_Profile_SavedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_saved = /** @type {((inputs?: Billing_Profile_SavedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_SavedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_saved(inputs)
	return __ro.billing_profile_saved(inputs)
});
/**
* | output |
* | --- |
* | "Company name" |
*
* @param {Billing_Profile_Company_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_company_name = /** @type {((inputs?: Billing_Profile_Company_NameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Company_NameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_company_name(inputs)
	return __ro.billing_profile_company_name(inputs)
});
/**
* | output |
* | --- |
* | "Full name" |
*
* @param {Billing_Profile_Full_NameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_full_name = /** @type {((inputs?: Billing_Profile_Full_NameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Full_NameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_full_name(inputs)
	return __ro.billing_profile_full_name(inputs)
});
/**
* | output |
* | --- |
* | "An error occurred" |
*
* @param {Billing_Profile_Error_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_error_title = /** @type {((inputs?: Billing_Profile_Error_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Error_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_error_title(inputs)
	return __ro.billing_profile_error_title(inputs)
});
/**
* | output |
* | --- |
* | "Loading billing information..." |
*
* @param {Billing_Profile_LoadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_loading = /** @type {((inputs?: Billing_Profile_LoadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_LoadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_loading(inputs)
	return __ro.billing_profile_loading(inputs)
});
/**
* | output |
* | --- |
* | "Please wait while we retrieve your profile details." |
*
* @param {Billing_Profile_Loading_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_loading_body = /** @type {((inputs?: Billing_Profile_Loading_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Loading_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_loading_body(inputs)
	return __ro.billing_profile_loading_body(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load billing information" |
*
* @param {Billing_Profile_Failed_LoadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_failed_load = /** @type {((inputs?: Billing_Profile_Failed_LoadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Failed_LoadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_failed_load(inputs)
	return __ro.billing_profile_failed_load(inputs)
});
/**
* | output |
* | --- |
* | "Retry Loading" |
*
* @param {Billing_Profile_Retry_LoadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_retry_loading = /** @type {((inputs?: Billing_Profile_Retry_LoadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Retry_LoadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_retry_loading(inputs)
	return __ro.billing_profile_retry_loading(inputs)
});
/**
* | output |
* | --- |
* | "Billing Entity" |
*
* @param {Billing_Profile_Billing_EntityInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_billing_entity = /** @type {((inputs?: Billing_Profile_Billing_EntityInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Billing_EntityInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_billing_entity(inputs)
	return __ro.billing_profile_billing_entity(inputs)
});
/**
* | output |
* | --- |
* | "Choose between an individual or business profile." |
*
* @param {Billing_Profile_Entity_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_entity_description = /** @type {((inputs?: Billing_Profile_Entity_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Entity_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_entity_description(inputs)
	return __ro.billing_profile_entity_description(inputs)
});
/**
* | output |
* | --- |
* | "Individual" |
*
* @param {Billing_Profile_IndividualInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_individual = /** @type {((inputs?: Billing_Profile_IndividualInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_IndividualInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_individual(inputs)
	return __ro.billing_profile_individual(inputs)
});
/**
* | output |
* | --- |
* | "Company" |
*
* @param {Billing_Profile_CompanyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_company = /** @type {((inputs?: Billing_Profile_CompanyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_CompanyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_company(inputs)
	return __ro.billing_profile_company(inputs)
});
/**
* | output |
* | --- |
* | "General Details" |
*
* @param {Billing_Profile_General_DetailsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_general_details = /** @type {((inputs?: Billing_Profile_General_DetailsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_General_DetailsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_general_details(inputs)
	return __ro.billing_profile_general_details(inputs)
});
/**
* | output |
* | --- |
* | "Billing email" |
*
* @param {Billing_Profile_Billing_EmailInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_billing_email = /** @type {((inputs?: Billing_Profile_Billing_EmailInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Billing_EmailInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_billing_email(inputs)
	return __ro.billing_profile_billing_email(inputs)
});
/**
* | output |
* | --- |
* | "Billing Address" |
*
* @param {Billing_Profile_Billing_AddressInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_billing_address = /** @type {((inputs?: Billing_Profile_Billing_AddressInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Billing_AddressInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_billing_address(inputs)
	return __ro.billing_profile_billing_address(inputs)
});
/**
* | output |
* | --- |
* | "Address line 1" |
*
* @param {Billing_Profile_Address_Line1Inputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_address_line1 = /** @type {((inputs?: Billing_Profile_Address_Line1Inputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Address_Line1Inputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_address_line1(inputs)
	return __ro.billing_profile_address_line1(inputs)
});
/**
* | output |
* | --- |
* | "Address line 2" |
*
* @param {Billing_Profile_Address_Line2Inputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_address_line2 = /** @type {((inputs?: Billing_Profile_Address_Line2Inputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Address_Line2Inputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_address_line2(inputs)
	return __ro.billing_profile_address_line2(inputs)
});
/**
* | output |
* | --- |
* | "City" |
*
* @param {Billing_Profile_CityInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_city = /** @type {((inputs?: Billing_Profile_CityInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_CityInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_city(inputs)
	return __ro.billing_profile_city(inputs)
});
/**
* | output |
* | --- |
* | "Region/state" |
*
* @param {Billing_Profile_Region_StateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_region_state = /** @type {((inputs?: Billing_Profile_Region_StateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Region_StateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_region_state(inputs)
	return __ro.billing_profile_region_state(inputs)
});
/**
* | output |
* | --- |
* | "Country" |
*
* @param {Billing_Profile_CountryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_country = /** @type {((inputs?: Billing_Profile_CountryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_CountryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_country(inputs)
	return __ro.billing_profile_country(inputs)
});
/**
* | output |
* | --- |
* | "Postal code" |
*
* @param {Billing_Profile_Postal_CodeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_postal_code = /** @type {((inputs?: Billing_Profile_Postal_CodeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Postal_CodeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_postal_code(inputs)
	return __ro.billing_profile_postal_code(inputs)
});
/**
* | output |
* | --- |
* | "Company Details" |
*
* @param {Billing_Profile_Company_DetailsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_company_details = /** @type {((inputs?: Billing_Profile_Company_DetailsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Company_DetailsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_company_details(inputs)
	return __ro.billing_profile_company_details(inputs)
});
/**
* | output |
* | --- |
* | "Fiscal code" |
*
* @param {Billing_Profile_Fiscal_CodeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_fiscal_code = /** @type {((inputs?: Billing_Profile_Fiscal_CodeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Fiscal_CodeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_fiscal_code(inputs)
	return __ro.billing_profile_fiscal_code(inputs)
});
/**
* | output |
* | --- |
* | "Registration number" |
*
* @param {Billing_Profile_Registration_NumberInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_registration_number = /** @type {((inputs?: Billing_Profile_Registration_NumberInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Registration_NumberInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_registration_number(inputs)
	return __ro.billing_profile_registration_number(inputs)
});
/**
* | output |
* | --- |
* | "Save billing information" |
*
* @param {Billing_Profile_Save_ButtonInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const billing_profile_save_button = /** @type {((inputs?: Billing_Profile_Save_ButtonInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Billing_Profile_Save_ButtonInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.billing_profile_save_button(inputs)
	return __ro.billing_profile_save_button(inputs)
});
/**
* | output |
* | --- |
* | "Datasets" |
*
* @param {Datasets_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_page_title = /** @type {((inputs?: Datasets_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_page_title(inputs)
	return __ro.datasets_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Dataset" |
*
* @param {Datasets_Detail_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_detail_page_title = /** @type {((inputs?: Datasets_Detail_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Detail_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_detail_page_title(inputs)
	return __ro.datasets_detail_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Datasets_Name_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_name_column = /** @type {((inputs?: Datasets_Name_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Name_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_name_column(inputs)
	return __ro.datasets_name_column(inputs)
});
/**
* | output |
* | --- |
* | "Schema" |
*
* @param {Datasets_Schema_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_schema_column = /** @type {((inputs?: Datasets_Schema_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Schema_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_schema_column(inputs)
	return __ro.datasets_schema_column(inputs)
});
/**
* | output |
* | --- |
* | "Fields" |
*
* @param {Datasets_Fields_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_fields_column = /** @type {((inputs?: Datasets_Fields_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Fields_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_fields_column(inputs)
	return __ro.datasets_fields_column(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Datasets_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_created_column = /** @type {((inputs?: Datasets_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_created_column(inputs)
	return __ro.datasets_created_column(inputs)
});
/**
* | output |
* | --- |
* | "Actions" |
*
* @param {Datasets_Actions_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_actions_column = /** @type {((inputs?: Datasets_Actions_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Actions_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_actions_column(inputs)
	return __ro.datasets_actions_column(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date ascending" |
*
* @param {Datasets_Sort_Created_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_sort_created_ascending = /** @type {((inputs?: Datasets_Sort_Created_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Sort_Created_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_sort_created_ascending(inputs)
	return __ro.datasets_sort_created_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date descending" |
*
* @param {Datasets_Sort_Created_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_sort_created_descending = /** @type {((inputs?: Datasets_Sort_Created_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Sort_Created_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_sort_created_descending(inputs)
	return __ro.datasets_sort_created_descending(inputs)
});
/**
* | output |
* | --- |
* | "Retry" |
*
* @param {Datasets_RetryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_retry = /** @type {((inputs?: Datasets_RetryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_RetryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_retry(inputs)
	return __ro.datasets_retry(inputs)
});
/**
* | output |
* | --- |
* | "Open" |
*
* @param {Datasets_OpenInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_open = /** @type {((inputs?: Datasets_OpenInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_OpenInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_open(inputs)
	return __ro.datasets_open(inputs)
});
/**
* | output |
* | --- |
* | "No datasets found" |
*
* @param {Datasets_No_Datasets_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_datasets_found = /** @type {((inputs?: Datasets_No_Datasets_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Datasets_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_datasets_found(inputs)
	return __ro.datasets_no_datasets_found(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} dataset on this page." |
*
* @param {Datasets_Showing_Datasets_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_showing_datasets_one = /** @type {((inputs: Datasets_Showing_Datasets_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Showing_Datasets_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_showing_datasets_one(inputs)
	return __ro.datasets_showing_datasets_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} datasets on this page." |
*
* @param {Datasets_Showing_Datasets_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_showing_datasets_other = /** @type {((inputs: Datasets_Showing_Datasets_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Showing_Datasets_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_showing_datasets_other(inputs)
	return __ro.datasets_showing_datasets_other(inputs)
});
/**
* | output |
* | --- |
* | "No datasets to show." |
*
* @param {Datasets_No_Datasets_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_datasets_to_show = /** @type {((inputs?: Datasets_No_Datasets_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Datasets_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_datasets_to_show(inputs)
	return __ro.datasets_no_datasets_to_show(inputs)
});
/**
* | output |
* | --- |
* | "Rows per page" |
*
* @param {Datasets_Rows_Per_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_rows_per_page = /** @type {((inputs?: Datasets_Rows_Per_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Rows_Per_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_rows_per_page(inputs)
	return __ro.datasets_rows_per_page(inputs)
});
/**
* | output |
* | --- |
* | "Previous page" |
*
* @param {Datasets_Previous_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_previous_page = /** @type {((inputs?: Datasets_Previous_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Previous_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_previous_page(inputs)
	return __ro.datasets_previous_page(inputs)
});
/**
* | output |
* | --- |
* | "Next page" |
*
* @param {Datasets_Next_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_next_page = /** @type {((inputs?: Datasets_Next_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Next_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_next_page(inputs)
	return __ro.datasets_next_page(inputs)
});
/**
* | output |
* | --- |
* | "{count} field" |
*
* @param {Datasets_Field_Count_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_field_count_one = /** @type {((inputs: Datasets_Field_Count_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Field_Count_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_field_count_one(inputs)
	return __ro.datasets_field_count_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} fields" |
*
* @param {Datasets_Field_Count_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_field_count_other = /** @type {((inputs: Datasets_Field_Count_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Field_Count_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_field_count_other(inputs)
	return __ro.datasets_field_count_other(inputs)
});
/**
* | output |
* | --- |
* | "Date range" |
*
* @param {Datasets_Date_RangeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_date_range = /** @type {((inputs?: Datasets_Date_RangeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Date_RangeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_date_range(inputs)
	return __ro.datasets_date_range(inputs)
});
/**
* | output |
* | --- |
* | "Any" |
*
* @param {Datasets_Any_DateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_any_date = /** @type {((inputs?: Datasets_Any_DateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Any_DateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_any_date(inputs)
	return __ro.datasets_any_date(inputs)
});
/**
* | output |
* | --- |
* | "{start} - {end}" |
*
* @param {Datasets_Date_Range_ValueInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_date_range_value = /** @type {((inputs: Datasets_Date_Range_ValueInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Date_Range_ValueInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_date_range_value(inputs)
	return __ro.datasets_date_range_value(inputs)
});
/**
* | output |
* | --- |
* | "Presets" |
*
* @param {Datasets_PresetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_presets = /** @type {((inputs?: Datasets_PresetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_PresetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_presets(inputs)
	return __ro.datasets_presets(inputs)
});
/**
* | output |
* | --- |
* | "Today" |
*
* @param {Datasets_TodayInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_today = /** @type {((inputs?: Datasets_TodayInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_TodayInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_today(inputs)
	return __ro.datasets_today(inputs)
});
/**
* | output |
* | --- |
* | "This week" |
*
* @param {Datasets_This_WeekInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_this_week = /** @type {((inputs?: Datasets_This_WeekInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_This_WeekInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_this_week(inputs)
	return __ro.datasets_this_week(inputs)
});
/**
* | output |
* | --- |
* | "This month" |
*
* @param {Datasets_This_MonthInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_this_month = /** @type {((inputs?: Datasets_This_MonthInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_This_MonthInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_this_month(inputs)
	return __ro.datasets_this_month(inputs)
});
/**
* | output |
* | --- |
* | "Clear" |
*
* @param {Datasets_ClearInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_clear = /** @type {((inputs?: Datasets_ClearInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_ClearInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_clear(inputs)
	return __ro.datasets_clear(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Datasets_ApplyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_apply = /** @type {((inputs?: Datasets_ApplyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_ApplyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_apply(inputs)
	return __ro.datasets_apply(inputs)
});
/**
* | output |
* | --- |
* | "Document ID" |
*
* @param {Datasets_Document_Id_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_document_id_column = /** @type {((inputs?: Datasets_Document_Id_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Document_Id_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_document_id_column(inputs)
	return __ro.datasets_document_id_column(inputs)
});
/**
* | output |
* | --- |
* | "Filename" |
*
* @param {Datasets_Filename_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_filename_column = /** @type {((inputs?: Datasets_Filename_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Filename_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_filename_column(inputs)
	return __ro.datasets_filename_column(inputs)
});
/**
* | output |
* | --- |
* | "Dataset not found" |
*
* @param {Datasets_Not_Found_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_not_found_title = /** @type {((inputs?: Datasets_Not_Found_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Not_Found_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_not_found_title(inputs)
	return __ro.datasets_not_found_title(inputs)
});
/**
* | output |
* | --- |
* | "This dataset does not exist." |
*
* @param {Datasets_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_not_found_body = /** @type {((inputs?: Datasets_Not_Found_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Not_Found_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_not_found_body(inputs)
	return __ro.datasets_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "View datasets" |
*
* @param {Datasets_View_DatasetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_view_datasets = /** @type {((inputs?: Datasets_View_DatasetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_View_DatasetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_view_datasets(inputs)
	return __ro.datasets_view_datasets(inputs)
});
/**
* | output |
* | --- |
* | "Preview document {documentId}" |
*
* @param {Datasets_Preview_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_preview_document = /** @type {((inputs: Datasets_Preview_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Preview_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_preview_document(inputs)
	return __ro.datasets_preview_document(inputs)
});
/**
* | output |
* | --- |
* | "No documents extracted for this dataset" |
*
* @param {Datasets_No_Documents_ExtractedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_documents_extracted = /** @type {((inputs?: Datasets_No_Documents_ExtractedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Documents_ExtractedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_documents_extracted(inputs)
	return __ro.datasets_no_documents_extracted(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} row on this page." |
*
* @param {Datasets_Showing_Rows_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_showing_rows_one = /** @type {((inputs: Datasets_Showing_Rows_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Showing_Rows_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_showing_rows_one(inputs)
	return __ro.datasets_showing_rows_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} rows on this page." |
*
* @param {Datasets_Showing_Rows_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_showing_rows_other = /** @type {((inputs: Datasets_Showing_Rows_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Showing_Rows_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_showing_rows_other(inputs)
	return __ro.datasets_showing_rows_other(inputs)
});
/**
* | output |
* | --- |
* | "No rows to show." |
*
* @param {Datasets_No_Rows_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_rows_to_show = /** @type {((inputs?: Datasets_No_Rows_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Rows_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_rows_to_show(inputs)
	return __ro.datasets_no_rows_to_show(inputs)
});
/**
* | output |
* | --- |
* | "CSV export" |
*
* @param {Datasets_Export_CsvInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_export_csv = /** @type {((inputs?: Datasets_Export_CsvInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Export_CsvInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_export_csv(inputs)
	return __ro.datasets_export_csv(inputs)
});
/**
* | output |
* | --- |
* | "XLSX export" |
*
* @param {Datasets_Export_XlsxInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_export_xlsx = /** @type {((inputs?: Datasets_Export_XlsxInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Export_XlsxInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_export_xlsx(inputs)
	return __ro.datasets_export_xlsx(inputs)
});
/**
* | output |
* | --- |
* | "Failed to export dataset" |
*
* @param {Datasets_Failed_ExportInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_failed_export = /** @type {((inputs?: Datasets_Failed_ExportInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Failed_ExportInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_failed_export(inputs)
	return __ro.datasets_failed_export(inputs)
});
/**
* | output |
* | --- |
* | "Invalid date" |
*
* @param {Datasets_Invalid_DateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_invalid_date = /** @type {((inputs?: Datasets_Invalid_DateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Invalid_DateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_invalid_date(inputs)
	return __ro.datasets_invalid_date(inputs)
});
/**
* | output |
* | --- |
* | "Missing document id" |
*
* @param {Datasets_Missing_Document_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_missing_document_id = /** @type {((inputs?: Datasets_Missing_Document_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Missing_Document_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_missing_document_id(inputs)
	return __ro.datasets_missing_document_id(inputs)
});
/**
* | output |
* | --- |
* | "Add dataset" |
*
* @param {Datasets_Add_DatasetInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_add_dataset = /** @type {((inputs?: Datasets_Add_DatasetInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Add_DatasetInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_add_dataset(inputs)
	return __ro.datasets_add_dataset(inputs)
});
/**
* | output |
* | --- |
* | "All datasets" |
*
* @param {Datasets_All_DatasetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_all_datasets = /** @type {((inputs?: Datasets_All_DatasetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_All_DatasetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_all_datasets(inputs)
	return __ro.datasets_all_datasets(inputs)
});
/**
* | output |
* | --- |
* | "Retry datasets" |
*
* @param {Datasets_Retry_DatasetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_retry_datasets = /** @type {((inputs?: Datasets_Retry_DatasetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Retry_DatasetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_retry_datasets(inputs)
	return __ro.datasets_retry_datasets(inputs)
});
/**
* | output |
* | --- |
* | "No datasets" |
*
* @param {Datasets_No_DatasetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_datasets = /** @type {((inputs?: Datasets_No_DatasetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_DatasetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_datasets(inputs)
	return __ro.datasets_no_datasets(inputs)
});
/**
* | output |
* | --- |
* | "Dataset actions" |
*
* @param {Datasets_Dataset_ActionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_dataset_actions = /** @type {((inputs?: Datasets_Dataset_ActionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Dataset_ActionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_dataset_actions(inputs)
	return __ro.datasets_dataset_actions(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Datasets_EditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_edit = /** @type {((inputs?: Datasets_EditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_EditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_edit(inputs)
	return __ro.datasets_edit(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Datasets_DeleteInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_delete = /** @type {((inputs?: Datasets_DeleteInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_DeleteInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_delete(inputs)
	return __ro.datasets_delete(inputs)
});
/**
* | output |
* | --- |
* | "Delete failed" |
*
* @param {Datasets_Delete_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_delete_failed = /** @type {((inputs?: Datasets_Delete_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Delete_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_delete_failed(inputs)
	return __ro.datasets_delete_failed(inputs)
});
/**
* | output |
* | --- |
* | "Delete dataset?" |
*
* @param {Datasets_Delete_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_delete_confirm_title = /** @type {((inputs?: Datasets_Delete_Confirm_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Delete_Confirm_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_delete_confirm_title(inputs)
	return __ro.datasets_delete_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete \"{name}\"?" |
*
* @param {Datasets_Delete_Confirm_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_delete_confirm_description = /** @type {((inputs: Datasets_Delete_Confirm_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Delete_Confirm_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_delete_confirm_description(inputs)
	return __ro.datasets_delete_confirm_description(inputs)
});
/**
* | output |
* | --- |
* | "New dataset" |
*
* @param {Datasets_Dialog_Title_NewInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_dialog_title_new = /** @type {((inputs?: Datasets_Dialog_Title_NewInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Dialog_Title_NewInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_dialog_title_new(inputs)
	return __ro.datasets_dialog_title_new(inputs)
});
/**
* | output |
* | --- |
* | "Edit dataset" |
*
* @param {Datasets_Dialog_Title_EditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_dialog_title_edit = /** @type {((inputs?: Datasets_Dialog_Title_EditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Dialog_Title_EditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_dialog_title_edit(inputs)
	return __ro.datasets_dialog_title_edit(inputs)
});
/**
* | output |
* | --- |
* | "Save changes" |
*
* @param {Datasets_Save_ChangesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_save_changes = /** @type {((inputs?: Datasets_Save_ChangesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Save_ChangesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_save_changes(inputs)
	return __ro.datasets_save_changes(inputs)
});
/**
* | output |
* | --- |
* | "Create dataset" |
*
* @param {Datasets_Create_DatasetInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_create_dataset = /** @type {((inputs?: Datasets_Create_DatasetInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Create_DatasetInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_create_dataset(inputs)
	return __ro.datasets_create_dataset(inputs)
});
/**
* | output |
* | --- |
* | "Selected schema" |
*
* @param {Datasets_Selected_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_selected_schema = /** @type {((inputs?: Datasets_Selected_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Selected_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_selected_schema(inputs)
	return __ro.datasets_selected_schema(inputs)
});
/**
* | output |
* | --- |
* | "Loading schemas" |
*
* @param {Datasets_Loading_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_loading_schemas = /** @type {((inputs?: Datasets_Loading_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Loading_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_loading_schemas(inputs)
	return __ro.datasets_loading_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Select schema" |
*
* @param {Datasets_Select_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_select_schema = /** @type {((inputs?: Datasets_Select_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Select_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_select_schema(inputs)
	return __ro.datasets_select_schema(inputs)
});
/**
* | output |
* | --- |
* | "No fields selected" |
*
* @param {Datasets_No_Fields_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_fields_selected = /** @type {((inputs?: Datasets_No_Fields_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Fields_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_fields_selected(inputs)
	return __ro.datasets_no_fields_selected(inputs)
});
/**
* | output |
* | --- |
* | "1 field selected" |
*
* @param {Datasets_One_Field_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_one_field_selected = /** @type {((inputs?: Datasets_One_Field_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_One_Field_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_one_field_selected(inputs)
	return __ro.datasets_one_field_selected(inputs)
});
/**
* | output |
* | --- |
* | "{count} fields selected" |
*
* @param {Datasets_Fields_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_fields_selected = /** @type {((inputs: Datasets_Fields_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Fields_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_fields_selected(inputs)
	return __ro.datasets_fields_selected(inputs)
});
/**
* | output |
* | --- |
* | "Collapse {label}" |
*
* @param {Datasets_Collapse_FieldInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_collapse_field = /** @type {((inputs: Datasets_Collapse_FieldInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Collapse_FieldInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_collapse_field(inputs)
	return __ro.datasets_collapse_field(inputs)
});
/**
* | output |
* | --- |
* | "Expand {label}" |
*
* @param {Datasets_Expand_FieldInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_expand_field = /** @type {((inputs: Datasets_Expand_FieldInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Expand_FieldInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_expand_field(inputs)
	return __ro.datasets_expand_field(inputs)
});
/**
* | output |
* | --- |
* | "Select {label}" |
*
* @param {Datasets_Select_FieldInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_select_field = /** @type {((inputs: Datasets_Select_FieldInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Select_FieldInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_select_field(inputs)
	return __ro.datasets_select_field(inputs)
});
/**
* | output |
* | --- |
* | "Dataset name" |
*
* @param {Datasets_Name_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_name_placeholder = /** @type {((inputs?: Datasets_Name_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Name_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_name_placeholder(inputs)
	return __ro.datasets_name_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Search schemas" |
*
* @param {Datasets_Search_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_search_schemas = /** @type {((inputs?: Datasets_Search_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Search_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_search_schemas(inputs)
	return __ro.datasets_search_schemas(inputs)
});
/**
* | output |
* | --- |
* | "No schemas found." |
*
* @param {Datasets_No_Schemas_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_schemas_found = /** @type {((inputs?: Datasets_No_Schemas_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_Schemas_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_schemas_found(inputs)
	return __ro.datasets_no_schemas_found(inputs)
});
/**
* | output |
* | --- |
* | "No fields" |
*
* @param {Datasets_No_FieldsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_no_fields = /** @type {((inputs?: Datasets_No_FieldsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_No_FieldsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_no_fields(inputs)
	return __ro.datasets_no_fields(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Datasets_CancelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_cancel = /** @type {((inputs?: Datasets_CancelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_CancelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_cancel(inputs)
	return __ro.datasets_cancel(inputs)
});
/**
* | output |
* | --- |
* | "JSON" |
*
* @param {Datasets_Json_BadgeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const datasets_json_badge = /** @type {((inputs?: Datasets_Json_BadgeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Datasets_Json_BadgeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.datasets_json_badge(inputs)
	return __ro.datasets_json_badge(inputs)
});
/**
* | output |
* | --- |
* | "Documents" |
*
* @param {Documents_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_page_title = /** @type {((inputs?: Documents_Page_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Page_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_page_title(inputs)
	return __ro.documents_page_title(inputs)
});
/**
* | output |
* | --- |
* | "New OCR job" |
*
* @param {Documents_New_Ocr_JobInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_new_ocr_job = /** @type {((inputs?: Documents_New_Ocr_JobInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_New_Ocr_JobInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_new_ocr_job(inputs)
	return __ro.documents_new_ocr_job(inputs)
});
/**
* | output |
* | --- |
* | "Search by filename..." |
*
* @param {Documents_Search_Filename_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_search_filename_placeholder = /** @type {((inputs?: Documents_Search_Filename_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Search_Filename_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_search_filename_placeholder(inputs)
	return __ro.documents_search_filename_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Search filename" |
*
* @param {Documents_Search_FilenameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_search_filename = /** @type {((inputs?: Documents_Search_FilenameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Search_FilenameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_search_filename(inputs)
	return __ro.documents_search_filename(inputs)
});
/**
* | output |
* | --- |
* | "Date range" |
*
* @param {Documents_Date_RangeInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_date_range = /** @type {((inputs?: Documents_Date_RangeInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Date_RangeInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_date_range(inputs)
	return __ro.documents_date_range(inputs)
});
/**
* | output |
* | --- |
* | "Any" |
*
* @param {Documents_Any_DateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_any_date = /** @type {((inputs?: Documents_Any_DateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Any_DateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_any_date(inputs)
	return __ro.documents_any_date(inputs)
});
/**
* | output |
* | --- |
* | "{start} - {end}" |
*
* @param {Documents_Date_Range_ValueInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_date_range_value = /** @type {((inputs: Documents_Date_Range_ValueInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Date_Range_ValueInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_date_range_value(inputs)
	return __ro.documents_date_range_value(inputs)
});
/**
* | output |
* | --- |
* | "Presets" |
*
* @param {Documents_PresetsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_presets = /** @type {((inputs?: Documents_PresetsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_PresetsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_presets(inputs)
	return __ro.documents_presets(inputs)
});
/**
* | output |
* | --- |
* | "Today" |
*
* @param {Documents_TodayInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_today = /** @type {((inputs?: Documents_TodayInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_TodayInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_today(inputs)
	return __ro.documents_today(inputs)
});
/**
* | output |
* | --- |
* | "This week" |
*
* @param {Documents_This_WeekInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_this_week = /** @type {((inputs?: Documents_This_WeekInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_This_WeekInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_this_week(inputs)
	return __ro.documents_this_week(inputs)
});
/**
* | output |
* | --- |
* | "This month" |
*
* @param {Documents_This_MonthInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_this_month = /** @type {((inputs?: Documents_This_MonthInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_This_MonthInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_this_month(inputs)
	return __ro.documents_this_month(inputs)
});
/**
* | output |
* | --- |
* | "Clear" |
*
* @param {Documents_ClearInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_clear = /** @type {((inputs?: Documents_ClearInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_ClearInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_clear(inputs)
	return __ro.documents_clear(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Documents_ApplyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_apply = /** @type {((inputs?: Documents_ApplyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_ApplyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_apply(inputs)
	return __ro.documents_apply(inputs)
});
/**
* | output |
* | --- |
* | "Filter by Collection" |
*
* @param {Documents_Filter_By_CollectionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_filter_by_collection = /** @type {((inputs?: Documents_Filter_By_CollectionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Filter_By_CollectionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_filter_by_collection(inputs)
	return __ro.documents_filter_by_collection(inputs)
});
/**
* | output |
* | --- |
* | "Filter by Schema" |
*
* @param {Documents_Filter_By_SchemaInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_filter_by_schema = /** @type {((inputs?: Documents_Filter_By_SchemaInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Filter_By_SchemaInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_filter_by_schema(inputs)
	return __ro.documents_filter_by_schema(inputs)
});
/**
* | output |
* | --- |
* | "Unknown collection" |
*
* @param {Documents_Unknown_CollectionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_unknown_collection = /** @type {((inputs?: Documents_Unknown_CollectionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Unknown_CollectionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_unknown_collection(inputs)
	return __ro.documents_unknown_collection(inputs)
});
/**
* | output |
* | --- |
* | "All collections" |
*
* @param {Documents_All_CollectionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_all_collections = /** @type {((inputs?: Documents_All_CollectionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_All_CollectionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_all_collections(inputs)
	return __ro.documents_all_collections(inputs)
});
/**
* | output |
* | --- |
* | "All schemas" |
*
* @param {Documents_All_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_all_schemas = /** @type {((inputs?: Documents_All_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_All_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_all_schemas(inputs)
	return __ro.documents_all_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Missing document id" |
*
* @param {Documents_Missing_Document_IdInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_missing_document_id = /** @type {((inputs?: Documents_Missing_Document_IdInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Missing_Document_IdInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_missing_document_id(inputs)
	return __ro.documents_missing_document_id(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load documents" |
*
* @param {Documents_Failed_Load_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_load_documents = /** @type {((inputs?: Documents_Failed_Load_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Load_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_load_documents(inputs)
	return __ro.documents_failed_load_documents(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load document" |
*
* @param {Documents_Failed_Load_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_load_document = /** @type {((inputs?: Documents_Failed_Load_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Load_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_load_document(inputs)
	return __ro.documents_failed_load_document(inputs)
});
/**
* | output |
* | --- |
* | "Failed to delete document" |
*
* @param {Documents_Failed_Delete_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_delete_document = /** @type {((inputs?: Documents_Failed_Delete_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Delete_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_delete_document(inputs)
	return __ro.documents_failed_delete_document(inputs)
});
/**
* | output |
* | --- |
* | "Failed to update document" |
*
* @param {Documents_Failed_Update_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_update_document = /** @type {((inputs?: Documents_Failed_Update_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Update_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_update_document(inputs)
	return __ro.documents_failed_update_document(inputs)
});
/**
* | output |
* | --- |
* | "Failed to delete documents" |
*
* @param {Documents_Failed_Delete_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_delete_documents = /** @type {((inputs?: Documents_Failed_Delete_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Delete_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_delete_documents(inputs)
	return __ro.documents_failed_delete_documents(inputs)
});
/**
* | output |
* | --- |
* | "Failed to move documents" |
*
* @param {Documents_Failed_Move_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_move_documents = /** @type {((inputs?: Documents_Failed_Move_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_Move_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_move_documents(inputs)
	return __ro.documents_failed_move_documents(inputs)
});
/**
* | output |
* | --- |
* | "Failed to download documents" |
*
* @param {Documents_Failed_DownloadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_download = /** @type {((inputs?: Documents_Failed_DownloadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_DownloadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_download(inputs)
	return __ro.documents_failed_download(inputs)
});
/**
* | output |
* | --- |
* | "Invalid date" |
*
* @param {Documents_Invalid_DateInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_invalid_date = /** @type {((inputs?: Documents_Invalid_DateInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Invalid_DateInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_invalid_date(inputs)
	return __ro.documents_invalid_date(inputs)
});
/**
* | output |
* | --- |
* | "Select all documents on this page" |
*
* @param {Documents_Select_All_On_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_select_all_on_page = /** @type {((inputs?: Documents_Select_All_On_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Select_All_On_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_select_all_on_page(inputs)
	return __ro.documents_select_all_on_page(inputs)
});
/**
* | output |
* | --- |
* | "Select {name}" |
*
* @param {Documents_Select_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_select_document = /** @type {((inputs: Documents_Select_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Select_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_select_document(inputs)
	return __ro.documents_select_document(inputs)
});
/**
* | output |
* | --- |
* | "Filename" |
*
* @param {Documents_Filename_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_filename_column = /** @type {((inputs?: Documents_Filename_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Filename_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_filename_column(inputs)
	return __ro.documents_filename_column(inputs)
});
/**
* | output |
* | --- |
* | "Collections" |
*
* @param {Documents_Collections_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collections_column = /** @type {((inputs?: Documents_Collections_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collections_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collections_column(inputs)
	return __ro.documents_collections_column(inputs)
});
/**
* | output |
* | --- |
* | "Pages" |
*
* @param {Documents_Pages_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_pages_column = /** @type {((inputs?: Documents_Pages_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Pages_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_pages_column(inputs)
	return __ro.documents_pages_column(inputs)
});
/**
* | output |
* | --- |
* | "Created" |
*
* @param {Documents_Created_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_created_column = /** @type {((inputs?: Documents_Created_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Created_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_created_column(inputs)
	return __ro.documents_created_column(inputs)
});
/**
* | output |
* | --- |
* | "File size" |
*
* @param {Documents_File_Size_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_file_size_column = /** @type {((inputs?: Documents_File_Size_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_File_Size_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_file_size_column(inputs)
	return __ro.documents_file_size_column(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date ascending" |
*
* @param {Documents_Sort_Created_AscendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_sort_created_ascending = /** @type {((inputs?: Documents_Sort_Created_AscendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Sort_Created_AscendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_sort_created_ascending(inputs)
	return __ro.documents_sort_created_ascending(inputs)
});
/**
* | output |
* | --- |
* | "Sort by created date descending" |
*
* @param {Documents_Sort_Created_DescendingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_sort_created_descending = /** @type {((inputs?: Documents_Sort_Created_DescendingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Sort_Created_DescendingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_sort_created_descending(inputs)
	return __ro.documents_sort_created_descending(inputs)
});
/**
* | output |
* | --- |
* | "Collection not found" |
*
* @param {Documents_Collection_Not_Found_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_not_found_title = /** @type {((inputs?: Documents_Collection_Not_Found_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Not_Found_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_not_found_title(inputs)
	return __ro.documents_collection_not_found_title(inputs)
});
/**
* | output |
* | --- |
* | "This collection does not exist." |
*
* @param {Documents_Collection_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_not_found_body = /** @type {((inputs?: Documents_Collection_Not_Found_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Not_Found_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_not_found_body(inputs)
	return __ro.documents_collection_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "View all documents" |
*
* @param {Documents_View_All_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_view_all_documents = /** @type {((inputs?: Documents_View_All_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_View_All_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_view_all_documents(inputs)
	return __ro.documents_view_all_documents(inputs)
});
/**
* | output |
* | --- |
* | "Retry" |
*
* @param {Documents_RetryInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_retry = /** @type {((inputs?: Documents_RetryInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_RetryInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_retry(inputs)
	return __ro.documents_retry(inputs)
});
/**
* | output |
* | --- |
* | "No documents found" |
*
* @param {Documents_No_Documents_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_documents_found = /** @type {((inputs?: Documents_No_Documents_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Documents_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_documents_found(inputs)
	return __ro.documents_no_documents_found(inputs)
});
/**
* | output |
* | --- |
* | "No documents match your filter criteria. Try resetting them to view your files." |
*
* @param {Documents_Empty_Filtered_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_empty_filtered_body = /** @type {((inputs?: Documents_Empty_Filtered_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Empty_Filtered_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_empty_filtered_body(inputs)
	return __ro.documents_empty_filtered_body(inputs)
});
/**
* | output |
* | --- |
* | "You haven't processed any documents yet. Upload a document to get started." |
*
* @param {Documents_Empty_Unfiltered_BodyInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_empty_unfiltered_body = /** @type {((inputs?: Documents_Empty_Unfiltered_BodyInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Empty_Unfiltered_BodyInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_empty_unfiltered_body(inputs)
	return __ro.documents_empty_unfiltered_body(inputs)
});
/**
* | output |
* | --- |
* | "Clear filters" |
*
* @param {Documents_Clear_FiltersInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_clear_filters = /** @type {((inputs?: Documents_Clear_FiltersInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Clear_FiltersInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_clear_filters(inputs)
	return __ro.documents_clear_filters(inputs)
});
/**
* | output |
* | --- |
* | "Process first document" |
*
* @param {Documents_Process_First_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_process_first_document = /** @type {((inputs?: Documents_Process_First_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Process_First_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_process_first_document(inputs)
	return __ro.documents_process_first_document(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} document on this page." |
*
* @param {Documents_Showing_Documents_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_showing_documents_one = /** @type {((inputs: Documents_Showing_Documents_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Showing_Documents_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_showing_documents_one(inputs)
	return __ro.documents_showing_documents_one(inputs)
});
/**
* | output |
* | --- |
* | "Showing {count} documents on this page." |
*
* @param {Documents_Showing_Documents_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_showing_documents_other = /** @type {((inputs: Documents_Showing_Documents_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Showing_Documents_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_showing_documents_other(inputs)
	return __ro.documents_showing_documents_other(inputs)
});
/**
* | output |
* | --- |
* | "No documents to show." |
*
* @param {Documents_No_Documents_To_ShowInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_documents_to_show = /** @type {((inputs?: Documents_No_Documents_To_ShowInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Documents_To_ShowInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_documents_to_show(inputs)
	return __ro.documents_no_documents_to_show(inputs)
});
/**
* | output |
* | --- |
* | "Rows per page" |
*
* @param {Documents_Rows_Per_PageInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_rows_per_page = /** @type {((inputs?: Documents_Rows_Per_PageInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Rows_Per_PageInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_rows_per_page(inputs)
	return __ro.documents_rows_per_page(inputs)
});
/**
* | output |
* | --- |
* | "Previous" |
*
* @param {Documents_PreviousInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_previous = /** @type {((inputs?: Documents_PreviousInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_PreviousInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_previous(inputs)
	return __ro.documents_previous(inputs)
});
/**
* | output |
* | --- |
* | "Next" |
*
* @param {Documents_NextInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_next = /** @type {((inputs?: Documents_NextInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_NextInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_next(inputs)
	return __ro.documents_next(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Documents_DeleteInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete = /** @type {((inputs?: Documents_DeleteInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_DeleteInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete(inputs)
	return __ro.documents_delete(inputs)
});
/**
* | output |
* | --- |
* | "Delete document?" |
*
* @param {Documents_Delete_Single_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_single_title = /** @type {((inputs?: Documents_Delete_Single_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Single_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_single_title(inputs)
	return __ro.documents_delete_single_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete \"{name}\"? This action cannot be undone." |
*
* @param {Documents_Delete_Single_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_single_description = /** @type {((inputs: Documents_Delete_Single_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Single_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_single_description(inputs)
	return __ro.documents_delete_single_description(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} document?" |
*
* @param {Documents_Delete_Bulk_Title_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_bulk_title_one = /** @type {((inputs: Documents_Delete_Bulk_Title_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Bulk_Title_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_bulk_title_one(inputs)
	return __ro.documents_delete_bulk_title_one(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} documents?" |
*
* @param {Documents_Delete_Bulk_Title_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_bulk_title_other = /** @type {((inputs: Documents_Delete_Bulk_Title_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Bulk_Title_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_bulk_title_other(inputs)
	return __ro.documents_delete_bulk_title_other(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} selected document? This action cannot be undone." |
*
* @param {Documents_Delete_Bulk_Description_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_bulk_description_one = /** @type {((inputs: Documents_Delete_Bulk_Description_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Bulk_Description_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_bulk_description_one(inputs)
	return __ro.documents_delete_bulk_description_one(inputs)
});
/**
* | output |
* | --- |
* | "Delete {count} selected documents? This action cannot be undone." |
*
* @param {Documents_Delete_Bulk_Description_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_bulk_description_other = /** @type {((inputs: Documents_Delete_Bulk_Description_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Bulk_Description_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_bulk_description_other(inputs)
	return __ro.documents_delete_bulk_description_other(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Documents_Selected_Count_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_selected_count_one = /** @type {((inputs: Documents_Selected_Count_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Selected_Count_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_selected_count_one(inputs)
	return __ro.documents_selected_count_one(inputs)
});
/**
* | output |
* | --- |
* | "{count} selected" |
*
* @param {Documents_Selected_Count_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_selected_count_other = /** @type {((inputs: Documents_Selected_Count_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Selected_Count_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_selected_count_other(inputs)
	return __ro.documents_selected_count_other(inputs)
});
/**
* | output |
* | --- |
* | "Download selected documents" |
*
* @param {Documents_Download_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_download_selected = /** @type {((inputs?: Documents_Download_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Download_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_download_selected(inputs)
	return __ro.documents_download_selected(inputs)
});
/**
* | output |
* | --- |
* | "Download" |
*
* @param {Documents_DownloadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_download = /** @type {((inputs?: Documents_DownloadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_DownloadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_download(inputs)
	return __ro.documents_download(inputs)
});
/**
* | output |
* | --- |
* | "Downloading..." |
*
* @param {Documents_DownloadingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_downloading = /** @type {((inputs?: Documents_DownloadingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_DownloadingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_downloading(inputs)
	return __ro.documents_downloading(inputs)
});
/**
* | output |
* | --- |
* | "Move" |
*
* @param {Documents_MoveInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_move = /** @type {((inputs?: Documents_MoveInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_MoveInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_move(inputs)
	return __ro.documents_move(inputs)
});
/**
* | output |
* | --- |
* | "Moving..." |
*
* @param {Documents_MovingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_moving = /** @type {((inputs?: Documents_MovingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_MovingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_moving(inputs)
	return __ro.documents_moving(inputs)
});
/**
* | output |
* | --- |
* | "Deleting..." |
*
* @param {Documents_DeletingInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_deleting = /** @type {((inputs?: Documents_DeletingInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_DeletingInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_deleting(inputs)
	return __ro.documents_deleting(inputs)
});
/**
* | output |
* | --- |
* | "Open actions for {name}" |
*
* @param {Documents_Open_Actions_ForInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_open_actions_for = /** @type {((inputs: Documents_Open_Actions_ForInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Open_Actions_ForInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_open_actions_for(inputs)
	return __ro.documents_open_actions_for(inputs)
});
/**
* | output |
* | --- |
* | "Preview" |
*
* @param {Documents_PreviewInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_preview = /** @type {((inputs?: Documents_PreviewInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_PreviewInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_preview(inputs)
	return __ro.documents_preview(inputs)
});
/**
* | output |
* | --- |
* | "Rename" |
*
* @param {Documents_RenameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_rename = /** @type {((inputs?: Documents_RenameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_RenameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_rename(inputs)
	return __ro.documents_rename(inputs)
});
/**
* | output |
* | --- |
* | "Failed to rename document" |
*
* @param {Documents_Failed_RenameInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_failed_rename = /** @type {((inputs?: Documents_Failed_RenameInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Failed_RenameInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_failed_rename(inputs)
	return __ro.documents_failed_rename(inputs)
});
/**
* | output |
* | --- |
* | "Rename {name}" |
*
* @param {Documents_Rename_FileInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_rename_file = /** @type {((inputs: Documents_Rename_FileInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Rename_FileInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_rename_file(inputs)
	return __ro.documents_rename_file(inputs)
});
/**
* | output |
* | --- |
* | "Preview {name}" |
*
* @param {Documents_Preview_FileInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_preview_file = /** @type {((inputs: Documents_Preview_FileInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Preview_FileInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_preview_file(inputs)
	return __ro.documents_preview_file(inputs)
});
/**
* | output |
* | --- |
* | "Download document" |
*
* @param {Documents_Download_Dialog_Title_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_download_dialog_title_one = /** @type {((inputs?: Documents_Download_Dialog_Title_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Download_Dialog_Title_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_download_dialog_title_one(inputs)
	return __ro.documents_download_dialog_title_one(inputs)
});
/**
* | output |
* | --- |
* | "Download {count} documents" |
*
* @param {Documents_Download_Dialog_Title_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_download_dialog_title_other = /** @type {((inputs: Documents_Download_Dialog_Title_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Download_Dialog_Title_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_download_dialog_title_other(inputs)
	return __ro.documents_download_dialog_title_other(inputs)
});
/**
* | output |
* | --- |
* | "Selected documents" |
*
* @param {Documents_Selected_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_selected_documents = /** @type {((inputs?: Documents_Selected_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Selected_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_selected_documents(inputs)
	return __ro.documents_selected_documents(inputs)
});
/**
* | output |
* | --- |
* | "Markdown" |
*
* @param {Documents_Format_MarkdownInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_format_markdown = /** @type {((inputs?: Documents_Format_MarkdownInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Format_MarkdownInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_format_markdown(inputs)
	return __ro.documents_format_markdown(inputs)
});
/**
* | output |
* | --- |
* | "HTML" |
*
* @param {Documents_Format_HtmlInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_format_html = /** @type {((inputs?: Documents_Format_HtmlInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Format_HtmlInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_format_html(inputs)
	return __ro.documents_format_html(inputs)
});
/**
* | output |
* | --- |
* | "JSON" |
*
* @param {Documents_Format_JsonInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_format_json = /** @type {((inputs?: Documents_Format_JsonInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Format_JsonInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_format_json(inputs)
	return __ro.documents_format_json(inputs)
});
/**
* | output |
* | --- |
* | "Preparing download..." |
*
* @param {Documents_Preparing_DownloadInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_preparing_download = /** @type {((inputs?: Documents_Preparing_DownloadInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Preparing_DownloadInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_preparing_download(inputs)
	return __ro.documents_preparing_download(inputs)
});
/**
* | output |
* | --- |
* | "No collections selected" |
*
* @param {Documents_No_Collections_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_collections_selected = /** @type {((inputs?: Documents_No_Collections_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Collections_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_collections_selected(inputs)
	return __ro.documents_no_collections_selected(inputs)
});
/**
* | output |
* | --- |
* | "1 collection selected" |
*
* @param {Documents_One_Collection_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_one_collection_selected = /** @type {((inputs?: Documents_One_Collection_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_One_Collection_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_one_collection_selected(inputs)
	return __ro.documents_one_collection_selected(inputs)
});
/**
* | output |
* | --- |
* | "{count} collections selected" |
*
* @param {Documents_Collections_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collections_selected = /** @type {((inputs: Documents_Collections_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collections_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collections_selected(inputs)
	return __ro.documents_collections_selected(inputs)
});
/**
* | output |
* | --- |
* | "Remove from all" |
*
* @param {Documents_Remove_From_AllInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_remove_from_all = /** @type {((inputs?: Documents_Remove_From_AllInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Remove_From_AllInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_remove_from_all(inputs)
	return __ro.documents_remove_from_all(inputs)
});
/**
* | output |
* | --- |
* | "Move documents" |
*
* @param {Documents_Move_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_move_documents = /** @type {((inputs?: Documents_Move_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Move_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_move_documents(inputs)
	return __ro.documents_move_documents(inputs)
});
/**
* | output |
* | --- |
* | "Replace collections for 1 selected document." |
*
* @param {Documents_Move_Description_OneInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_move_description_one = /** @type {((inputs?: Documents_Move_Description_OneInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Move_Description_OneInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_move_description_one(inputs)
	return __ro.documents_move_description_one(inputs)
});
/**
* | output |
* | --- |
* | "Replace collections for {count} selected documents." |
*
* @param {Documents_Move_Description_OtherInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_move_description_other = /** @type {((inputs: Documents_Move_Description_OtherInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Move_Description_OtherInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_move_description_other(inputs)
	return __ro.documents_move_description_other(inputs)
});
/**
* | output |
* | --- |
* | "Collections" |
*
* @param {Documents_Collections_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collections_label = /** @type {((inputs?: Documents_Collections_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collections_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collections_label(inputs)
	return __ro.documents_collections_label(inputs)
});
/**
* | output |
* | --- |
* | "Search collections" |
*
* @param {Documents_Search_CollectionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_search_collections = /** @type {((inputs?: Documents_Search_CollectionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Search_CollectionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_search_collections(inputs)
	return __ro.documents_search_collections(inputs)
});
/**
* | output |
* | --- |
* | "Loading collections" |
*
* @param {Documents_Loading_CollectionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_loading_collections = /** @type {((inputs?: Documents_Loading_CollectionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Loading_CollectionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_loading_collections(inputs)
	return __ro.documents_loading_collections(inputs)
});
/**
* | output |
* | --- |
* | "No collections found." |
*
* @param {Documents_No_Collections_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_collections_found = /** @type {((inputs?: Documents_No_Collections_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Collections_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_collections_found(inputs)
	return __ro.documents_no_collections_found(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Documents_CancelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_cancel = /** @type {((inputs?: Documents_CancelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_CancelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_cancel(inputs)
	return __ro.documents_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Collections" |
*
* @param {Documents_Collections_Nav_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collections_nav_label = /** @type {((inputs?: Documents_Collections_Nav_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collections_Nav_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collections_nav_label(inputs)
	return __ro.documents_collections_nav_label(inputs)
});
/**
* | output |
* | --- |
* | "Add collection" |
*
* @param {Documents_Add_CollectionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_add_collection = /** @type {((inputs?: Documents_Add_CollectionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Add_CollectionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_add_collection(inputs)
	return __ro.documents_add_collection(inputs)
});
/**
* | output |
* | --- |
* | "All documents" |
*
* @param {Documents_All_DocumentsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_all_documents = /** @type {((inputs?: Documents_All_DocumentsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_All_DocumentsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_all_documents(inputs)
	return __ro.documents_all_documents(inputs)
});
/**
* | output |
* | --- |
* | "Retry collections" |
*
* @param {Documents_Retry_CollectionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_retry_collections = /** @type {((inputs?: Documents_Retry_CollectionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Retry_CollectionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_retry_collections(inputs)
	return __ro.documents_retry_collections(inputs)
});
/**
* | output |
* | --- |
* | "No collections" |
*
* @param {Documents_No_CollectionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_collections = /** @type {((inputs?: Documents_No_CollectionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_CollectionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_collections(inputs)
	return __ro.documents_no_collections(inputs)
});
/**
* | output |
* | --- |
* | "Collection actions" |
*
* @param {Documents_Collection_ActionsInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_actions = /** @type {((inputs?: Documents_Collection_ActionsInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_ActionsInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_actions(inputs)
	return __ro.documents_collection_actions(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Documents_EditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_edit = /** @type {((inputs?: Documents_EditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_EditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_edit(inputs)
	return __ro.documents_edit(inputs)
});
/**
* | output |
* | --- |
* | "Delete failed" |
*
* @param {Documents_Delete_FailedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_failed = /** @type {((inputs?: Documents_Delete_FailedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_FailedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_failed(inputs)
	return __ro.documents_delete_failed(inputs)
});
/**
* | output |
* | --- |
* | "Delete collection?" |
*
* @param {Documents_Delete_Collection_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_collection_title = /** @type {((inputs?: Documents_Delete_Collection_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Collection_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_collection_title(inputs)
	return __ro.documents_delete_collection_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete \"{name}\"? Documents remain available in All documents." |
*
* @param {Documents_Delete_Collection_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_delete_collection_description = /** @type {((inputs: Documents_Delete_Collection_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Delete_Collection_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_delete_collection_description(inputs)
	return __ro.documents_delete_collection_description(inputs)
});
/**
* | output |
* | --- |
* | "New collection" |
*
* @param {Documents_Collection_Dialog_Title_NewInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_dialog_title_new = /** @type {((inputs?: Documents_Collection_Dialog_Title_NewInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Dialog_Title_NewInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_dialog_title_new(inputs)
	return __ro.documents_collection_dialog_title_new(inputs)
});
/**
* | output |
* | --- |
* | "Edit collection" |
*
* @param {Documents_Collection_Dialog_Title_EditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_dialog_title_edit = /** @type {((inputs?: Documents_Collection_Dialog_Title_EditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Dialog_Title_EditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_dialog_title_edit(inputs)
	return __ro.documents_collection_dialog_title_edit(inputs)
});
/**
* | output |
* | --- |
* | "Group documents by name and optional schema filters." |
*
* @param {Documents_Collection_Dialog_Description_NewInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_dialog_description_new = /** @type {((inputs?: Documents_Collection_Dialog_Description_NewInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Dialog_Description_NewInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_dialog_description_new(inputs)
	return __ro.documents_collection_dialog_description_new(inputs)
});
/**
* | output |
* | --- |
* | "Update the collection name and schema filters." |
*
* @param {Documents_Collection_Dialog_Description_EditInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_dialog_description_edit = /** @type {((inputs?: Documents_Collection_Dialog_Description_EditInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Dialog_Description_EditInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_dialog_description_edit(inputs)
	return __ro.documents_collection_dialog_description_edit(inputs)
});
/**
* | output |
* | --- |
* | "Save changes" |
*
* @param {Documents_Save_ChangesInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_save_changes = /** @type {((inputs?: Documents_Save_ChangesInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Save_ChangesInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_save_changes(inputs)
	return __ro.documents_save_changes(inputs)
});
/**
* | output |
* | --- |
* | "Create collection" |
*
* @param {Documents_Create_CollectionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_create_collection = /** @type {((inputs?: Documents_Create_CollectionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Create_CollectionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_create_collection(inputs)
	return __ro.documents_create_collection(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Documents_Name_ColumnInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_name_column = /** @type {((inputs?: Documents_Name_ColumnInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Name_ColumnInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_name_column(inputs)
	return __ro.documents_name_column(inputs)
});
/**
* | output |
* | --- |
* | "Invoices, reports, receipts" |
*
* @param {Documents_Collection_Name_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_name_placeholder = /** @type {((inputs?: Documents_Collection_Name_PlaceholderInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Name_PlaceholderInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_name_placeholder(inputs)
	return __ro.documents_collection_name_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Schemas" |
*
* @param {Documents_Schemas_LabelInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_schemas_label = /** @type {((inputs?: Documents_Schemas_LabelInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Schemas_LabelInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_schemas_label(inputs)
	return __ro.documents_schemas_label(inputs)
});
/**
* | output |
* | --- |
* | "No schemas selected" |
*
* @param {Documents_No_Schemas_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_schemas_selected = /** @type {((inputs?: Documents_No_Schemas_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Schemas_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_schemas_selected(inputs)
	return __ro.documents_no_schemas_selected(inputs)
});
/**
* | output |
* | --- |
* | "1 schema selected" |
*
* @param {Documents_One_Schema_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_one_schema_selected = /** @type {((inputs?: Documents_One_Schema_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_One_Schema_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_one_schema_selected(inputs)
	return __ro.documents_one_schema_selected(inputs)
});
/**
* | output |
* | --- |
* | "{count} schemas selected" |
*
* @param {Documents_Schemas_SelectedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_schemas_selected = /** @type {((inputs: Documents_Schemas_SelectedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Schemas_SelectedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_schemas_selected(inputs)
	return __ro.documents_schemas_selected(inputs)
});
/**
* | output |
* | --- |
* | "Search schemas" |
*
* @param {Documents_Search_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_search_schemas = /** @type {((inputs?: Documents_Search_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Search_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_search_schemas(inputs)
	return __ro.documents_search_schemas(inputs)
});
/**
* | output |
* | --- |
* | "Loading schemas" |
*
* @param {Documents_Loading_SchemasInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_loading_schemas = /** @type {((inputs?: Documents_Loading_SchemasInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Loading_SchemasInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_loading_schemas(inputs)
	return __ro.documents_loading_schemas(inputs)
});
/**
* | output |
* | --- |
* | "No schemas found." |
*
* @param {Documents_No_Schemas_FoundInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_schemas_found = /** @type {((inputs?: Documents_No_Schemas_FoundInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Schemas_FoundInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_schemas_found(inputs)
	return __ro.documents_no_schemas_found(inputs)
});
/**
* | output |
* | --- |
* | "Only matching schema documents are shown in this collection." |
*
* @param {Documents_Collection_Schema_HintInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_collection_schema_hint = /** @type {((inputs?: Documents_Collection_Schema_HintInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Collection_Schema_HintInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_collection_schema_hint(inputs)
	return __ro.documents_collection_schema_hint(inputs)
});
/**
* | output |
* | --- |
* | "Document preview" |
*
* @param {Documents_Preview_Fallback_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_preview_fallback_title = /** @type {((inputs?: Documents_Preview_Fallback_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Preview_Fallback_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_preview_fallback_title(inputs)
	return __ro.documents_preview_fallback_title(inputs)
});
/**
* | output |
* | --- |
* | "Review extracted markdown and JSON for this document." |
*
* @param {Documents_Preview_DescriptionInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_preview_description = /** @type {((inputs?: Documents_Preview_DescriptionInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Preview_DescriptionInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_preview_description(inputs)
	return __ro.documents_preview_description(inputs)
});
/**
* | output |
* | --- |
* | "Rename document" |
*
* @param {Documents_Rename_Document_TitleInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_rename_document_title = /** @type {((inputs?: Documents_Rename_Document_TitleInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Rename_Document_TitleInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_rename_document_title(inputs)
	return __ro.documents_rename_document_title(inputs)
});
/**
* | output |
* | --- |
* | "Loading document..." |
*
* @param {Documents_Loading_DocumentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_loading_document = /** @type {((inputs?: Documents_Loading_DocumentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Loading_DocumentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_loading_document(inputs)
	return __ro.documents_loading_document(inputs)
});
/**
* | output |
* | --- |
* | "Copy Markdown" |
*
* @param {Documents_Copy_MarkdownInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_copy_markdown = /** @type {((inputs?: Documents_Copy_MarkdownInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Copy_MarkdownInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_copy_markdown(inputs)
	return __ro.documents_copy_markdown(inputs)
});
/**
* | output |
* | --- |
* | "Copy HTML" |
*
* @param {Documents_Copy_HtmlInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_copy_html = /** @type {((inputs?: Documents_Copy_HtmlInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Copy_HtmlInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_copy_html(inputs)
	return __ro.documents_copy_html(inputs)
});
/**
* | output |
* | --- |
* | "Copy JSON" |
*
* @param {Documents_Copy_JsonInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_copy_json = /** @type {((inputs?: Documents_Copy_JsonInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_Copy_JsonInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_copy_json(inputs)
	return __ro.documents_copy_json(inputs)
});
/**
* | output |
* | --- |
* | "Copied" |
*
* @param {Documents_CopiedInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_copied = /** @type {((inputs?: Documents_CopiedInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_CopiedInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_copied(inputs)
	return __ro.documents_copied(inputs)
});
/**
* | output |
* | --- |
* | "No JSON annotation available." |
*
* @param {Documents_No_Json_AnnotationInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_json_annotation = /** @type {((inputs?: Documents_No_Json_AnnotationInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Json_AnnotationInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_json_annotation(inputs)
	return __ro.documents_no_json_annotation(inputs)
});
/**
* | output |
* | --- |
* | "No markdown content available." |
*
* @param {Documents_No_Markdown_ContentInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_markdown_content = /** @type {((inputs?: Documents_No_Markdown_ContentInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Markdown_ContentInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_markdown_content(inputs)
	return __ro.documents_no_markdown_content(inputs)
});
/**
* | output |
* | --- |
* | "No document preview available." |
*
* @param {Documents_No_Preview_AvailableInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_no_preview_available = /** @type {((inputs?: Documents_No_Preview_AvailableInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_No_Preview_AvailableInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_no_preview_available(inputs)
	return __ro.documents_no_preview_available(inputs)
});
/**
* | output |
* | --- |
* | "Close" |
*
* @param {Documents_CloseInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_close = /** @type {((inputs?: Documents_CloseInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_CloseInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_close(inputs)
	return __ro.documents_close(inputs)
});
/**
* | output |
* | --- |
* | "More" |
*
* @param {Documents_MoreInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_more = /** @type {((inputs?: Documents_MoreInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_MoreInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_more(inputs)
	return __ro.documents_more(inputs)
});
/**
* | output |
* | --- |
* | "Open" |
*
* @param {Documents_OpenInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_open = /** @type {((inputs?: Documents_OpenInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_OpenInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_open(inputs)
	return __ro.documents_open(inputs)
});
/**
* | output |
* | --- |
* | "Share" |
*
* @param {Documents_ShareInputs} inputs
* @param {{ locale?: "en" | "ro" }} options
* @returns {LocalizedString}
*/
export const documents_share = /** @type {((inputs?: Documents_ShareInputs, options?: { locale?: "en" | "ro" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Documents_ShareInputs, { locale?: "en" | "ro" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.documents_share(inputs)
	return __ro.documents_share(inputs)
});
