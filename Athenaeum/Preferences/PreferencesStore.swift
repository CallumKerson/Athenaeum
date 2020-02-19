//
//  PreferencesStore.swift
//  Athenaeum
//
//  Created by Callum Kerson on 19/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation
import Combine

final class PreferencesStore: ObservableObject {
    
    static var global = PreferencesStore()
    
    let objectWillChange = PassthroughSubject<Void, Never>()
    
    private static var defaultLibraryPath: URL {
        PreferencesStore.userMusicPath().appendingPathComponent("Athanaeum")
    }

    private static var defaultImportPath: URL {
        defaultLibraryPath.appendingPathComponent("import")
    }

    @UserDefault(key: .libraryPath, defaultValue: defaultLibraryPath)
    var libraryPath: URL

    @UserDefault(key: .useImport, defaultValue: false)
    var useImportDirectory: Bool
    
    @UserDefault(key: .importPath, defaultValue: defaultImportPath)
    var importPath: URL
    
    @UserDefault(key: .goodReadsAPIKey, defaultValue: "")
    var goodReadsAPIKey: String


    private var didChangeCancellable: AnyCancellable?

    init() {
        didChangeCancellable = NotificationCenter.default
            .publisher(for: UserDefaults.didChangeNotification)
            .map { _ in () }
            .receive(on: DispatchQueue.main)
            .subscribe(objectWillChange)
    }
    
    private static func userMusicPath() -> URL {
        if let path = FileManager.default.urls(for: .musicDirectory, in: .userDomainMask).first {
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
