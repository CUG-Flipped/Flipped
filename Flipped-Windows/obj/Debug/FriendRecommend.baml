   M S B A M L   `   `   `  ����  Flipped-Win10#    Flipped_Win10.FriendRecommend   5      W SPresentationCore, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35\ XPresentationFramework, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35R NWindowsBase, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35D 9http://schemas.microsoft.com/winfx/2006/xaml/presentation    V RMaterialDesignThemes.Wpf, Version=3.1.2.1465, Culture=neutral, PublicKeyToken=nullFmaterialDesign1http://materialdesigninxaml.net/winfx/xaml/themes  5      :d2http://schemas.microsoft.com/expression/blend/2008  5      R NSystem.Xaml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089:x,http://schemas.microsoft.com/winfx/2006/xaml     5      Jmc;http://schemas.openxmlformats.org/markup-compatibility/2006     5      -   5          � Title$  FriendRecommend��5	   	   $	��435��6!   $	��425��6.    � WindowStyle$
 None=�6:    � AllowsTransparency .  � WindowStartupLocation$ CenterScreen=�5
       � 
ResizeMode$ NoResize=�6*    � Icon$ win.ico��6q   '  !System.Windows.Shell.WindowChrome   WindowChrome 5        5   
      GlassFrameThickness$ 0q�6      ResizeBorderThickness$ 10q�6/   	   CaptionHeight$	 50Q�6J   ��5      �� 5   
   
 �� 
TargetType"
 �6   .R�5      �  �Name$ backGrid��6   .+�5      ��  ��Name$ contentBorder��6   $��10q�60    �� CornerRadius$ 10��6<   ���#fefefc J� Effect 5      3  -System.Windows.Media.Effects.DropShadowEffect  5         
BlurRadius$ 6Q�6+      ShadowDepth$ 0Q�6:      Color$ #B7B7B7��6J    J� OpacityMask 5      !� 5       !� Visual 5      �� 5   "   
���   � J� SnapsToDevicePixels .  5   )   �� *CornerRadius+ �� RelativeSource ��  �� AncestorType" ����5   )   �� *ActualWidth+ �� " ����5   )   �� *ActualHeight+ �� " ��5      ._�5       �� .u�5!      �� 5"      5#      5$       �� Triggers 5%      P� 5&       P� Property �@�IsActive P� Value . . �5'      ��  �� 
TargetName$ contentBorder��6    �� Property �@J�Effect �� Value 5(        5)      $ 9Q�6/   $ 0Q�6>   $ #B2B2B2��6N   5+      5,      5.      .��51      � .+�52   
   �� -   5        $
��Left=�52      $	��Top=�6G      WindowCloseButtonStyle#��[  6_   $��30��6�   $��60��6�      IsHitTestVisibleInChrome . $��	345,0,0,0q�6�   �� 53   
   -   5         ��Name$ recommendFriendAvtar��53      $
��Left=�6.   $	��325��6I   $��	20,30,0,0q�6V   $	��Top=�6i   $	��330��6�   �� 55   
    �� Orientation$ 
Horizontal=�6   $��	0,365,0,0q�6.   .+�56      �� -   5         . *MaterialDesignFloatingActionMiniDarkButton#��[ 56        � ToolTip$  Refresh��6Y   $��20,0,0,0q�6k   .��57      '  !MaterialDesignThemes.Wpf.PackIcon  !   Kind$! Refresh=�6*   $��24��69   $��24��6E   58      �� 59      -   5        #��[ 59      $  Dislike��6Y   $��20,0,0,0q�6k   .��5:        $! HeartOff=�6*   $��24��6:   $��24��6F   5;      �� 5<      -   5        #��[ 5<      $  Collect��6Y   $��20,0,0,0q�6k   .��5=        -   5        "  Name$" 
collectBtn��5=   *   $! 
StarBorder=�6>   $��24��6P   $��24��6\   5>      �� 5?      -   5        #��[ 5?      $
  like��6Y   $��20,0,0,0q�6h   .��5@        $! 
HeartPulse=�6*   $��24��6<   $��24��6H   5A      �� 5B      -	   5        #��[ 5B      $  detailed information��6Y   $��20,0,0,0q�6x   .��5C        $! InformationVariant=�6*   $��24��6D   $��24��6P   5D      �� 5E      -
   5        #��[ 5E      $  
quick chat��6Y   $��20,0,0,0q�6n   .��5F        $! ChatProcessing=�6*   $��24��6@   $��24��6L   5G      5H      )  #MaterialDesignThemes.Wpf.DrawerHost  5I   
   -   5        #  Name$# Drawer��5I   $     MaterialDesignDivider#���  64   $
��Left=�6j   $	��Top=�6�   $	��415��6�   $	��405��6�   $   RightDrawerContent$ 5J      �� 5K      
�������$	��210��60   $
��Left=�6<   .+�5L      #  MaterialDesignThemes.Wpf.Chip  -   5        %  Name$% UserNameChip��5L   *   &   Icon$& U��6@   $  Username��6I   $��30,10q�6\     5M      -   5        $% 	EmailChip��5M   *   $& E��6=   $  Email��6F   $��30,10q�6V     5N      -   5        $% NameChip��5N   *   $& R��6<   $  RealName��6E   $��30,10q�6X     5O      -   5        $% ProfessionChip��5O   *   $& P��6B   $  
Profession��6K   $��30,10q�6`     5P      -   5        $% AgeChip��5P   *   $& A��6;   $	  Age��6D   $��30,10q�6R     5Q      -   5        $% 
RegionChip��5Q   *   $& R��6>   $  Region��6G   $��30,10q�6X     5R      -   5        $% 	HobbyChip��5R   *   $& H��6=   $  Hobby��6F   $��30,10q�6V   �� 5S      -   5        $��Center=�5S      $
��0,10q�6:   .��6a   Back6g   5T      5V      5X      5Y      