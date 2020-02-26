/**
 PreferencesStore.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

protocol PreferencesStore: ObservableObject {
    var libraryPath: URL { get set }

    var useImportDirectory: Bool { get set }

    var importPath: URL { get set }

    var goodReadsAPIKey: String { get set }
}
