/**
 PreviewAudiobooks.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

#if DEBUG

    let previewAudiobooks: [MockAudiobook] = load("previewAudiobooks.json")

    class MockLibrary: Library {
        var 🎧📚: [Audiobook] = previewAudiobooks

        func shelve(book _: Audiobook) {
            // non functional
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

    class MockAudiobook: Audiobook, Decodable {
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
    }
#endif
