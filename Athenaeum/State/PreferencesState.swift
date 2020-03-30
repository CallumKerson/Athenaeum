/**
 PreferencesState.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct PreferencesState: AppState, Codable, Equatable {
    var libraryURL: URL
    var autoImport: Bool
    var goodReadsAPIKey: String

    init(libraryURL: URL, autoImport: Bool, goodReadsAPIKey: String) {
        self.libraryURL = libraryURL
        self.autoImport = autoImport
        self.goodReadsAPIKey = goodReadsAPIKey
    }

    init() {
        if let presetLibraryURL = UserDefaults.standard
            .url(forKey: PreferencesKey.libraryPath.rawValue) {
            self.libraryURL = presetLibraryURL
        } else {
            self.libraryURL = PreferencesState.userMusicPath().appendingPathComponent("Athanaeum")
        }

        self.autoImport = UserDefaults.standard.bool(forKey: PreferencesKey.autoImport.rawValue)

        if let presetGoodReadsAPIKey = UserDefaults.standard
            .string(forKey: PreferencesKey.goodReadsAPIKey.rawValue) {
            self.goodReadsAPIKey = presetGoodReadsAPIKey
        } else {
            self.goodReadsAPIKey = ""
        }
    }

    private static func userMusicPath() -> URL {
        if let path = FileManager.default.urls(for: .musicDirectory,
                                               in: .userDomainMask).first {
            return path
        } else {
            return URL(string: NSHomeDirectory())!
        }
    }
}

struct PreferencesKey: RawRepresentable {
    let rawValue: String
}

extension PreferencesKey {
    static let libraryPath = PreferencesKey(rawValue: "libraryPath")
    static let autoImport = PreferencesKey(rawValue: "autoImport")
    static let goodReadsAPIKey = PreferencesKey(rawValue: "goodReadsAPIKey")
}
