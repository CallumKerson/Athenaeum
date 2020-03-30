/**
 NavigationMasterViewModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine

class NavigationMasterViewModel: ObservableObject {
    let store: Store<GlobalAppState>

    init(store: Store<GlobalAppState>) {
        self.store = store
    }
}
