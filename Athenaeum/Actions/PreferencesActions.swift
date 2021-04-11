/**
 PreferencesActions.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

struct PreferencesActions {
//    struct UpdateAutoImportPreference: AsyncAction {
//        let updateValueTo: Bool
//
//        func execute(state _: AppState?, dispatch: @escaping DispatchFunction) {
//            UserDefaults.standard
//                .set(self.updateValueTo, forKey: PreferencesKey.autoImport.rawValue)
//            dispatch(SetUpdatedAutoImportPreference(updatedValue: self.updateValueTo))
//        }
//    }

    struct SetUpdatedAutoImportPreference: Action {
        let updatedValue: Bool
    }

//    struct UpdateGoodReadsAPIKeyPreference: AsyncAction {
//        let updateValueTo: String
//
//        func execute(state _: AppState?, dispatch: @escaping DispatchFunction) {
//            UserDefaults.standard
//                .set(self.updateValueTo, forKey: PreferencesKey.goodReadsAPIKey.rawValue)
//            dispatch(SetUpdatedGoodReadsAPIKeyPreference(updatedValue: self.updateValueTo))
//        }
//    }

    struct SetUpdatedGoodReadsAPIKeyPreference: Action {
        let updatedValue: String
    }

    struct SetUpdatedPodcastAuthorPreference: Action {
        let updatedValue: String
    }

    struct SetUpdatedPodcastEmailPreference: Action {
        let updatedValue: String
    }

    struct SetUpdatedPodcastHostURL: Action {
        let updatedValue: String
    }
}
