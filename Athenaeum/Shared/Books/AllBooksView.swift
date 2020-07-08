/**
 AllBooksView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct AllBooksView: View {
    @EnvironmentObject private var model: AthenaeumModel

    var body: some View {
        BooksListView(books: model.books)
            .navigationTitle("All Books")
    }
}

struct AllBooksView_Previews: PreviewProvider {
    static var previews: some View {
        AllBooksView()
            .environmentObject(AthenaeumModel())
    }
}
