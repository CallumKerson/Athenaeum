/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import RealmSwift

class Library {
    static var global = Library()

    let libraryURL: URL
    var importPath: URL?
    let audiobooks: Results<Audiobook>
    let database: Realm

    init() {
        libraryURL = Preferences.libraryPath()
        if Preferences.getBool(for: .useImport) ?? false {
            importPath = Preferences.importPath()
        }
//        try! FileManager.default.removeItem(at: Realm.Configuration.defaultConfiguration.fileURL!)
        database = try! Realm()
        audiobooks = database.objects(Audiobook.self)

//        audiobooks = Library.getAudiobooks(path: libraryURL.path)
    }

    func importAudiobook(fileURL: URL) {
        let newBook = Audiobook.getBookFromFile(fileURL: fileURL)
        let destination = libraryURL
            .appendingPathComponent(newBook.author.removeIllegalCharacters(), isDirectory: true)
            .appendingPathComponent(newBook.title.removeIllegalCharacters())
            .appendingPathExtension("m4b")
        try! FileManager.default.moveItemCreatingIntermediaryDirectories(at: fileURL, to: destination)
    }

    private func getAudiobooks() -> [Audiobook] {
        var books = [Audiobook]()

        let enumerator = FileManager.default.enumerator(atPath: libraryURL.path)
        let filePaths = enumerator?.allObjects as! [String]
        let audiobookFilePaths = filePaths.filter { $0.contains(".m4b") }
        for audiobookPath in audiobookFilePaths {
            print(audiobookPath)
            books.append(Audiobook(value: "\(libraryURL.path)/\(audiobookPath)"))
        }

        books.sort {
            $0.author < $1.author
        }

        return books
    }
}

extension FileManager {
    func moveItemCreatingIntermediaryDirectories(at: URL, to: URL) throws {
        let parentPath = to.deletingLastPathComponent()
        if !fileExists(atPath: parentPath.path) {
            try createDirectory(at: parentPath, withIntermediateDirectories: true, attributes: nil)
        }
        try moveItem(at: at, to: to)
    }
}

extension String {
    func removeIllegalCharacters() -> String {
        var invalidCharacters = CharacterSet(charactersIn: ":/")
        invalidCharacters.formUnion(.newlines)
        invalidCharacters.formUnion(.illegalCharacters)
        invalidCharacters.formUnion(.controlCharacters)

        return components(separatedBy: invalidCharacters)
            .joined(separator: "")
    }
}
