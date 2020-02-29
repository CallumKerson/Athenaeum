/**
 PreferencesReducer.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

func preferencesStateReducer(state: PreferencesState, action: Action) -> PreferencesState {
    var state = state
    switch action {
    case let action as PreferencesActions.SetUpdatedAutoImportPreference:
        state.autoImport = action.updatedValue
    case let action as PreferencesActions.SetUpdatedGoodReadsAPIKeyPreference:
        state.goodReadsAPIKey = action.updatedValue
    default:
        break
    }
    return state
}
