//
//  InfoImageView.h
//  Tinder
//
//  Created by Layer on 2020/5/14.
//  Copyright © 2020 Layer. All rights reserved.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@class PageControl;

@interface InfoImageView : UIView <UIScrollViewDelegate>

@property (strong, nonatomic) UIScrollView *scrollView;
@property (assign, nonatomic) CGFloat width;
@property (assign, nonatomic) CGFloat height;
@property (strong, nonatomic) NSMutableArray *dataArr;

@property (assign, nonatomic) float heightPag;

@property (strong, nonatomic) PageControl *pageConl;

- (id)initWithFrame:(CGRect)frame addImageArr:(NSMutableArray*)addImageArr;

@end

NS_ASSUME_NONNULL_END
