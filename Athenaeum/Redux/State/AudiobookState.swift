/**
 AudiobookState.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct AudiobookState: AppState, Codable, Equatable {
    var audiobooks: [UUID: AudioBook] = [:]
    var importsInProgress: Set<UUID> = Set()
    var selectedAudiobook: AudioBook?
}
