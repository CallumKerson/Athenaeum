/**
 AudiobookListView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AudiobookListView: View {
    let books: [Audiobook]

    var body: some View {
        List(books) { book in
            AudiobookCellView(book: book)
        }
    }
}

struct AudiobookList_Previews: PreviewProvider {
    static var previews: some View {
        AudiobookListView(books: [
            Audiobook.getBookFromFile(path: "/Users/ckerson/Music/TWoK.m4b"),
            Audiobook.getBookFromFile(path: "/Users/ckerson/Music/The Gift.m4b"),
            Audiobook.getBookFromFile(path: "/Users/ckerson/Music/Gothe F--ktoSleep_ep6.m4b"),
        ])
    }
}
