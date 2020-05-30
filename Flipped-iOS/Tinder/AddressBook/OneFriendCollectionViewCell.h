//
//  OneFriendCollectionViewCell.h
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@interface OneFriendCollectionViewCell : UICollectionViewCell

@property (strong, nonatomic) UIImageView *headImg;
@property (strong, nonatomic) UILabel *nameLab;

-(void)installDataDic:(NSDictionary *)dic;

@end

NS_ASSUME_NONNULL_END
