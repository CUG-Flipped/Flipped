//
//  OneFriendTableViewCell.h
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@interface OneFriendTableViewCell : UITableViewCell < UICollectionViewDelegate, UICollectionViewDataSource>

@property (strong, nonatomic) UICollectionView *collection;

@end

NS_ASSUME_NONNULL_END
