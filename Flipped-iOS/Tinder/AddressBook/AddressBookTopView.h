//
//  AddressBookTopView.h
//  Tinder
//
//  Created by Layer on 2020/5/23.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

@protocol AddressBookTopViewDelegate <NSObject>

- (void)clickAddressBookTopViewButton:(NSInteger)tag;

@end

NS_ASSUME_NONNULL_BEGIN

@interface AddressBookTopView : UIView

@property (strong, nonatomic) UIButton *backBtn;
@property (strong, nonatomic) UIImageView *titleImg;
@property (strong, nonatomic) UIButton *messageBtn;
@property (strong, nonatomic) UIButton *listBtn;
@property (strong, nonatomic) UILabel *lineLab;

@property (weak, nonatomic) id<AddressBookTopViewDelegate> delegate;

@end

NS_ASSUME_NONNULL_END
