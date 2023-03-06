/**
 BookCellView.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import SwiftUI

struct BookCellView: View {
    var book: Book

    var body: some View {
        VStack(alignment: .leading) {
            Text(book.title)
                .font(.headline)
                .lineLimit(3)
            HStack {
                Text(book.authorString)
                    .font(.caption)
                    .foregroundColor(.gray)
                if let releaseDate = book.shortReleaseDate {
                    Spacer()
                    Text(releaseDate)
                        .font(.caption)
                        .foregroundColor(.gray)
                }
            }
        }
        .padding(.vertical, 4)
    }
}

struct BookCellView_Previews: PreviewProvider {
    static var previews: some View {
        BookCellView(book: MockController.book)
    }
}
