//
//  Audiobook+Persistable.swift
//  Athenaeum
//
//  Created by Callum Kerson on 17/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation

import Foundation
//import RealmSwift

extension Audiobook: Entity {
    
    private var storableAudiobook: PersistableAudiobook {
        let realmAudiobook = PersistableAudiobook()
        realmAudiobook.uuid = file.sha256HashOfContents
        realmAudiobook.title = title
        realmAudiobook.author = author
        realmAudiobook.filePath = file.path
        realmAudiobook.narrator = narrator
        realmAudiobook.summary = summary
        realmAudiobook.seriesTitle = series?.title
        realmAudiobook.seriesEntry = series?.entry
        return realmAudiobook
    }
    
    func toStorable() -> PersistableAudiobook {
        return storableAudiobook
    }
    
    
}
