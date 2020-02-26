/**
 PreferencesStore.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

protocol PreferencesStore: ObservableObject {
    var libraryPath: URL { get set }

    var useImportDirectory: Bool { get set }

    var goodReadsAPIKey: String { get set }
}
