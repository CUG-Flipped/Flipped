//
//  OneFriendTableViewCell.h
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol FriendCollectDelegate <NSObject>

- (void)clickFriendCollectiomItemTag:(NSInteger)tag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface OneFriendTableViewCell : UITableViewCell < UICollectionViewDelegate, UICollectionViewDataSource>

@property (strong, nonatomic) UICollectionView *collection;
@property (strong, nonatomic) NSArray *collList;
@property (weak, nonatomic) id<FriendCollectDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
