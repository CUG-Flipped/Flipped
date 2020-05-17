//
//  InfoImageView.m
//  Tinder
//
//  Created by Layer on 2020/5/14.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "InfoImageView.h"

@implementation InfoImageView

- (id)initWithFrame:(CGRect)frame addImageArr:(NSMutableArray*)addImageArr {
    if (self = [super initWithFrame:frame]) {
//        if (SCREEN_HEIGHT > 736)  _heightPag = 64;
//        else  _heightPag = 40;
        _heightPag = 10;
        
        _width = self.frame.size.width;
        _height = self.frame.size.height;
        _dataArr = [NSMutableArray arrayWithArray:addImageArr];
        
        if (_dataArr.count == 0) {
            
        }
        else {
            [_dataArr addObject:addImageArr[0]];
            [_dataArr insertObject:[addImageArr lastObject] atIndex:0];
        }
        [self addSubview: self.scrollView];
        [self addSubview:self.pageConl];
    }
    return self;
}

- (UIScrollView*)scrollView {
    if (!_scrollView) {
        _scrollView = [[UIScrollView alloc] initWithFrame:CGRectMake(0, 0, _width, _height)]; // 大小
        _scrollView.contentSize = CGSizeMake(_width * _dataArr.count, _height); // 滑动范围
        _scrollView.contentOffset = CGPointMake(_width, 0); // 起始页
        _scrollView.pagingEnabled = YES; // 分页效果
        _scrollView.bounces = YES; // 弹簧效果
        _scrollView.delegate = self;
        _scrollView.showsVerticalScrollIndicator = NO; // 禁止水平滑动
        _scrollView.showsHorizontalScrollIndicator = NO; // 隐藏滑动条

        if (_dataArr.count == 0) {
            UIImageView *imageV = [[UIImageView alloc] initWithFrame:CGRectMake(_width, 0, _width, _height)];
            imageV.image = [UIImage imageNamed:@"background_1"];
            [_scrollView addSubview: imageV];
        }
        else {
            for (int i = 0; i < _dataArr.count; i++) {
                UIImageView *imageV = [[UIImageView alloc] initWithFrame:CGRectMake(_width, 0, _width, _height)];
                [imageV sd_setImageWithURL:[NSURL URLWithString: _dataArr[i]] placeholderImage:nil];
                [_scrollView addSubview: imageV];
//                UIImageView *imageV = [[UIImageView alloc] initWithFrame:CGRectMake(_width * i, 0, _width, _height)];
//                imageV.image = [UIImage imageNamed:_dataArr[i]];
//                [_scrollView addSubview: imageV];
            }
        }
    }
    return _scrollView;
}

-(PageControl*)pageConl {
    if (!_pageConl) {
        _pageConl = [[PageControl alloc] initWithFrame:CGRectMake(5, _heightPag, _width - 10, 10)];
        _pageConl.numberOfPages = _dataArr.count - 2;
        _pageConl.currentPage = 0;
    }
    return _pageConl;
}

- (void)scrollViewDidEndDecelerating:(UIScrollView *)scrollView {
    CGPoint curOfset = scrollView.contentOffset;
    if (curOfset.x == (_dataArr.count - 1) * _width) {
        _scrollView.contentOffset = CGPointMake(_width, 0);
    }
    if (curOfset.x == 0) {
        _scrollView.contentOffset = CGPointMake((_dataArr.count - 2) * _width, 0);
    }
    CGPoint newPoint = _scrollView.contentOffset;
    NSInteger temPage = newPoint.x / _width;
    _pageConl.currentPage = temPage - 1;
}


@end
