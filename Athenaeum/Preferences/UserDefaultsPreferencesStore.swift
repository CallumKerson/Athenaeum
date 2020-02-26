/**
 PreferencesStore.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

final class UserDefaultsPreferencesStore: PreferencesStore {
    static var global = UserDefaultsPreferencesStore()

    let objectWillChange = ObservableObjectPublisher()

    private static var defaultLibraryPath: URL {
        UserDefaultsPreferencesStore.userMusicPath()
            .appendingPathComponent("Athanaeum")
    }

    private static var defaultImportPath: URL {
        defaultLibraryPath.appendingPathComponent("import")
    }

    @UserDefault(key: .libraryPath, defaultValue: defaultLibraryPath)
    var libraryPath: URL {
        willSet { self.objectWillChange.send() }
    }

    @UserDefault(key: .useImport, defaultValue: false)
    var useImportDirectory: Bool {
        willSet { self.objectWillChange.send() }
    }

    @UserDefault(key: .goodReadsAPIKey, defaultValue: "")
    var goodReadsAPIKey: String {
        willSet { self.objectWillChange.send() }
    }

    private var didChangeCancellable: AnyCancellable?

    init() {
        self.didChangeCancellable = NotificationCenter.default
            .publisher(for: UserDefaults.didChangeNotification)
            .map { _ in () }
            .receive(on: DispatchQueue.main)
            .sink(receiveValue: { _ in self.objectWillChange.send() })
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

extension Key {
    static let libraryPath = Key(rawValue: "libraryPath")
    static let importPath = Key(rawValue: "importPath")
    static let useImport = Key(rawValue: "useImport")
    static let goodReadsAPIKey = Key(rawValue: "goodReadsAPIKey")
}
