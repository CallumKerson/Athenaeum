//
//  AudiobookObject.swift
//  Athenaeum
//
//  Created by Callum Kerson on 17/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation
import RealmSwift

final class PersistableAudiobook: Object, Persistable {
    
    @objc dynamic var id = 0
    @objc dynamic var title = ""
    @objc dynamic var author = ""
    @objc dynamic var filePath = ""
    @objc dynamic var narrator: String?
    @objc dynamic var publicationDate: String?
    @objc dynamic var isbn: String?
    @objc dynamic var summary: String?
    @objc dynamic var entry: String?
    dynamic var series: SeriesObject?
    
    override public static func primaryKey() -> String? {
        return "id"
    }
    
    var model: Audiobook {
        get {
            return Audiobook(title: title,
            author: author,
            file: URL(fileURLWithPath: filePath),
            narrator: narrator,
            publicationDate: publicationDate,
            isbn: isbn,
            summary: summary,
            entry: entry,
            series: series?.model)
        }
    }
}
