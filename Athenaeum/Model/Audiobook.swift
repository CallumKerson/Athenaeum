
//
//  Audiobook.swift
//  Athenaeum
//
//  Created by Callum Kerson on 21/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation

protocol Audiobook {
    var title: String {get set}
    var author: String {get set}
    var location: URL {get set}
    var narrator: String? {get set}
    var publicationDate: String? {get set}
    var isbn: String? {get set}
    var summary: String? {get set}
    var series: Series? {get set}
    
    func getCover() -> Data?
}
