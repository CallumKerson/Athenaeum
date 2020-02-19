/**
 Preferences.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class PreferencesDeprecated {
//    struct Key: RawRepresentable {
//        let rawValue: String
//
//        static let libraryPath = Key(rawValue: "libraryPath")
//        static let importPath = Key(rawValue: "importPath")
//        static let useImport = Key(rawValue: "useImport")
//        static let goodReadsAPIKey = Key(rawValue: "goodReadsAPIKey")
//    }

    static var defaultLibraryPath: URL {
        PreferencesDeprecated.userMusicPath().appendingPathComponent("Athanaeum")
    }

    static var defaultImportPath: URL {
        defaultLibraryPath.appendingPathComponent("import")
    }

    static let defaultPreferences: [String: Any] = [
        Key.libraryPath.rawValue: defaultLibraryPath,
        Key.importPath.rawValue: defaultImportPath,
        Key.useImport.rawValue: false,
    ]

    private static let ud = UserDefaults.standard

    static func getURL(for key: Key) -> URL? {
        ud.url(forKey: key.rawValue)
    }

    static func set(_ value: URL, for key: Key) {
        ud.set(value, forKey: key.rawValue)
    }

    static func getBool(for key: Key) -> Bool? {
        ud.bool(forKey: key.rawValue)
    }

    static func set(_ value: Bool, for key: Key) {
        ud.set(value, forKey: key.rawValue)
    }

    static func getString(for key: Key) -> String? {
        ud.string(forKey: key.rawValue)
    }

    static func set(_ value: String, for key: Key) {
        ud.set(value, forKey: key.rawValue)
    }

    static func libraryPath() -> URL {
        log.debug("Getting library path")
        if isTestEnvironment() {
            let url = URL(fileURLWithPath: NSTemporaryDirectory()).appendingPathComponent("Athenaeum_test")
            createDirectoryAt(url)
            return url
        } else {
            let libraryPath = PreferencesDeprecated.getURL(for: .libraryPath) ?? defaultLibraryPath
            log.debug("Library path is \(libraryPath.path)")
            createDirectoryAt(libraryPath)
            return libraryPath
        }
    }

    static func importPath() -> URL {
        if isTestEnvironment() {
            let url = URL(fileURLWithPath: NSTemporaryDirectory()).appendingPathComponent("Athenaeum_import_test")
            createDirectoryAt(url)
            return url
        } else {
            let importPath = PreferencesDeprecated.getURL(for: .importPath) ?? defaultImportPath
            log.debug("Import path is \(importPath.path)")
            createDirectoryAt(importPath)
            return importPath
        }
    }

    private static func createDirectoryAt(_ url: URL) {
        var isDir = ObjCBool(true)
        if FileManager.default.fileExists(atPath: url.path, isDirectory: &isDir) == false {
            log.debug("Creating directory \(url.path)")
            do {
                try FileManager.default.createDirectory(at: url, withIntermediateDirectories: true, attributes: nil)
            } catch {
                log.error("Failed to create directory \(error)")
            }
        }
    }

    private static func userMusicPath() -> URL {
        if let path = FileManager.default.urls(for: .musicDirectory, in: .userDomainMask).first {
            return path
        } else {
            return URL(string: NSHomeDirectory())!
        }
    }

    static func isTestEnvironment() -> Bool {
        if ProcessInfo.processInfo.arguments.contains("UI-TEST") {
            return true
        }

        return ProcessInfo.processInfo.environment["TEST"] != nil
    }
}
