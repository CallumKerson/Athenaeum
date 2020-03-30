/**
 AppReducer.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

func appStateReducer(state: GlobalAppState, action: Action) -> GlobalAppState {
    var state = state
    state.audiobookState = audiobookStateReducer(state: state.audiobookState, action: action)
    state.preferencesState = preferencesStateReducer(state: state.preferencesState, action: action)
    return state
}
