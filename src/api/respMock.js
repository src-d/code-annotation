const smallDiff = `diff --git a/src/components/uast/Node.js b/src/components/uast/Node.js
index 933d16c..e27afe3 100644
--- a/src/components/uast/Node.js
+++ b/src/components/uast/Node.js
@@ -1,295 +1,299 @@
 import React, { Component } from 'react';
 import styled from 'styled-components';
 import { font, border, background } from '../../styling/variables';
 import { connect } from 'react-redux';
+import { nodeToggle } from '../../state/ast';

 const INDENT_SIZE = 20;
 const WHITE_SPACE = 5;
 const LINE_HEIGHT = 29;

 const PROPERTY_VALUE_SEPARATOR = "':'";
 const STRING_LIMITER = "'\\27'";
 const COLLAPSIBLE_COLLAPSED = "'+'";
 const COLLAPSIBLE_EXTENDED = "'-'";

 export const StyledItem = styled.div\`
   margin-left: \${INDENT_SIZE}px;
   min-width: 400px;
   background: \${props => (props.highlighted ? background.highlight : 'none')};
 \`;

 const StyledTitle = styled.div\`
   min-height: \${LINE_HEIGHT}px;

   &::before {
     content: '';
     width: 15px;
     display: inline-block;
   }
 \`;

 export const StyledCollapsibleTitle = StyledTitle.extend\`
   cursor: pointer;

   &::before {
     content: \${props =>
       props.collapsed ? COLLAPSIBLE_COLLAPSED : COLLAPSIBLE_EXTENDED};
     color: \${props =>
       props.collapsed ? font.color.accentDark : font.color.accentLight};
   }

   &:hover > span {
     border-bottom: 1px dashed \${border.light};
   }
 \`;

 export const StyledCollapsibleContent = styled.div\`
   display: \${props => (props.collapsed ? 'none' : 'block')};
   border-left: 1px solid \${border.smooth};
   margin-left: 4px;
 \`;

 const Label = styled.summary\`
   display: inline;
   font-family: \${font.family.code};
   color: \${font.color.light};
 \`;

 const StyledValue = styled.span\`
   color: \${props =>
     props.type === 'number' || props.type === 'object'
       ? font.color.accentDark
       : props.type === 'string' ? font.color.accentLight : font.color.dark};

   &::before {
     content: \${props => (props.type === 'string' ? STRING_LIMITER : '')};
   }

   &::after {
     content: \${props => (props.type === 'string' ? STRING_LIMITER : '')};
   }
 \`;

 const StyledPropertyName = styled.span\`
   margin-right: \${WHITE_SPACE}px;

   &::after {
     content: \${PROPERTY_VALUE_SEPARATOR};
     color: \${font.color.light};
   }
 \`;

 export function Value({ value }) {
   return (
     <StyledValue type={typeof value}>
       {value !== null ? value : 'null'}
     </StyledValue>
   );
 }

 export function PropertyName({ name }) {
   return name ? <StyledPropertyName>{name}</StyledPropertyName> : null;
 }

 export class Node extends Component {
   constructor() {
     super();

     this.onNodeSelected = this.onNodeSelected.bind(this);
   }

   onNodeSelected(e) {
     const { StartPosition: start, EndPosition: end } = this.props.node;

     let from, to;
     if (start && start.Line && start.Col) {
       from = {
         line: start.Line - 1,
         ch: start.Col - 1,
       };
     }

     if (end && end.Line && end.Col) {
       to = {
         line: end.Line - 1,
         ch: end.Col - 1,
       };
     }

     this.props.onNodeSelected && this.props.onNodeSelected(from, to);
     e.stopPropagation();
   }

   render() {
-    const { node, showLocations } = this.props;
+    const { node, showLocations, toggle } = this.props;

     if (!node) {
       return null;
     }

     return (
       <CollapsibleItem
         label="Node"
         highlighted={node.highlighted}
         collapsed={!node.expanded}
+        toggle={() => toggle(node.id)}
         onNodeSelected={this.onNodeSelected}
       >
         <Property name="internal_type" value={node.InternalType} />
         <Properties properties={node.Properties} />
         <Property name="token" value={node.Token} />
         {showLocations ? (
           <Position name="start_position" position={node.StartPosition} />
         ) : null}
         {showLocations ? (
           <Position name="end_position" position={node.EndPosition} />
         ) : null}
         <Roles roles={node.Roles} />
         <Children
           ids={node.Children}
           onNodeSelected={this.props.onNodeSelected}
         />
       </CollapsibleItem>
     );
   }
 }

 const mapStateToProps = (state, ownProps) => {
   return {
     node: state.code.ast[ownProps.id],
     showLocations: state.options.showLocations,
   };
 };

-const ConnectedNode = connect(mapStateToProps)(Node);
+const ConnectedNode = connect(mapStateToProps, {
+  toggle: nodeToggle,
+})(Node);
 export default ConnectedNode;

 export function Properties({ properties }) {
   if (properties && Object.keys(properties).length > 0) {
     return (
       <CollapsibleItem name="properties" label="map<string, string>">
         {Object.keys(properties).map((name, i) => (
           <Property key={i} name={name} value={properties[name]} />
         ))}
       </CollapsibleItem>
     );
   }

   return null;
 }

 export function Property({ name, value }) {
   if (typeof value !== 'undefined') {
     return (
       <StyledItem>
         <StyledTitle>
           {name ? <PropertyName name={name} /> : null}
           <Value value={value} />
         </StyledTitle>
       </StyledItem>
     );
   }

   return null;
 }

 export class Children extends Component {
   render() {
     const { ids, onNodeSelected } = this.props;

     if (Array.isArray(ids)) {
       return (
         <CollapsibleItem name="children" label="[]Node" collapsed={false}>
           {ids.map(id => (
             <ConnectedNode key={id} id={id} onNodeSelected={onNodeSelected} />
           ))}
         </CollapsibleItem>
       );
     }

     return null;
   }
 }

 function coordinates(position) {
   if (!position) {
     return [];
   }

   const values = ['Offset', 'Line', 'Col'];

   return values
     .filter(name => typeof position[name] !== 'undefined')
     .map((name, i) => (
       <Property key={i} name={name.toLowerCase()} value={position[name]} />
     ));
 }

 export function Position({ name, position }) {
   const coords = coordinates(position);
   if (position && coordinates.length > 0) {
     return (
       <CollapsibleItem name={name} label="Position">
         {coords}
       </CollapsibleItem>
     );
   }

   return <Property name={name} value={null} />;
 }

 export function Roles({ roles }) {
   if (Array.isArray(roles)) {
     return (
       <CollapsibleItem name="roles" label="[]Role">
         {roles.map((role, i) => <Property key={i} value={role} />)}
       </CollapsibleItem>
     );
   }

   return null;
 }

 export class CollapsibleItem extends Component {
   constructor(props) {
     super(props);
     // component can be controlled or uncontrolled depends on did we pass collapsed state on init or not
     // similar to react input components
     this.controlled = typeof this.props.collapsed !== 'undefined';
     this.state = {
       collapsed: this.props.collapsed,
     };
   }

   componentWillReceiveProps(nextProps) {
     if (!this.controlled) {
       return;
     }
     this.setState({ collapsed: nextProps.collapsed });
   }

   toggle() {
-    this.setState({ collapsed: !this.state.collapsed });
-  }
-
-  expand() {
-    this.setState({ collapsed: false });
+    if (!this.controlled) {
+      this.setState({ collapsed: !this.state.collapsed });
+      return;
+    }
+    this.props.toggle();
   }

   render() {
     const { name, label, children, onNodeSelected, highlighted } = this.props;
     return (
       <StyledItem onMouseMove={onNodeSelected} highlighted={highlighted}>
         <StyledCollapsibleTitle
           collapsed={this.state.collapsed}
           onClick={this.toggle.bind(this)}
         >
           <PropertyName name={name} />
           <Label>{label}</Label>
         </StyledCollapsibleTitle>
         <StyledCollapsibleContent collapsed={this.state.collapsed}>
           {children}
         </StyledCollapsibleContent>
       </StyledItem>
     );
   }
 }`;

export default smallDiff;
