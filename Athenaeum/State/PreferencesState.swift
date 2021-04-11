/**
 PreferencesState.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct PreferencesState: AppState, Codable, Equatable {
    var libraryURL: URL
    var autoImport: Bool
    var goodReadsAPIKey: String
    var podcastAuthor: String
    var podcastEmail: String
    var podcastHostURL: String

    init(libraryURL: URL = PreferencesState.userMusicPath().appendingPathComponent("Athanaeum"),
         autoImport: Bool = false,
         goodReadsAPIKey: String = "",
         podcastAuthor: String = "",
         podcastEmail _: String = "",
         podcastHostURL: String = "")
    {
        self.libraryURL = libraryURL
        self.autoImport = autoImport
        self.goodReadsAPIKey = goodReadsAPIKey
        self.podcastAuthor = podcastAuthor
        self.podcastEmail = podcastAuthor
        self.podcastHostURL = podcastHostURL
    }

    private static func userMusicPath() -> URL {
        if let path = FileManager.default.urls(for: .musicDirectory,
                                               in: .userDomainMask).first
        {
            return path
        } else {
            return URL(string: NSHomeDirectory())!
        }
    }
}
