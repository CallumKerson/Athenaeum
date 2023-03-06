/**
 BookDetailView.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Combine
import SwiftUI

struct BookDetailView: View {
    @ObservedObject var viewModel: BookDetailViewModel
    @State private var showingSheet = false
    let onSave: () -> Void

    @ViewBuilder
    var body: some View {
        if viewModel.loading {
            ProgressView()
                .onAppear(perform: { viewModel.reload() })
        } else if let error = viewModel.error {
            Label(error.description, systemImage: "exclamationmark.triangle")
        } else if let item = viewModel.book {
            HStack {
                loaded(book: item)
                Spacer()
            }
        }
    }

    func loaded(book: Book) -> some View {
        VStack {
            VStack(alignment: .leading) {
                titleAndEditButton(book: book)

                VStack(alignment: .leading) {
                    Text("By \(book.authorString)")
                        .font(.title2)

                    if let series = book.series {
                        Text(series.description)
                            .font(.title3)
                    }
                }.padding(.vertical)

                if let releaseDate = book.mediumReleaseDate {
                    Text("Released \(releaseDate)")
                        .font(.caption)
                }

                if let summary = book.summary {
                    Text(summary).padding(.vertical)
                }
            }
            .padding()
            Spacer()
        }
    }

    func titleAndEditButton(book: Book) -> some View {
        HStack {
            Text(book.title)
                .font(.largeTitle)
                .lineLimit(3)
                .padding(.bottom)
            Spacer()
            Button("Edit") {
                showingSheet.toggle()
            }
            .sheet(isPresented: $showingSheet) {
                EditView(viewModel: EditViewModel(
                    booksLogicController: viewModel.booksLogicController,
                    originalBook: book
                )) {
                    viewModel.reload()
                    onSave()
                }
            }
            .padding(.bottom)
        }
    }
}

// struct BookDetailView_Previews: PreviewProvider {
//    static var previews: some View {
//        BookDetailView(viewModel: nil).loaded(book: Book(id: "001", title: "A New Book", author: [Person(givenNames: "An", familyName: "Author")], summary: "This is a book", releaseDate: ISO8601DateFormatter().date(from: "2021-04-25T08:00:00Z"), series: nil))
//    }
// }
//
