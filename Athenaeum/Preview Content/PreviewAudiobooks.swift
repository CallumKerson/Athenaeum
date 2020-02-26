/**
 PreviewAudiobooks.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import RealmSwift

#if DEBUG

    let previewAudiobooks: [MockAudiobook] = load("previewAudiobooks.json")

    class MockRepo: Repository {
        var items: [MockAudiobook]

        init() {
            self.items = previewAudiobooks
        }

        func insert(item _: MockAudiobook) throws {
            //
        }

        func update(item _: MockAudiobook) throws {
            //
        }

        func delete(item _: MockAudiobook) throws {
            //
        }

        func objectExists(item _: MockAudiobook) throws -> Bool {
            false
        }
    }

    func load<T: Decodable>(_ filename: String) -> T {
        let data: Data

        guard let file = Bundle.main
            .url(forResource: filename, withExtension: nil)
        else {
            fatalError("Couldn't find \(filename) in main bundle.")
        }

        do {
            data = try Data(contentsOf: file)
        } catch {
            fatalError("Couldn't load \(filename) from main bundle:\n\(error)")
        }

        do {
            let decoder = JSONDecoder()
            return try decoder.decode(T.self, from: data)
        } catch {
            fatalError("Couldn't parse \(filename) as \(T.self):\n\(error)")
        }
    }

    class MockAudiobook: Audiobook, Decodable, Entity {
        var title: String = "Audiobook Title"

        var author: String = "Some Author"

        var location: URL = URL(string: "https://www.goodreads.com")!

        var narrator: String?

        var publicationDate: String?

        var isbn: String?

        var summary: String?

        var series: Series?

        func getCover() -> Data? {
            if let filePath = Bundle.main.url(forResource: title,
                                              withExtension: "jpg") {
                do {
                    return try Data(contentsOf: filePath)
                } catch {
                    return nil
                }
            }
            return nil
        }

        func toStorable() -> StorableMockAudiobook {
            StorableMockAudiobook()
        }
    }

    final class StorableMockAudiobook: Object, Storable {
        @objc dynamic var uuid = ""
        @objc dynamic var title = ""
        @objc dynamic var author = ""
        @objc dynamic var filePath = ""
        @objc dynamic var narrator: String?
        @objc dynamic var publicationDate: String?
        @objc dynamic var isbn: String?
        @objc dynamic var summary: String?
        @objc dynamic var seriesEntry: String?
        @objc dynamic var seriesTitle: String?

        public override static func primaryKey() -> String? {
            "uuid"
        }

        var model: AudiobookFile {
            var series: Series?
            if let seriesTitle = seriesTitle, let seriesEntry = seriesEntry {
                series = Series(title: seriesTitle, entry: seriesEntry)
            }
            return AudiobookFile(title: self.title,
                                 author: self.author,
                                 file: URL(fileURLWithPath: self.filePath),
                                 narrator: self.narrator,
                                 publicationDate: self.publicationDate,
                                 isbn: self.isbn,
                                 summary: self.summary,
                                 series: series)
        }
    }
#endif
