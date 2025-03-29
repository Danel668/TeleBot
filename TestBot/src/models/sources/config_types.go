/* This file is autogen, don`t edit! */

package sources

type Config struct {
    Section Section `json:"section"`
    PeriodicTask PeriodicTask `json:"periodic_task"`
}

type Section struct {
    SectionRecommendation SectionRecommendation `json:"section_recommendation"`
    SectionAdmin SectionAdmin `json:"section_admin"`
    SectionRegistration SectionRegistration `json:"section_registration"`
    SectionCalories SectionCalories `json:"section_calories"`
    SectionReminder SectionReminder `json:"section_reminder"`
    SectionMyProgress SectionMyProgress `json:"section_my_progress"`
}

type SectionRecommendation struct {
    ButtonRecommendation string `json:"button_recommendation"`
    TextAfterButtonRecommendation string `json:"text_after_button_recommendation"`
    TextAfterGettingRecommendation string `json:"text_after_getting_recommendation"`
    TextErrorAfterGettingRecommendation string `json:"text_error_after_getting_recommendation"`
}

type SectionAdmin struct {
    AdminModeSwitching string `json:"admin_mode_switching"`
    TextAfterSwitchingAdminMode string `json:"text_after_switching_admin_mode"`
    TextAfterBanUser string `json:"text_after_ban_user"`
    TextAfterAssignRole string `json:"text_after_assign_role"`
    TextAfterChangingConfig string `json:"text_after_changing_config"`
}

type SectionRegistration struct {
    ButtonRegistration string `json:"button_registration"`
    TextAfterButtonRegistration string `json:"text_after_button_registration"`
    TextAfterSuccessfulRegistration string `json:"text_after_successful_registration"`
    TextAfterFailedRegistration string `json:"text_after_failed_registration"`
}

type SectionCalories struct {
    ButtonCalories string `json:"button_calories"`
    TextAffterButtonCalories string `json:"text_affter_button_calories"`
    TextIfRequestExceedsCondition string `json:"text_if_request_exceeds_condition"`
    InlineButtonChangeRequest string `json:"inline_button_change_request"`
    InlineButtonChangeRequestUnique string `json:"inline_button_change_request_unique"`
}

type SectionReminder struct {
    ButtonReminder string `json:"button_reminder"`
    ButtonAddReminder string `json:"button_add_reminder"`
    TextErrorAfterCheckRegistration string `json:"text_error_after_check_registration"`
    TextAfterButtonReminder string `json:"text_after_button_reminder"`
    TextAfterButtonAddReminder string `json:"text_after_button_add_reminder"`
    ButtonReminderManual string `json:"button_reminder_manual"`
    TextAfterButtonReminderManual string `json:"text_after_button_reminder_manual"`
    TextInvalidPatternAddReminder string `json:"text_invalid_pattern_add_reminder"`
    DifferenceBetweenSendingAndExpiringM int `json:"difference_between_sending_and_expiring_m"`
    TextAfterAddingReminder string `json:"text_after_adding_reminder"`
    TextInvalidTime string `json:"text_invalid_time"`
    ButtonGetMyReminders string `json:"button_get_my_reminders"`
    TextIfRemindersIsEmpty string `json:"text_if_reminders_is_empty"`
}

type SectionMyProgress struct {
    ButtonMyProgress string `json:"button_my_progress"`
    ButtonMyProgressManual string `json:"button_my_progress_manual"`
    ButtonMyProgressSetGoal string `json:"button_my_progress_set_goal"`
    ButtonMyProgressAddRation string `json:"button_my_progress_add_ration"`
    ButtonMyProgressGetMyRations string `json:"button_my_progress_get_my_rations"`
    TextError string `json:"text_error"`
    TextAfterButtonMyProgress string `json:"text_after_button_my_progress"`
    TextAfterButtonMyProgressManual string `json:"text_after_button_my_progress_manual"`
    TextAfterButtonAddRation string `json:"text_after_button_add_ration"`
    TextAfterButtonMyProgressSetGoal string `json:"text_after_button_my_progress_set_goal"`
    TextErrorGoalFormat string `json:"text_error_goal_format"`
    TextAfterInsertingGoal string `json:"text_after_inserting_goal"`
    TextAfterButtonMyProgressGetMyRations string `json:"text_after_button_my_progress_get_my_rations"`
    TextErrorFormatGetMyRations string `json:"text_error_format_get_my_rations"`
    TextEmptyRations string `json:"text_empty_rations"`
    ButtonMyProgressGetMyRationsForLastTime string `json:"button_my_progress_get_my_rations_for_last_time"`
    LastTimeGetRations int `json:"last_time_get_rations"`
}

type PeriodicTask struct {
    SendRemindersSettings SendRemindersSettings `json:"send_reminders_settings"`
    ClearExpiredStateSettings ClearExpiredStateSettings `json:"clear_expired_state_settings"`
    DumpToFileSettings DumpToFileSettings `json:"dump_to_file_settings"`
    LoadConfigSettings LoadConfigSettings `json:"load_config_settings"`
    DeleteExpiredRemindersSettings DeleteExpiredRemindersSettings `json:"delete_expired_reminders_settings"`
}

type SendRemindersSettings struct {
    FrequencyS int `json:"frequency_s"`
    LockName string `json:"lock_name"`
    AnswerPlaceholder string `json:"answer_placeholder"`
}

type ClearExpiredStateSettings struct {
    FrequencyS int `json:"frequency_s"`
    TtlS int `json:"ttl_s"`
}

type DumpToFileSettings struct {
    FrequencyS int `json:"frequency_s"`
}

type LoadConfigSettings struct {
    FrequencyS int `json:"frequency_s"`
}

type DeleteExpiredRemindersSettings struct {
    FrequencyS int `json:"frequency_s"`
    LockName string `json:"lock_name"`
}
