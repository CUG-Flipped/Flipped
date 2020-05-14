//
//  HomeViewController.m
//  Tinder
//
//  Created by Layer on 2020/5/8.
//  Copyright © 2020 Layer. All rights reserved.
//

#import "HomeViewController.h"

@interface HomeViewController () <HomeTopDelegate, HomeBottomDelegate, CardViewDelegate>

@property (strong, nonatomic) HomeTopView* topView;
@property (strong, nonatomic) HomeBottomView* bottomView;
@property (strong, nonatomic) CardView* cardView;

@property (strong, nonatomic) NSMutableArray* dataSources; // 存放图片
@property (strong, nonatomic) NSMutableArray* cards; // 存放CardView
@property (strong, nonatomic) CardView* topCard;  // 顶CardView
@property (strong, nonatomic) CardView* bottomCard; // 底CardView

@property (strong, nonatomic) NSMutableArray* listArr; // 存放当前请求获得的数据

@property (assign, nonatomic) NSInteger page; // 当前页

@property (assign, nonatomic) NSString* type; // 请求类型

@end

@implementation HomeViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    self.cards = [[NSMutableArray alloc] init];
    self.dataSources = [[NSMutableArray alloc] initWithObjects:@"background_1", @"background_2", @"background_3", @"background_4", nil];
    
//    self.view.backgroundColor = [UIColor orangeColor];
//    NSLog(@"%@",[UserDefaults getUserDefaults:@"token"]);
    _page = 0;
    _type = 0;
    [self requestHomeDatas];
    [self creatTopView];
    [self creatBottomView];
//    [self creatCardView];
}

- (void)requestHomeDatas {
    [HttpData requestHomeListToken:[UserDefaults getUserDefaults:@"token"] page:_page pageSize:5 success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"] integerValue] == 200) {
                self.listArr = json[@"data"];
                [self creatCardView];
            }
            else {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:json[@"message"]];
            }
        }
    } failure:^(NSError * _Nonnull err) {
        NSLog(@"%@",err);
    }];
}

- (void)creatTopView {
    self.topView = [[HomeTopView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH, self.navigationController.navigationBar.frame.size.height + 80)];
    self.topView.backgroundColor = [UIColor whiteColor];
    [self.view addSubview: self.topView];
    self.topView.delegate = self;
}

- (void)creatBottomView {
    self.bottomView = [[HomeBottomView alloc] init];
    self.bottomView.backgroundColor = [UIColor whiteColor];
    [self.view addSubview:self.bottomView];
    _bottomView.sd_layout.leftSpaceToView(self.view, 0).rightSpaceToView(self.view, 0).bottomSpaceToView(self.view, 0).heightIs(110);
    self.bottomView.delegate = self;
}

- (void)creatCardView {
//    for (int i = 0; i < self.dataSources.count; i++) {
    for (int i = 0; i < self.listArr.count; i++) {
        self.cardView = [[CardView alloc] initWithFrame:CGRectMake(0, 0, SCREEN_WIDTH - 30, SCREEN_HEIGHT - 260)];
        self.cardView.tag = i;
        self.cardView.delegate = self;
        
        self.cardView.backgroundColor = [UIColor whiteColor];
        self.cardView.center = CGPointMake(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2 + 8);
        
        int index = i;
//        if (index == self.dataSources.count) {
        if (index == self.listArr.count) {
            index = 0;
        }
        
        _cardView.transform = CGAffineTransformMakeScale(1 - 0.1 * index, 1 - 0.1 * index); // 缩放
        
//        _cardView.image.image = [UIImage imageNamed: self.dataSources[i]];
        [_cardView homeArrData:_listArr[i]];
        
        [self.cards addObject:_cardView];
        
        [self.view addSubview:_cardView];
        [self.view sendSubviewToBack:_cardView];
        
        UIPanGestureRecognizer *pan = [[UIPanGestureRecognizer alloc] initWithTarget:self action:@selector(panHandle:)];
        [_cardView addGestureRecognizer:pan];
        
        _cardView.userInteractionEnabled = NO;
        
        if (i == 0) {
            _cardView.userInteractionEnabled = YES;
            self.topCard = _cardView;
        }
        if (i == self.dataSources.count - 1) {
            self.bottomCard = _cardView;
        }
    }
}

- (void)panHandle:(UIPanGestureRecognizer*)pan {
    CardView* cv = ( CardView*)pan.view;
    if (pan.state == UIGestureRecognizerStateBegan) {

    }
    else if (pan.state == UIGestureRecognizerStateChanged) {
        // 手指移动后，相对位置偏移
        CGPoint transLocation = [pan translationInView:cv];
        // 修改坐标中心
        cv.center = CGPointMake(cv.center.x + transLocation.x, cv.center.y + transLocation.y);
        // 移动旋转弧度
        CGFloat xOffset = (cv.center.x - SCREEN_WIDTH / 2.0) / (SCREEN_WIDTH / 2.0);
        CGFloat rotation = 3.14 / 4.0 * xOffset;
        cv.transform = CGAffineTransformMakeRotation(rotation);
        
        if (cv.frame.origin.x < 5) {
            cv.like.alpha = 0;
            cv.dislike.alpha = 1;
        }
        else {
            cv.like.alpha = 1;
            cv.dislike.alpha = 0;
        }
        
        // 移动动画
        [pan setTranslation:CGPointZero inView:cv];
        [self changeCardView:fabs(xOffset)];
    }
    else if (pan.state == UIGestureRecognizerStateEnded) {
        if (cv.center.x > 100 && cv.center.x < SCREEN_WIDTH - 100) { // 复位
            [UIView animateWithDuration:0.3 animations:^{
                cv.center = CGPointMake(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2 + 8);
                cv.transform = CGAffineTransformMakeRotation(0);
                [self changeCardView:0];
            }];
        }
        else if (cv.center.x <= 100) { // 偏左
            [UIView animateWithDuration:0.5 animations:^{
                cv.center = CGPointMake(-500, SCREEN_HEIGHT / 2 + 8);
            }];
            [self changeCardView:1];
            [self performSelector:@selector(cardRemove) withObject:nil afterDelay:0.5];
        }
        else if (cv.center.x >= SCREEN_WIDTH - 100) { // 偏右
            [UIView animateWithDuration:0.5 animations:^{
                cv.center = CGPointMake(SCREEN_WIDTH + 500, SCREEN_HEIGHT / 2 + 8);
            }];
            [self changeCardView:1];
            [self performSelector:@selector(cardRemove) withObject:nil afterDelay:0.5];
        }
        cv.like.alpha = 0;
        cv.dislike.alpha = 0;
    }
}

- (void)changeCardView:(CGFloat)x {
    for (UIView* card in self.cards) {
        if (card != self.topCard) {
            CGFloat temp = 1 - 0.1 * card.tag + 0.1 * x;
            card.transform = CGAffineTransformMakeScale(temp, temp);
        }
    }
}

- (void)cardRemove {
    [self.cards removeObject:self.topCard];
//    [self.cards addObject:self.topCard];
    
    for (int i = 0; i < self.cards.count; i++) {
        UIView* card = self.cards[i];
        card.tag = i;
    }
    
    self.topCard.userInteractionEnabled = NO;
    self.topCard.center = CGPointMake(SCREEN_WIDTH / 2, SCREEN_HEIGHT / 2 + 8);
    self.topCard.transform = CGAffineTransformMakeScale(1 - 0.1 * 2, 1 - 0.1 * 2);
    
//    [self.view sendSubviewToBack:self.topCard];
    [self.topCard removeFromSuperview];
    
    self.bottomCard = self.topCard;
    self.topCard = self.cards.firstObject;
    self.topCard.userInteractionEnabled = YES;
    
    if (_cards.count == 0) {
        _page += 1;
        [self requestHomeDatas];
    }
    
}


- (void)clickButton:(NSInteger)btnTag {
    if (btnTag == 1) {
        NSLog(@"left");
    }
    if (btnTag == 2) {
        NSLog(@"right");
    }
}

- (void)homeBottomClickBtn:(UIButton *)btn {
    switch (btn.tag) {
        case 1: // 刷新
            [self refreshBtnClick];
            break;
        case 2: // dislake
            [self dislikeBtnClick:btn];
            break;
        case 3: // 收藏
            _type = @"1";
            [self typeBtnClick];
            [self dislikeBtnClick:btn];
            break;
        case 4: // like
            _type = @"2";
            [self typeBtnClick];
            [self likeBtnClick:btn];
            break;
        case 5: // 聊天
            _type = @"3";
            [self typeBtnClick];
            break;
        default:
            break;
    }
}

- (void)likeBtnClick:(UIButton*)btn {
    CardView *temp = self.topCard;
    temp.like.alpha = 1;
    [UIView animateWithDuration:0.75 animations:^{
        btn.enabled = NO;
        self.topCard.center = CGPointMake(2 * SCREEN_WIDTH, SCREEN_HEIGHT / 2);
        self.topCard.transform = CGAffineTransformMakeRotation(3.14 / 2);
        [self changeCardView:1];
    } completion:^(BOOL finished) {
        btn.enabled = YES;
        temp.like.alpha = 0;
        [self cardRemove];
    }];
}

- (void)dislikeBtnClick:(UIButton*)btn {
    CardView *temp = self.topCard;
    temp.dislike.alpha = 1;
    [UIView animateWithDuration:0.75 animations:^{
        btn.enabled = NO;
        self.topCard.center = CGPointMake(-2 * SCREEN_WIDTH, SCREEN_HEIGHT / 2);
        self.topCard.transform = CGAffineTransformMakeRotation(-3.14 / 2);
        [self changeCardView:1];
    } completion:^(BOOL finished) {
        btn.enabled = YES;
        temp.dislike.alpha = 0;
        [self cardRemove];
    }];
}

- (void)refreshBtnClick {
    for (int i = 0; i < self.cards.count; i++) {
        CardView* card = self.cards[i];
        [card removeFromSuperview];
    }
    [self.cards removeAllObjects];
    _page += 1;
    [self requestHomeDatas];
}

- (void)typeBtnClick {
    
    CardView* card = self.cards[0];
    
    [HttpData requestHomeCollectToken:[UserDefaults getUserDefaults:@"token"] user_id: self.listArr[card.tag][@"id"] type:_type success:^(id  _Nonnull json) {
        if ([json isKindOfClass:[NSDictionary class]]) {
            if ([json[@"status"] integerValue] == 200) {
                if ([self->_type isEqualToString:@"1"]) {
                    [SVProgressHUD dismissWithDelay:1.5];
                    [SVProgressHUD showSuccessWithStatus:@"已收藏"];
                }
                else if ([self->_type isEqualToString:@"2"]) {
                    // 喜欢
                }
                else if ([self->_type isEqualToString:@"3"]) {
                    // 聊天
                }
            }
            else {
                [SVProgressHUD dismissWithDelay:1.5];
                [SVProgressHUD showErrorWithStatus:json[@"message"]];
            }
        }
    } failure:^(NSError * _Nonnull err) {
        [SVProgressHUD dismissWithDelay:1.5];
        [SVProgressHUD showErrorWithStatus:@"操作失败！"];
    }];
}

- (void) clickInfoBtn:(UIButton*)btn {
    NSLog(@"123");
}

@end
