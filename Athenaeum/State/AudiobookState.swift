/**
 AudiobookState.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct AudiobookState: AppState, Codable, Equatable {
    var audiobooks: [UUID: Loadable<Audiobook>] = [:]
    var selectedAudiobookID: UUID?
    var fixMatchDialogDisplayed: Bool = false
}
