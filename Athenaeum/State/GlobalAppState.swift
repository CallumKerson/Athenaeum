/**
 GlobalAppState.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct GlobalAppState: AppState, Codable, Equatable {
    var audiobookState: AudiobookState
    var preferencesState: PreferencesState

    init(
        audiobookState: AudiobookState = AudiobookState(),
        preferencesState: PreferencesState = PreferencesState()
    ) {
        self.audiobookState = audiobookState
        self.preferencesState = preferencesState
    }
}
