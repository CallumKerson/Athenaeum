/**
 NavigationDetailView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import GoodReadsKit
import SwiftUI

struct NavigationDetailView: View {
    @ObservedObject var viewModel: NavigationDetailViewModel

    init(_ viewModel: NavigationDetailViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        ScrollView {
            Unwrap(viewModel.audiobook) { book in
                VStack(alignment: HorizontalAlignment.leading, spacing: 12) {
                    HStack(alignment: VerticalAlignment.center, spacing: 24) {
                        Unwrap(book.getCover()) { coverData in
                            DetailCover(data: coverData)
                                .padding(8)
                        }
                        VStack(alignment: HorizontalAlignment.leading) {
                            if self.viewModel.isGoodReadsConfigured {
                                HStack {
                                    Spacer()
                                    Button(action: {
                                        self.viewModel.showFixMatchDialog()
                                    }) {
                                        Text("Fix Match")
                                            .frame(maxWidth: 100, maxHeight: 24)
                                    }
                                }
                            }
                            Spacer()
                            Unwrap(book.metadata) { metadata in
                                BookText(metadata: metadata)
                            }

                            Spacer()
                        }
                    }

                    Unwrap(book.metadata?.summary) { summary in
                        Divider()
                        SummaryView(summary: summary)
                            .lineLimit(nil)
                    }
                }
            }
        }
        .blur(radius: viewModel.isImporting ? 10 : 0)
        .padding()
        .overlay(
            VStack {
                if viewModel.isImporting {
                    ActivityIndicator()
                        .frame(width: 100, height: 100)
                } else {
                    EmptyView()
                }
            }
        )
        .frame(minWidth: 800, maxWidth: 800, minHeight: 500)
        .sheet(isPresented: $viewModel.fixMatchDialogDisplayed) {
            Unwrap(self.viewModel.audiobook) { audiobook in
                FixMatchModalView(viewModel: FixMatchModalViewModel(audiobook: audiobook,
                                                                    store: self.viewModel.store))
            }
        }
    }
}

struct BookText: View {
    let metadata: BookMetadata

    var body: some View {
        VStack(alignment: HorizontalAlignment.leading) {
            VStack(alignment: HorizontalAlignment.leading) {
                Text(metadata.title).font(.title)
                Unwrap(metadata.series) { series in
                    Text("Book \(series.entry.asString) of \(series.title)")
                }
            }.padding(.bottom)

            Unwrap(metadata.authors?.author) { author in
                Text("Written by ") + Text(author)
            }
            Unwrap(metadata.publicationDate?.getDateAsString) { releaseDate in
                Text("Released on ") + Text(releaseDate)
            }
            Unwrap(metadata.narrators?.author) { narrator in
                Text("Narrated by ") + Text(narrator)
            }
            Unwrap(metadata.isbn) { isbn in
                Text("ISBN ") + Text(isbn)
            }
        }
        .font(.subheadline)
    }
}

struct DetailCover: View {
    let data: Data

    var body: some View {
        Image(nsImage: NSImage(data: data)!)
            .resizable()
            .scaledToFit()
            .frame(width: 400, height: 400)
            .clipShape(RoundedRectangle(cornerRadius: 10))
            .shadow(radius: 5)
    }
}

#if DEBUG
    struct NavigationDetailView_Previews: PreviewProvider {
        static var previews: some View {
            NavigationDetailView(NavigationDetailViewModel(id: sampleAudiobook.id,
                                                           store: sampleStore))
        }
    }
#endif
