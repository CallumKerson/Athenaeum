/**
 StorableAudiobook.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation
import RealmSwift

final class StorableAudiobook: Object, Storable {
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
