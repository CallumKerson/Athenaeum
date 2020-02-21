//
//  Library.swift
//  Athenaeum
//
//  Created by Callum Kerson on 21/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation

protocol Library: ObservableObject {
    var 🎧📚: [Audiobook] {get set}
    func shelve(book: Audiobook)
}
